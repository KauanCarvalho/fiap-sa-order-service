package dto

import (
	"time"
)

type ClientOutput struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	CPF       string     `json:"cpf"`
	CognitoID NullString `json:"cognito_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
