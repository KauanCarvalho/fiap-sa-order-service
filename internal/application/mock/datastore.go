package mock

import (
	"context"
	"errors"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/ports"
)

type DatastoreMock struct {
	PingFn               func(ctx context.Context) error
	CreateClientFn       func(ctx context.Context, client *entities.Client) error
	GetClientByCpfFn     func(ctx context.Context, cpf string) (*entities.Client, error)
	GetClientByIDFn      func(ctx context.Context, id uint) (*entities.Client, error)
	CreateOrderFn        func(ctx context.Context, order *entities.Order) error
	UpdateOrderStatusFn  func(ctx context.Context, orderID uint, status string) error
	GetPaginatedOrdersFn func(ctx context.Context, filter ports.Filter) ([]*entities.Order, error)
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

func (m *DatastoreMock) GetClientByID(ctx context.Context, id uint) (*entities.Client, error) {
	if m.GetClientByIDFn != nil {
		return m.GetClientByIDFn(ctx, id)
	}

	return nil, ErrFunctionNotImplemented
}

func (m *DatastoreMock) CreateOrder(ctx context.Context, order *entities.Order) error {
	if m.CreateOrderFn != nil {
		return m.CreateOrderFn(ctx, order)
	}

	return ErrFunctionNotImplemented
}

func (m *DatastoreMock) UpdateOrderStatus(ctx context.Context, orderID uint, status string) error {
	if m.UpdateOrderStatusFn != nil {
		return m.UpdateOrderStatusFn(ctx, orderID, status)
	}

	return ErrFunctionNotImplemented
}

func (m *DatastoreMock) GetPaginatedOrders(ctx context.Context, filter ports.Filter) ([]*entities.Order, error) {
	if m.GetPaginatedOrdersFn != nil {
		return m.GetPaginatedOrdersFn(ctx, filter)
	}

	return nil, ErrFunctionNotImplemented
}
