package entities

import (
	"time"
)

type OrderItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;notNull"`
	OrderID   uint      `gorm:"notNull"`
	SKU       string    `gorm:"type:varchar(50);notNull"`
	Quantity  int       `gorm:"notNull"`
	Price     float64   `gorm:"type:decimal(10,2);notNull"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;notNull"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;notNull"`
	Order     Order     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}
