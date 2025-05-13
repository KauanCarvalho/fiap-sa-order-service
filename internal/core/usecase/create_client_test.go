package usecase_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateClientUseCase_Run(t *testing.T) {
	t.Run("should create client with valid data", func(t *testing.T) {
		prepareTestDatabase()

		input := dto.ClientInputCreate{
			Name: "João Silva",
			CPF:  "12345678901",
		}

		client, err := cc.Run(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, client)

		assert.Equal(t, "João Silva", client.Name)
		assert.Equal(t, "12345678901", client.CPF)
	})

	t.Run("should create client with valid data and cognitoID", func(t *testing.T) {
		prepareTestDatabase()

		input := dto.ClientInputCreate{
			Name:      "João Silva",
			CPF:       "12345678901",
			CognitoID: uuid.New().String(),
		}

		client, err := cc.Run(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, client)

		assert.Equal(t, "João Silva", client.Name)
		assert.Equal(t, "12345678901", client.CPF)
		assert.Equal(t, input.CognitoID, client.CognitoID.String)
	})

	t.Run("should return error if client with same CPF already exists", func(t *testing.T) {
		prepareTestDatabase()

		input := dto.ClientInputCreate{
			Name: "Maria Souza",
			CPF:  "07644959092", // CPF already exists in the fixtures.
		}

		client, err := cc.Run(ctx, input)
		require.Error(t, err)
		require.Nil(t, client)
		assert.ErrorIs(t, err, datastore.ErrExistingRecord)
	})
}
