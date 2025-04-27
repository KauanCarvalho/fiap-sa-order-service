package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	useCaseDTO "github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckoutHandler_Create(t *testing.T) {
	prepareTestDatabase()

	t.Run("successful order creation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/products/test-sku", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"price": 49.99}`))
		}))
		defer server.Close()

		engine := setupTestRouter(server.URL)

		reqBody := `{
			"client_id": 1,
			"items": [{"sku": "test-sku", "quantity": 2}]
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ginEngine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("inavalid params", func(t *testing.T) {
		input := useCaseDTO.OrderInputCreate{
			ClientID: 1,
			Items:    []useCaseDTO.OrderItemInputCreate{},
		}

		body, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ginEngine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("client not found", func(t *testing.T) {
		input := useCaseDTO.OrderInputCreate{
			ClientID: 99999,
			Items: []useCaseDTO.OrderItemInputCreate{
				{
					SKU:      "ABC123",
					Quantity: 1,
				},
			},
		}

		body, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		ginEngine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("product not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/products/INVALIDSKU", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		engine := setupTestRouter(server.URL)

		input := useCaseDTO.OrderInputCreate{
			ClientID: 1,
			Items: []useCaseDTO.OrderItemInputCreate{
				{
					SKU:      "INVALIDSKU",
					Quantity: 1,
				},
			},
		}

		body, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("product service internal server error (500)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/products/ANYSKU", r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		engine := setupTestRouter(server.URL)

		input := useCaseDTO.OrderInputCreate{
			ClientID: 1,
			Items: []useCaseDTO.OrderItemInputCreate{
				{
					SKU:      "ANYSKU",
					Quantity: 1,
				},
			},
		}

		body, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
