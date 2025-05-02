package payment_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/payment"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestAuthorizePayment(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/payments/authorize", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			bodyBytes, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			body := string(bodyBytes)
			assert.NoError(t, err)
			assert.Contains(t, body, `"amount":123.45`)
			assert.Contains(t, body, `"external_reference":"order-123"`)
			assert.Contains(t, body, `"payment_method":"pix"`)

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{
				"amount": 123.45,
				"status": "AUTHORIZED",
				"external_reference": "order-123",
				"provider": "mockpay",
				"payment_method": "pix",
				"qr_code": "some-qr-code"
			}`))
		}))
		defer server.Close()

		cfg := config.Config{PaymentServiceURL: server.URL}
		client := payment.NewClient(cfg)

		resp, err := client.AuthorizePayment(ctx, 123.45, "order-123", "pix")

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.InEpsilon(t, 123.45, resp.Amount, 0.01)
		assert.Equal(t, "AUTHORIZED", resp.Status)
		assert.Equal(t, "order-123", resp.ExternalReference)
		assert.Equal(t, "pix", resp.PaymentMethod)
		assert.Equal(t, "some-qr-code", resp.QRCode)
	})

	t.Run("server returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cfg := config.Config{PaymentServiceURL: server.URL}
		client := payment.NewClient(cfg)

		resp, err := client.AuthorizePayment(ctx, 100.0, "any-ref", "card")

		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{invalid-json}`))
		}))
		defer server.Close()

		cfg := config.Config{PaymentServiceURL: server.URL}
		client := payment.NewClient(cfg)

		resp, err := client.AuthorizePayment(ctx, 100.0, "any-ref", "card")

		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("empty response body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer server.Close()

		cfg := config.Config{PaymentServiceURL: server.URL}
		client := payment.NewClient(cfg)

		resp, err := client.AuthorizePayment(ctx, 100.0, "ref", "card")

		require.NoError(t, err)
		assert.Equal(t, &payment.Response{}, resp)
	})
}
