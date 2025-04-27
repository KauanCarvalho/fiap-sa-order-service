package validation

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func HandleValidationError(err error) []ErrorResponse {
	var errs []ErrorResponse

	if validationErrors, ok := err.(validator.ValidationErrors); ok { //nolint:errorlint // This is not necessary here.
		for _, fieldError := range validationErrors {
			errorResponse := ErrorResponse{
				Field:   fieldError.Field(),
				Message: getErrorMessage(fieldError),
			}

			errs = append(errs, errorResponse)
		}
	} else {
		basicError := ErrorResponse{
			Field:   "generic error response",
			Message: err.Error(),
		}

		errs = append(errs, basicError)
	}

	return errs
}

func getErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "gt":
		return "This field must be greater than " + fieldError.Param()
	case "max":
		return "This field must have at most " + fieldError.Param() + " elements"
	case "min":
		return "This field must have at least " + fieldError.Param() + " elements"
	default:
		return fieldError.Error()
	}
}
