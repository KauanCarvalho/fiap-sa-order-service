package dto

import "github.com/go-playground/validator/v10"

type OrderInputCreate struct {
	CognitoID string                 `json:"cognito_id"`
	ClientID  uint                   `json:"client_id"`
	Items     []OrderItemInputCreate `json:"items" validate:"required,min=1,dive"` //nolint:golines // 1 line
}
type OrderItemInputCreate struct {
	SKU      string `json:"sku"      validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

func ValidateOrderCreate(input OrderInputCreate) error {
	validate := validator.New()
	return validate.Struct(input)
}
