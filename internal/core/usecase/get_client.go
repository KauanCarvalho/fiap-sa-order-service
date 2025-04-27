package usecase

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
)

type GetClientUseCase interface {
	Run(ctx context.Context, cpf string) (*entities.Client, error)
}

type getClientUseCase struct {
	ds domain.Datastore
}

func NewGetClientUseCase(ds domain.Datastore) GetClientUseCase {
	return &getClientUseCase{ds: ds}
}

func (c *getClientUseCase) Run(ctx context.Context, cpf string) (*entities.Client, error) {
	return c.ds.GetClientByCpf(ctx, cpf)
}
