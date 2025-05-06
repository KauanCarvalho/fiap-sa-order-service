package bdd_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/cucumber/godog"
)

var (
	checkoutRequestBody string
	productMockServer   *httptest.Server
	paymentMockServer   *httptest.Server
)

func startMockExternalServices() {
	productMockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "INVALIDSKU") {
			http.NotFound(w, r)
			return
		}
		if strings.Contains(r.URL.Path, "FORCE500") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"price": 49.99}`))
	}))

	paymentMockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "FORCE_PAYMENT_500") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"amount": 99.98,
			"status": "pending",
			"external_reference": "order-123",
			"provider": "mockpay",
			"payment_method": "pix",
			"qr_code": "some-qr-code"
		}`))
	}))

	engine = setupTestRouter(productMockServer.URL, paymentMockServer.URL)
}

func stopMockExternalServices() {
	productMockServer.Close()
	paymentMockServer.Close()
}

func givenAClientWithIDExists(id string) error {
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid client ID: %v", err)
	}
	client := entities.Client{
		ID: uint(parsedID),
	}
	return sqlDB.Create(&client).Error
}

func iSendAPOSTRequestToCheckoutWithBody(body string) error {
	checkoutRequestBody = body
	req, err := http.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewBufferString(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	return nil
}

func theResponseShouldContainFieldWithValue(field, expected string) error {
	body := recorder.Body.String()
	if !strings.Contains(body, fmt.Sprintf(`"%s":"%s"`, field, expected)) &&
		!strings.Contains(body, fmt.Sprintf(`"%s":%s`, field, expected)) {
		return fmt.Errorf("expected response to contain field %q with value %q, got %s", field, expected, body)
	}
	return nil
}

func theProductServiceIsReturning500() error {
	productMockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "FORCE500") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"price": 49.99}`))
	}))

	engine = setupTestRouter(productMockServer.URL, paymentMockServer.URL)
	return nil
}

func thePaymentServiceIsReturning500() error {
	paymentMockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "FORCE_PAYMENT_500") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{
			"amount": 99.98,
			"status": "pending",
			"external_reference": "order-123",
			"provider": "mockpay",
			"payment_method": "pix",
			"qr_code": "some-qr-code"
		}`))
	}))

	engine = setupTestRouter(productMockServer.URL, paymentMockServer.URL)
	return nil
}

func InitializeScenarioCheckout(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ResetAndLoadFixtures()
		startMockExternalServices()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		stopMockExternalServices()
		return ctx, nil
	})

	ctx.Step(`^a client with ID "([^"]*)" exists$`, givenAClientWithIDExists)
	ctx.Step(`^I send a POST request to "/api/v1/checkout" with body:$`, iSendAPOSTRequestToCheckoutWithBody)
	ctx.Step(`^the response status should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response should contain "([^"]*)" with value "([^"]*)"$`, theResponseShouldContainFieldWithValue)
	ctx.Step(`^the product service is returning 500$`, theProductServiceIsReturning500)
	ctx.Step(`^the payment service is returning 500$`, thePaymentServiceIsReturning500)
}
