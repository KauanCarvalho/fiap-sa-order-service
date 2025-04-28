package product_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/product"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestGetProduct(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/products/test-sku", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			w.WriteHeader(http.StatusOK)
			jsonResponse := `{
				"price": 99.99
			}`
			_, err := w.Write([]byte(jsonResponse))
			assert.NoError(t, err)
		}))
		defer server.Close()

		cfg := config.Config{ProductServiceURL: server.URL}
		client := product.NewClient(cfg)

		resp, err := client.GetProduct(ctx, "test-sku")

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.InEpsilon(t, 99.99, resp.Price, 0.01)
	})

	t.Run("product not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		cfg := config.Config{ProductServiceURL: server.URL}
		client := product.NewClient(cfg)

		resp, err := client.GetProduct(ctx, "nonexistent-sku")

		require.Error(t, err)
		require.ErrorIs(t, err, product.ErrSKUNotFound)
		assert.Nil(t, resp)
	})

	t.Run("server error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cfg := config.Config{ProductServiceURL: server.URL}
		client := product.NewClient(cfg)

		resp, err := client.GetProduct(ctx, "any-sku")

		require.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("invalid json response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{invalid-json}`))
		}))
		defer server.Close()

		cfg := config.Config{ProductServiceURL: server.URL}
		client := product.NewClient(cfg)

		resp, err := client.GetProduct(ctx, "sku-invalid-json")

		require.Error(t, err)
		assert.Nil(t, resp)
	})
}
