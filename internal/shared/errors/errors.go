package errors

import "fmt"

type InternalError struct {
	Message string
	Reason  error
}

func (e *InternalError) Error() string {
	if e.Reason != nil {
		return fmt.Sprintf("Internal Error: %s | Cause: %s", e.Message, e.Reason.Error())
	}

	return fmt.Sprintf("Internal Error: %s", e.Message)
}

func NewInternalError(message string, reason error) error {
	return &InternalError{
		Message: message,
		Reason:  reason,
	}
}
