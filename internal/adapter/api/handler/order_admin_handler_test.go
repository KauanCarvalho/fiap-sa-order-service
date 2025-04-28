package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderAdminHandler_UpdateOrderStatus(t *testing.T) {
	prepareTestDatabase()

	t.Run("success", func(t *testing.T) {
		order := &entities.Order{
			ClientID:  1,
			Status:    "pending",
			Price:     100.00,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			OrderItems: []entities.OrderItem{
				{
					SKU:      "SKU123",
					Quantity: 2,
					Price:    50.00,
				},
			},
		}

		err := ds.CreateOrder(ctx, order)
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/admin/orders/"+strconv.Itoa(int(order.ID))+"/ready", nil)
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		var updatedOrder entities.Order
		err = sqlDB.First(&updatedOrder, order.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "ready", updatedOrder.Status)
	})

	t.Run("invalid order ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/admin/orders/invalid-id/ready", nil)
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"field":"orderID"`)
		assert.Contains(t, w.Body.String(), `"message":"Invalid order ID"`)
	})

	t.Run("invalid status", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/admin/orders/1/invalid-status", nil)
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"field":"status"`)
		assert.Contains(t, w.Body.String(), `"message":"Invalid order status"`)
	})

	t.Run("order not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/admin/orders/99999/ready", nil)
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), `"field":"orderID"`)
		assert.Contains(t, w.Body.String(), `"message":"Order not found"`)
	})
}

func TestOrderAdminHandler_GetPaginatedOrders(t *testing.T) {
	prepareTestDatabase()

	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/orders?page=0&pageSize=10", nil)
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		assert.Contains(t, w.Body.String(), "pending")
		assert.Contains(t, w.Body.String(), "preparing")
		assert.Contains(t, w.Body.String(), "ready")
		assert.NotContains(t, w.Body.String(), "delivered")
	})

	t.Run("no orders found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/orders?page=2&pageSize=10", nil)
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "[]", w.Body.String())
	})

	t.Run("invalid page", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/orders?page=invalid", nil)
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid pageSize", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/admin/orders?pageSize=invalid", nil)
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
