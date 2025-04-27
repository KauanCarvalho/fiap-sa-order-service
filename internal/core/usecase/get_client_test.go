package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetClientUseCase_Run(t *testing.T) {
	t.Run("should return client if CPF exists", func(t *testing.T) {
		prepareTestDatabase()

		client, err := gc.Run(ctx, "20681201002")
		require.NoError(t, err)
		require.NotNil(t, client)

		assert.Equal(t, "20681201002", client.CPF)
		assert.NotEmpty(t, client.Name)
	})

	t.Run("should return error if CPF does not exist", func(t *testing.T) {
		prepareTestDatabase()

		client, err := gc.Run(ctx, "00000000000")
		require.Error(t, err)
		require.Nil(t, client)
	})
}
