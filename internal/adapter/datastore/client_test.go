package datastore_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateClient(t *testing.T) {
	t.Run("successfully creates a client", func(t *testing.T) {
		prepareTestDatabase()

		client := &entities.Client{
			Name: "João Silva",
			CPF:  "12345678900",
		}

		err := ds.CreateClient(ctx, client)
		require.NoError(t, err)
		assert.NotZero(t, client.ID)
	})

	t.Run("fail to create client with duplicate CPF", func(t *testing.T) {
		prepareTestDatabase()

		client := &entities.Client{
			Name: "Maria Oliveira",
			CPF:  "12345678900",
		}

		err := ds.CreateClient(ctx, client)
		require.NoError(t, err)

		client.ID = 0

		err = ds.CreateClient(ctx, client)
		require.Error(t, err)
		require.ErrorIs(t, err, datastore.ErrExistingRecord)
	})
}

func TestGetClientByCpf(t *testing.T) {
	prepareTestDatabase()

	t.Run("successfully gets client by CPF", func(t *testing.T) {
		client := &entities.Client{
			Name: "Fernanda Souza",
			CPF:  "98765432100",
		}
		require.NoError(t, sqlDB.Create(client).Error)

		found, err := ds.GetClientByCpf(ctx, "98765432100")
		require.NoError(t, err)
		assert.Equal(t, client.Name, found.Name)
		assert.Equal(t, client.CPF, found.CPF)
	})

	t.Run("returns error when client not found", func(t *testing.T) {
		found, err := ds.GetClientByCpf(ctx, "00000000000")
		assert.Nil(t, found)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
