package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientHandler_Create(t *testing.T) {
	prepareTestDatabase()

	t.Run("success", func(t *testing.T) {
		body := map[string]string{
			"name": "John Doe",
			"cpf":  "12345678909",
		}
		bodyJSON, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewReader(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"John Doe"`)
		assert.Contains(t, w.Body.String(), `"cpf":"12345678909"`)
	})

	t.Run("success with cognitoID", func(t *testing.T) {
		body := map[string]string{
			"name":       "John Doe",
			"cpf":        "12345678910",
			"cognito_id": "cognito-id-123",
		}
		bodyJSON, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewReader(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"John Doe"`)
		assert.Contains(t, w.Body.String(), `"cpf":"12345678910"`)
		assert.Contains(t, w.Body.String(), `"cognito_id":"cognito-id-123"`)
	})

	t.Run("invalid body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewReader([]byte(`invalid json`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"field":"body"`)
		assert.Contains(t, w.Body.String(), `"message":"Invalid request body"`)
	})

	t.Run("missing fields validation error", func(t *testing.T) {
		body := map[string]string{}
		bodyJSON, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewReader(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("conflict cpf exists", func(t *testing.T) {
		body := map[string]string{
			"name": "Jane Doe",
			"cpf":  "98765432100",
		}
		bodyJSON, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewReader(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		req, _ = http.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewReader(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), `"field":"cpf"`)
		assert.Contains(t, w.Body.String(), `"message":"cpf already exists"`)
	})
}

func TestClientHandler_GetClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body := map[string]string{
			"name": "Alice Smith",
			"cpf":  "11122233344",
		}
		bodyJSON, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewReader(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		req, _ = http.NewRequest(http.MethodGet, "/api/v1/clients/11122233344", nil)
		w = httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Alice Smith"`)
		assert.Contains(t, w.Body.String(), `"cpf":"11122233344"`)
	})

	t.Run("not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/clients/00000000000", nil)
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
