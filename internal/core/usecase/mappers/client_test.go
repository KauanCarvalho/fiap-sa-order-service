package mappers_test

import (
	"testing"
	"time"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/mappers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToClientDTO(t *testing.T) {
	t.Run("should map client correctly", func(t *testing.T) {
		createdAt := time.Now().Add(-time.Hour)
		updatedAt := time.Now()

		client := entities.Client{
			ID:        1,
			Name:      "Foo",
			CPF:       "12345678900",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		result := mappers.ToClientDTO(client)

		require.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Foo", result.Name)
		assert.Equal(t, "12345678900", result.CPF)
		assert.Equal(t, createdAt, result.CreatedAt)
		assert.Equal(t, updatedAt, result.UpdatedAt)
	})

	t.Run("should handle client with empty name", func(t *testing.T) {
		client := entities.Client{
			ID:   2,
			Name: "",
			CPF:  "00000000000",
		}

		result := mappers.ToClientDTO(client)

		require.Equal(t, uint(2), result.ID)
		assert.Empty(t, result.Name)
		assert.Equal(t, "00000000000", result.CPF)
	})
}
