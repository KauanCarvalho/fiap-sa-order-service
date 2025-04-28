package usecase_test

import (
	"testing"
	"time"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateOrderUseCase_Run(t *testing.T) {
	t.Run("should successfully update order status", func(t *testing.T) {
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

		err := ds.CreateOrder(ctx, order)
		require.NoError(t, err)

		var createdOrder entities.Order
		err = sqlDB.Preload("OrderItems").First(&createdOrder, order.ID).Error
		require.NoError(t, err)

		updateOrderUC := usecase.NewUpdateOrderUseCase(ds)
		err = updateOrderUC.Run(ctx, createdOrder.ID, "shipped")
		require.NoError(t, err)

		var updatedOrder entities.Order
		err = sqlDB.Preload("OrderItems").First(&updatedOrder, createdOrder.ID).Error
		require.NoError(t, err)

		assert.Equal(t, "shipped", updatedOrder.Status)
	})

	t.Run("should return error if order does not exist", func(t *testing.T) {
		prepareTestDatabase()

		updateOrderUC := usecase.NewUpdateOrderUseCase(ds)
		err := updateOrderUC.Run(ctx, 99999, "shipped")
		require.Error(t, err)
	})
}
