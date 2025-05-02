package domain

import (
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/ports"

	"gorm.io/gorm"
)

type Datastore interface {
	ports.HealthCheckRepository
	ports.ClientRepository
	ports.OrderRepository

	GetDB() *gorm.DB
}
