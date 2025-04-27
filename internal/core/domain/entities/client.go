package entities

import (
	"time"
)

type Client struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	Name      string
	CPF       string    `gorm:"uniqueIndex"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
