package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/application/mock"

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
