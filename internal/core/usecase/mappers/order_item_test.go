package mappers_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/mappers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToOrderItemDTO(t *testing.T) {
	t.Run("should map line item correctly", func(t *testing.T) {
		orderItem := entities.OrderItem{
			ID:       1,
			OrderID:  1,
			SKU:      "ABC123",
			Quantity: 2,
			Price:    100.50,
		}

		result := mappers.ToOrderItemDTO(orderItem)

		require.Equal(t, "ABC123", result.SKU)
		assert.Equal(t, 2, result.Quantity)
		assert.InEpsilon(t, 100.50, result.Price, 0.01)
	})

	t.Run("should map line item with zero quantity and price", func(t *testing.T) {
		orderItem := entities.OrderItem{
			ID:       2,
			OrderID:  1,
			SKU:      "XYZ789",
			Quantity: 0,
			Price:    0.00,
		}

		result := mappers.ToOrderItemDTO(orderItem)

		require.Equal(t, "XYZ789", result.SKU)
		assert.Equal(t, 0, result.Quantity)
		assert.Zero(t, result.Price)
	})
}
