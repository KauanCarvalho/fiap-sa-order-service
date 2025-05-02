package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/config"
	externalErrors "github.com/KauanCarvalho/fiap-sa-order-service/internal/shared/errors"
)

type ClientInterface interface {
	AuthorizePayment(ctx context.Context, amount float64, externalReference, paymentMethod string) (*Response, error)
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

const clientTimeout = 10 * time.Second

var ErrProblemToAuthorizePayment = errors.New("problem to authorize payment")

func NewClient(cfg config.Config) *Client {
	return &Client{
		BaseURL:    cfg.PaymentServiceURL,
		HTTPClient: &http.Client{Timeout: clientTimeout},
	}
}

func (c *Client) createRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+endpoint, body)
	if err != nil {
		return nil, externalErrors.NewExternalError("error creating request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) doRequest(req *http.Request, response any) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return externalErrors.NewExternalError("error making request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusInternalServerError {
			return ErrProblemToAuthorizePayment
		}
		return externalErrors.NewExternalError("error response from server", errors.New(resp.Status))
	}

	if resp.ContentLength == 0 {
		return nil
	}

	if errDecode := json.NewDecoder(resp.Body).Decode(response); errDecode != nil {
		return externalErrors.NewExternalError("error decoding response", errDecode)
	}

	return nil
}

func (c *Client) AuthorizePayment(
	ctx context.Context,
	amount float64,
	externalReference, paymentMethod string,
) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, clientTimeout)
	defer cancel()

	body := map[string]any{
		"amount":             amount,
		"external_reference": externalReference,
		"payment_method":     paymentMethod,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, externalErrors.NewExternalError("error marshalling request body", err)
	}

	req, err := c.createRequest(
		ctx,
		http.MethodPost,
		"/v1/payments/authorize",
		io.NopCloser(bytes.NewReader(bodyBytes)))
	if err != nil {
		return nil, err
	}

	var response Response
	if errRequest := c.doRequest(req, &response); errRequest != nil {
		return nil, errRequest
	}

	return &response, nil
}
