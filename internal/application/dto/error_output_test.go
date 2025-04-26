package dto_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/shared/validation"

	"github.com/stretchr/testify/assert"
)

func TestSimpleAPIErrorsOutput(t *testing.T) {
	details := "Error details"
	field := "field"
	message := "Error message"

	t.Run("should return APIErrorsOutput with one error", func(t *testing.T) {
		result := dto.SimpleAPIErrorsOutput(details, field, message)

		t.Run("should contain exactly one error", func(t *testing.T) {
			assert.Len(t, result.Errors, 1, "Must have only one error")
		})

		t.Run("should match error details", func(t *testing.T) {
			assert.Equal(t, details, result.Errors[0].Details, "Details not match")
		})

		t.Run("should match error field", func(t *testing.T) {
			assert.Equal(t, field, result.Errors[0].Field, "Field not match")
		})

		t.Run("should match error message", func(t *testing.T) {
			assert.Equal(t, message, result.Errors[0].Message, "Message not match")
		})
	})
}

func TestErrorsFromValidationErrors(t *testing.T) {
	validationErrors := []validation.ErrorResponse{
		{Field: "name", Message: "Name is required"},
		{Field: "cpf", Message: "CPF is invalid"},
	}

	t.Run("should map validation errors to APIErrorsOutput", func(t *testing.T) {
		result := dto.ErrorsFromValidationErrors(validationErrors)

		t.Run("should contain the correct number of errors", func(t *testing.T) {
			assert.Len(t, result.Errors, len(validationErrors), "Number of errors should match")
		})

		t.Run("should map field correctly", func(t *testing.T) {
			assert.Equal(t, "name", result.Errors[0].Field, "Field name doesn't match")
			assert.Equal(t, "cpf", result.Errors[1].Field, "Field cpf doesn't match")
		})

		t.Run("should map message correctly", func(t *testing.T) {
			assert.Equal(t, "Name is required", result.Errors[0].Message, "Message doesn't match for name")
			assert.Equal(t, "CPF is invalid", result.Errors[1].Message, "Message doesn't match for cpf")
		})
	})
}
