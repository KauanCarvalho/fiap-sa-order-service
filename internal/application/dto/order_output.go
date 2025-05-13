package dto

import "time"

type PaymentOutput struct {
	Status        string `json:"status,omitempty"`
	QRCode        string `json:"qr_code,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
}

type OrderOutput struct {
	ID        uint              `json:"id"`
	ClientID  uint              `json:"client_id"`
	Status    string            `json:"status"`
	Price     float64           `json:"price"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Items     []OrderItemOutput `json:"items"`
	Payment   PaymentOutput     `json:"payment,omitempty"`
}
