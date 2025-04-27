package usecase

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/ports"
)

type GetPaginatedOrdersUseCase interface {
	Run(ctx context.Context, filter ports.Filter) ([]*entities.Order, error)
}

type getPaginatedOrdersUseCase struct {
	ds domain.Datastore
}

func NewGetPaginatedOrdersUseCase(ds domain.Datastore) GetPaginatedOrdersUseCase {
	return &getPaginatedOrdersUseCase{ds: ds}
}

func (c *getPaginatedOrdersUseCase) Run(ctx context.Context, filter ports.Filter) ([]*entities.Order, error) {
	return c.ds.GetPaginatedOrders(ctx, filter)
}
