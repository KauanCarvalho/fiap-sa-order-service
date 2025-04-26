package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/application/mock"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"

	"github.com/stretchr/testify/require"
)

func TestDatastoreMock_Ping(t *testing.T) {
	t.Run("when PingFn is defined, it returns the custom error", func(t *testing.T) {
		expectedErr := errors.New("db unavailable")

		ds := &mock.DatastoreMock{
			PingFn: func(_ context.Context) error {
				return expectedErr
			},
		}

		err := ds.Ping(ctx)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("when PingFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		err := ds.Ping(ctx)
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
	})
}

func TestDatastoreMock_CreateClient(t *testing.T) {
	t.Run("when CreateClientFn is defined, it returns the custom error", func(t *testing.T) {
		expectedErr := errors.New("db unavailable")

		ds := &mock.DatastoreMock{
			CreateClientFn: func(_ context.Context, _ *entities.Client) error {
				return expectedErr
			},
		}

		err := ds.CreateClient(ctx, nil)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("when CreateClientFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		err := ds.CreateClient(ctx, nil)
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
	})
}

func TestDatastoreMock_GetClientByCpf(t *testing.T) {
	t.Run("when GetClientByCpfFn is defined, it returns the custom error", func(t *testing.T) {
		expectedErr := errors.New("db unavailable")

		ds := &mock.DatastoreMock{
			GetClientByCpfFn: func(_ context.Context, _ string) (*entities.Client, error) {
				return nil, expectedErr
			},
		}

		client, err := ds.GetClientByCpf(ctx, "")
		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, client)
	})

	t.Run("when GetClientByCpfFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		client, err := ds.GetClientByCpf(ctx, "")
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
		require.Nil(t, client)
	})
}
