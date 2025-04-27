package usecase

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain"
)

type UpdateOrderUseCase interface {
	Run(ctx context.Context, orderID uint, status string) error
}

type updateOrderUseCase struct {
	ds domain.Datastore
}

func NewUpdateOrderUseCase(ds domain.Datastore) UpdateOrderUseCase {
	return &updateOrderUseCase{
		ds: ds,
	}
}

func (u *updateOrderUseCase) Run(ctx context.Context, orderID uint, status string) error {
	return u.ds.UpdateOrderStatus(ctx, orderID, status)
}
