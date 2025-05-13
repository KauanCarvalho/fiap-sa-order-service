package dto_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/dto"

	"github.com/stretchr/testify/assert"
)

func TestValidateOrderCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		input := dto.OrderInputCreate{
			ClientID: 1,
			Items: []dto.OrderItemInputCreate{
				{
					SKU:      "item123",
					Quantity: 2,
				},
			},
		}

		err := dto.ValidateOrderCreate(input)
		assert.NoError(t, err)
	})

	t.Run("missing items", func(t *testing.T) {
		input := dto.OrderInputCreate{
			ClientID: 1,
			Items:    nil,
		}

		err := dto.ValidateOrderCreate(input)
		assert.Error(t, err)
	})

	t.Run("missing sku in item", func(t *testing.T) {
		input := dto.OrderInputCreate{
			ClientID: 1,
			Items: []dto.OrderItemInputCreate{
				{
					SKU:      "",
					Quantity: 2,
				},
			},
		}

		err := dto.ValidateOrderCreate(input)
		assert.Error(t, err)
	})

	t.Run("missing items", func(t *testing.T) {
		input := dto.OrderInputCreate{
			ClientID: 1,
			Items:    []dto.OrderItemInputCreate{},
		}

		err := dto.ValidateOrderCreate(input)
		assert.Error(t, err)
	})

	t.Run("missing quantity in item", func(t *testing.T) {
		input := dto.OrderInputCreate{
			ClientID: 1,
			Items: []dto.OrderItemInputCreate{
				{
					SKU:      "item123",
					Quantity: 0,
				},
			},
		}

		err := dto.ValidateOrderCreate(input)
		assert.Error(t, err)
	})
}
