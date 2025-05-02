package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/application/mock"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/ports"
	"gorm.io/gorm"

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

func TestDatastoreMock_GetClientById(t *testing.T) {
	t.Run("when GetClientByIDFn is defined, it returns the custom error", func(t *testing.T) {
		expectedErr := errors.New("db unavailable")

		ds := &mock.DatastoreMock{
			GetClientByIDFn: func(_ context.Context, _ uint) (*entities.Client, error) {
				return nil, expectedErr
			},
		}

		client, err := ds.GetClientByID(ctx, 1)
		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, client)
	})

	t.Run("when GetClientByIDfn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		client, err := ds.GetClientByID(ctx, 1)
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
		require.Nil(t, client)
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

func TestDatastoreMock_CreateOrderTx(t *testing.T) {
	t.Run("when CreateOrderTxFn is defined, it returns the custom error", func(t *testing.T) {
		expectedErr := errors.New("db unavailable")

		ds := &mock.DatastoreMock{
			CreateOrderTxFn: func(_ context.Context, _ *gorm.DB, _ *entities.Order) error {
				return expectedErr
			},
		}

		err := ds.CreateOrderTx(ctx, nil, nil)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("when CreateOrderFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		err := ds.CreateOrderTx(ctx, nil, nil)
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
	})
}

func TestDatastoreMock_UpdateOrderStatus(t *testing.T) {
	t.Run("when UpdateOrderStatusFn is defined, it returns the custom error", func(t *testing.T) {
		expectedErr := errors.New("db unavailable")

		ds := &mock.DatastoreMock{
			UpdateOrderStatusFn: func(_ context.Context, _ uint, _ string) error {
				return expectedErr
			},
		}

		err := ds.UpdateOrderStatus(ctx, 1, "pending")
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("when UpdateOrderStatusFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		err := ds.UpdateOrderStatus(ctx, 1, "pending")
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
	})
}

func TestDatastoreMock_GetPaginatedOrders(t *testing.T) {
	t.Run("when GetPaginatedOrdersFn is defined, it returns the custom error", func(t *testing.T) {
		expectedErr := errors.New("db unavailable")

		ds := &mock.DatastoreMock{
			GetPaginatedOrdersFn: func(_ context.Context, _ ports.Filter) ([]*entities.Order, error) {
				return nil, expectedErr
			},
		}

		orders, err := ds.GetPaginatedOrders(ctx, ports.Filter{})
		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, orders)
	})

	t.Run("when GetPaginatedOrdersFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		orders, err := ds.GetPaginatedOrders(ctx, ports.Filter{})
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
		require.Nil(t, orders)
	})
}

func TestDatastoreMock_GetDB(t *testing.T) {
	t.Run("returns the database instance", func(t *testing.T) {
		ds := &mock.DatastoreMock{
			GetDBFn: func() *gorm.DB {
				return nil
			},
		}

		db := ds.GetDB()
		require.Nil(t, db)
	})
}

func TestDatastoreMock_GetDBFnNotDefined(t *testing.T) {
	t.Run("when GetDBFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		db := ds.GetDB()
		require.Nil(t, db)
	})
}
