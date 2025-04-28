package entities_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"

	"github.com/stretchr/testify/assert"
)

func TestOrder_CalculateTotal(t *testing.T) {
	t.Run("successfully calculates the total order price", func(t *testing.T) {
		order := &entities.Order{
			OrderItems: []entities.OrderItem{
				{SKU: "sku-1", Quantity: 2, Price: 10.0},
				{SKU: "sku-2", Quantity: 1, Price: 20.0},
				{SKU: "sku-3", Quantity: 3, Price: 5.0},
			},
		}

		order.CalculateTotal()

		expectedTotal := (2 * 10.0) + (1 * 20.0) + (3 * 5.0)
		assert.InEpsilon(t, expectedTotal, order.Price, 0.01)
	})

	t.Run("calculates total for empty order items", func(t *testing.T) {
		order := &entities.Order{}

		order.CalculateTotal()

		assert.Zero(t, order.Price)
	})
}
