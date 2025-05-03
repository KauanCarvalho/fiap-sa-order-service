package entities

import (
	"time"
)

type Payment struct {
	Status        string `gorm:"-"`
	QRCode        string `gorm:"-"`
	PaymentMethod string `gorm:"-"`
}
type Order struct {
	ID         uint        `gorm:"primaryKey;autoIncrement"`
	ClientID   uint        `gorm:"not null"`
	Status     string      `gorm:"type:varchar(20);not null"`
	Price      float64     `gorm:"type:decimal(10,2);not null"`
	CreatedAt  time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Client     Client      `gorm:"foreignKey:ClientID;references:ID"`
	Payment    Payment     `gorm:"-"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}

func (o *Order) CalculateTotal() {
	var total float64
	for _, item := range o.OrderItems {
		total += float64(item.Quantity) * item.Price
	}
	o.Price = total
}
