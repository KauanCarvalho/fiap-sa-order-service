package mappers_test

import (
	"testing"
	"time"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/mappers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToOrderDTO(t *testing.T) {
	t.Run("should map order correctly with items", func(t *testing.T) {
		createdAt := time.Now().Add(-time.Hour)
		updatedAt := time.Now()

		order := entities.Order{
			ID:        1,
			ClientID:  1,
			Status:    "Pending",
			Price:     200.00,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			OrderItems: []entities.OrderItem{
				{
					ID:       1,
					OrderID:  1,
					SKU:      "ABC123",
					Quantity: 2,
					Price:    100.50,
				},
				{
					ID:       2,
					OrderID:  1,
					SKU:      "XYZ789",
					Quantity: 1,
					Price:    99.50,
				},
			},
		}

		result := mappers.ToOrderDTO(order)

		require.Equal(t, uint(1), result.ID)
		assert.Equal(t, uint(1), result.ClientID)
		assert.Equal(t, "Pending", result.Status)
		assert.InEpsilon(t, 200.00, result.Price, 0.01)
		assert.Equal(t, createdAt, result.CreatedAt)
		assert.Equal(t, updatedAt, result.UpdatedAt)

		require.Len(t, result.Items, 2)
		assert.Equal(t, "ABC123", result.Items[0].SKU)
		assert.Equal(t, 2, result.Items[0].Quantity)
		assert.InEpsilon(t, 100.50, result.Items[0].Price, 0.01)

		assert.Equal(t, "XYZ789", result.Items[1].SKU)
		assert.Equal(t, 1, result.Items[1].Quantity)
		assert.InEpsilon(t, 99.50, result.Items[1].Price, 0.01)
	})

	t.Run("should map order with no items", func(t *testing.T) {
		createdAt := time.Now().Add(-time.Hour)
		updatedAt := time.Now()

		order := entities.Order{
			ID:         1,
			ClientID:   1,
			Status:     "Completed",
			Price:      150.00,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			OrderItems: []entities.OrderItem{},
		}

		result := mappers.ToOrderDTO(order)

		require.Equal(t, uint(1), result.ID)
		assert.Equal(t, uint(1), result.ClientID)
		assert.Equal(t, "Completed", result.Status)
		assert.InEpsilon(t, 150.00, result.Price, 0.01)
		assert.Equal(t, createdAt, result.CreatedAt)
		assert.Equal(t, updatedAt, result.UpdatedAt)

		assert.Empty(t, result.Items)
	})
}
