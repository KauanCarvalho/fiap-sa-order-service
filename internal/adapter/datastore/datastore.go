package datastore

import (
	"time"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/domain"

	"gorm.io/gorm"
)

const DefaultConnectionTimeout = 5 * time.Second

type datastore struct {
	db *gorm.DB
}

func NewDatastore(db *gorm.DB) domain.Datastore {
	return &datastore{db: db}
}
