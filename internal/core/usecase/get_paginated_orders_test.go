package usecase_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/ports"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPaginatedOrdersUseCase_Run(t *testing.T) {
	prepareTestDatabase()

	t.Run("should return paginated orders", func(t *testing.T) {
		filter := ports.Filter{Limit: 2, Offset: 0}

		result, err := gp.Run(ctx, filter)
		require.NoError(t, err)
		require.Len(t, result, 2)

		assert.Equal(t, "ready", result[0].Status)
		assert.Equal(t, "preparing", result[1].Status)
	})

	t.Run("should return empty result if no orders match the filter", func(t *testing.T) {
		filter := ports.Filter{Limit: 10, Offset: 999}

		result, err := gp.Run(ctx, filter)
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("should return orders with default filter values", func(t *testing.T) {
		filter := ports.Filter{Limit: 10, Offset: 0}

		result, err := gp.Run(ctx, filter)
		require.NoError(t, err)
		require.Len(t, result, 3)

		assert.Equal(t, "ready", result[0].Status)
		assert.Equal(t, "preparing", result[1].Status)
		assert.Equal(t, "pending", result[2].Status)
	})
}
