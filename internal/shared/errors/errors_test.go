package errors_test

import (
	"errors"
	"testing"

	sharedErrors "github.com/KauanCarvalho/fiap-sa-order-service/internal/shared/errors"

	"github.com/stretchr/testify/assert"
)

func TestInternalError_Error(t *testing.T) {
	t.Run("returns error with message only", func(t *testing.T) {
		internalErr := sharedErrors.NewInternalError("Something went wrong", nil)

		assert.Equal(t, "Internal Error: Something went wrong", internalErr.Error())
	})

	t.Run("returns error with message and cause", func(t *testing.T) {
		causeErr := errors.New("database connection failed")
		internalErr := sharedErrors.NewInternalError("Something went wrong", causeErr)

		expectedError := "Internal Error: Something went wrong | Cause: database connection failed"
		assert.Equal(t, expectedError, internalErr.Error())
	})
}

func TestNewInternalError(t *testing.T) {
	t.Run("creates InternalError with no cause", func(t *testing.T) {
		internalErr := sharedErrors.NewInternalError("Something went wrong", nil)

		assert.Equal(t, "Something went wrong", internalErr.(*sharedErrors.InternalError).Message)
		assert.NoError(t, internalErr.(*sharedErrors.InternalError).Reason)
	})

	t.Run("creates InternalError with a cause", func(t *testing.T) {
		causeErr := errors.New("database connection failed")
		internalErr := sharedErrors.NewInternalError("Something went wrong", causeErr)

		assert.Equal(t, "Something went wrong", internalErr.(*sharedErrors.InternalError).Message)
		assert.Equal(t, "database connection failed", internalErr.(*sharedErrors.InternalError).Reason.Error())
	})
}

func TestExteralError_Error(t *testing.T) {
	t.Run("returns error with message only", func(t *testing.T) {
		externalErr := sharedErrors.NewExternalError("Something went wrong", nil)

		assert.Equal(t, "External Error: Something went wrong", externalErr.Error())
	})

	t.Run("returns error with message and cause", func(t *testing.T) {
		causeErr := errors.New("database connection failed")
		externalErr := sharedErrors.NewExternalError("Something went wrong", causeErr)

		expectedError := "External Error: Something went wrong | Cause: database connection failed"
		assert.Equal(t, expectedError, externalErr.Error())
	})
}

func TestNewExternalError(t *testing.T) {
	t.Run("creates ExternalError with no cause", func(t *testing.T) {
		externalErr := sharedErrors.NewExternalError("Something went wrong", nil)

		assert.Equal(t, "Something went wrong", externalErr.(*sharedErrors.ExternalError).Message)
		assert.NoError(t, externalErr.(*sharedErrors.ExternalError).Reason)
	})

	t.Run("creates ExternalError with a cause", func(t *testing.T) {
		causeErr := errors.New("database connection failed")
		externalErr := sharedErrors.NewExternalError("Something went wrong", causeErr)

		assert.Equal(t, "Something went wrong", externalErr.(*sharedErrors.ExternalError).Message)
		assert.Equal(t, "database connection failed", externalErr.(*sharedErrors.ExternalError).Reason.Error())
	})
}
