package datastore_test

import (
	"testing"
	"time"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	prepareTestDatabase()

	t.Run("successfully creates an order", func(t *testing.T) {
		client, err := ds.GetClientByCpf(ctx, "20681201002")
		require.NoError(t, err)

		order := &entities.Order{
			ClientID:  client.ID,
			Status:    "pending",
			Price:     200.00,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			OrderItems: []entities.OrderItem{
				{
					SKU:      "SKU123",
					Quantity: 2,
					Price:    100.00,
				},
			},
		}

		err = ds.CreateOrderTx(ctx, ds.GetDB(), order)
		require.NoError(t, err)

		assert.NotZero(t, order.ID)

		var createdOrder entities.Order
		err = sqlDB.Preload("OrderItems").First(&createdOrder, order.ID).Error
		require.NoError(t, err)

		assert.Len(t, createdOrder.OrderItems, 1, "Order must have 1 item")
		assert.Equal(t, "SKU123", createdOrder.OrderItems[0].SKU)
		assert.Equal(t, 2, createdOrder.OrderItems[0].Quantity)
		assert.InEpsilon(t, 100.00, createdOrder.OrderItems[0].Price, 0.01)
	})

	t.Run("fail to create order with missing required fields", func(t *testing.T) {
		order := &entities.Order{
			ClientID: 0,
			Status:   "pending",
			Price:    200.00,
		}

		err := ds.CreateOrderTx(ctx, ds.GetDB(), order)
		require.Error(t, err)
	})
}

func TestUpdateOrderStatus(t *testing.T) {
	prepareTestDatabase()

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

	err := ds.CreateOrderTx(ctx, ds.GetDB(), order)
	require.NoError(t, err)
	require.NotZero(t, order.ID)

	t.Run("successfully updates order status", func(t *testing.T) {
		newStatus := "ready"
		err = ds.UpdateOrderStatus(ctx, order.ID, newStatus)
		require.NoError(t, err)

		var updatedOrder entities.Order
		err = sqlDB.First(&updatedOrder, order.ID).Error
		require.NoError(t, err)

		assert.Equal(t, newStatus, updatedOrder.Status)
	})

	t.Run("fail to update non-existing order", func(t *testing.T) {
		err = ds.UpdateOrderStatus(ctx, 99999, "ready")
		require.ErrorIs(t, err, datastore.ErrOrderNotFound)
	})
}

func TestGetPaginatedOrders(t *testing.T) {
	prepareTestDatabase()

	t.Run("successfully retrieves paginated orders", func(t *testing.T) {
		filter := ports.Filter{
			Limit:  2,
			Offset: 0,
		}

		result, err := ds.GetPaginatedOrders(ctx, filter)
		require.NoError(t, err)
		require.Len(t, result, 2)

		assert.Equal(t, "ready", result[0].Status)
		assert.Equal(t, "preparing", result[1].Status)
	})

	t.Run("retrieves empty result with invalid offset", func(t *testing.T) {
		filter := ports.Filter{
			Limit:  10,
			Offset: 999,
		}

		result, err := ds.GetPaginatedOrders(ctx, filter)
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("successfully retrieves orders with default filter values", func(t *testing.T) {
		filter := ports.Filter{
			Limit:  0,
			Offset: -1,
		}

		result, err := ds.GetPaginatedOrders(ctx, filter)
		require.NoError(t, err)
		require.Len(t, result, 3)

		assert.Equal(t, "ready", result[0].Status)
		assert.Equal(t, "preparing", result[1].Status)
		assert.Equal(t, "pending", result[2].Status)
	})

	t.Run("successfully retrieves orders filtered by status", func(t *testing.T) {
		filter := ports.Filter{
			Limit:  3,
			Offset: 0,
		}

		result, err := ds.GetPaginatedOrders(ctx, filter)
		require.NoError(t, err)
		require.Len(t, result, 3)

		for _, order := range result {
			assert.NotEqual(t, "delivered", order.Status)
		}
	})

	t.Run("correctly orders by status and created_at", func(t *testing.T) {
		filter := ports.Filter{
			Limit:  3,
			Offset: 0,
		}

		result, err := ds.GetPaginatedOrders(ctx, filter)
		require.NoError(t, err)
		require.Len(t, result, 3)

		assert.Equal(t, "ready", result[0].Status)
		assert.Equal(t, "preparing", result[1].Status)
		assert.Equal(t, "pending", result[2].Status)
	})
}
