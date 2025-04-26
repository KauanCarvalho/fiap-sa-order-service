package usecase

import (
	"context"
	"errors"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/dto"
	"gorm.io/gorm"
)

type CreateClientUseCase interface {
	Run(ctx context.Context, input dto.ClientInputCreate) (*entities.Client, error)
}

type createClientsUseCase struct {
	ds domain.Datastore
}

func NewCreateClientUseCase(ds domain.Datastore) CreateClientUseCase {
	return &createClientsUseCase{ds: ds}
}

func (c *createClientsUseCase) Run(ctx context.Context, input dto.ClientInputCreate) (*entities.Client, error) {
	existentClient, err := c.ds.GetClientByCpf(ctx, input.CPF)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existentClient != nil {
		return nil, datastore.ErrExistingRecord
	}

	client := &entities.Client{
		Name: input.Name,
		CPF:  input.CPF,
	}

	err = c.ds.CreateClient(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}
