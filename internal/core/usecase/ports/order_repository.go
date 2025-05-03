package ports

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"

	"gorm.io/gorm"
)

type Filter struct {
	Limit  int
	Offset int
}

type OrderRepository interface {
	CreateOrderTx(ctx context.Context, tx *gorm.DB, order *entities.Order) error
	UpdateOrderStatus(ctx context.Context, orderID uint, status string) error
	GetPaginatedOrders(ctx context.Context, filter Filter) ([]*entities.Order, error)
}
