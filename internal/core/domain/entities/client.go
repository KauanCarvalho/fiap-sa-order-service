package entities

import (
	"database/sql"
	"time"
)

type Client struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	Name      string
	CognitoID sql.NullString
	CPF       string    `gorm:"uniqueIndex"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
