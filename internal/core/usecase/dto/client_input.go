package dto

import (
	"github.com/go-playground/validator"
)

type ClientInputCreate struct {
	Name      string `json:"name"       validate:"required,max=100"`
	CPF       string `json:"cpf"        validate:"required"`
	CognitoID string `json:"cognito_id"`
}

func ValidateClientCreate(input ClientInputCreate) error {
	return validator.New().Struct(input)
}
