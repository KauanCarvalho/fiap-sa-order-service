package mock

import (
	"context"
	"errors"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
)

type DatastoreMock struct {
	PingFn           func(ctx context.Context) error
	CreateClientFn   func(ctx context.Context, client *entities.Client) error
	GetClientByCpfFn func(ctx context.Context, cpf string) (*entities.Client, error)
}

var ErrFunctionNotImplemented = errors.New("function not implemented")

func (m *DatastoreMock) Ping(ctx context.Context) error {
	if m.PingFn != nil {
		return m.PingFn(ctx)
	}

	return ErrFunctionNotImplemented
}

func (m *DatastoreMock) CreateClient(ctx context.Context, client *entities.Client) error {
	if m.CreateClientFn != nil {
		return m.CreateClientFn(ctx, client)
	}

	return ErrFunctionNotImplemented
}

func (m *DatastoreMock) GetClientByCpf(ctx context.Context, cpf string) (*entities.Client, error) {
	if m.GetClientByCpfFn != nil {
		return m.GetClientByCpfFn(ctx, cpf)
	}

	return nil, ErrFunctionNotImplemented
}
