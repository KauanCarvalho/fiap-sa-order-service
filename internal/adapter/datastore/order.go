package datastore

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/ports"
	internalErrors "github.com/KauanCarvalho/fiap-sa-order-service/internal/shared/errors"

	"gorm.io/gorm"
)

func (d *datastore) CreateOrderTx(ctx context.Context, tx *gorm.DB, order *entities.Order) error {
	if err := tx.WithContext(ctx).Create(order).Error; err != nil {
		if isDuplicateErr(err) {
			return ErrExistingRecord
		}
		return internalErrors.NewInternalError("failed to create order", err)
	}
	return nil
}

func (d *datastore) UpdateOrderStatus(ctx context.Context, orderID uint, status string) error {
	result := d.db.WithContext(ctx).Model(&entities.Order{}).Where("id = ?", orderID).Update("status", status)
	if result.Error != nil {
		return internalErrors.NewInternalError("failed to update order status", result.Error)
	} else if result.RowsAffected == 0 {
		return ErrOrderNotFound
	}

	return nil
}

func (d *datastore) GetPaginatedOrders(ctx context.Context, filter ports.Filter) ([]*entities.Order, error) {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	var orders []*entities.Order

	err := d.db.Model(&entities.Order{}).
		WithContext(ctx).
		Preload("OrderItems").
		Where("status IN ?", []string{"ready", "preparing", "pending"}).
		Order(
			"CASE " +
				"WHEN status = 'ready' THEN 1 " +
				"WHEN status = 'preparing' THEN 2 " +
				"WHEN status = 'pending' THEN 3 " +
				"ELSE 4 END, " +
				"created_at ASC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}
