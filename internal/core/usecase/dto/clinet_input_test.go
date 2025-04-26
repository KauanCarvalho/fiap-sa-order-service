package dto_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/dto"

	"github.com/stretchr/testify/assert"
)

func TestValidateClientCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		input := dto.ClientInputCreate{
			Name: "Alice Smith",
			CPF:  "12345678900",
		}

		err := dto.ValidateClientCreate(input)
		assert.NoError(t, err)
	})

	t.Run("missing name", func(t *testing.T) {
		input := dto.ClientInputCreate{
			Name: "",
			CPF:  "12345678900",
		}

		err := dto.ValidateClientCreate(input)
		assert.Error(t, err)
	})

	t.Run("missing cpf", func(t *testing.T) {
		input := dto.ClientInputCreate{
			Name: "Alice Smith",
			CPF:  "",
		}

		err := dto.ValidateClientCreate(input)
		assert.Error(t, err)
	})

	t.Run("name exceeds max length", func(t *testing.T) {
		input := dto.ClientInputCreate{
			Name: generateLongString(101),
			CPF:  "12345678900",
		}

		err := dto.ValidateClientCreate(input)
		assert.Error(t, err)
	})
}

func generateLongString(length int) string {
	s := ""
	for range length {
		s += "a"
	}
	return s
}
