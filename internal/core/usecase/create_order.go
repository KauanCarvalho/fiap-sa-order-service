package usecase

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/product"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/dto"
)

type CreateOrderUseCase interface {
	Run(ctx context.Context, input dto.OrderInputCreate) (*entities.Order, error)
}

type createOrderUseCase struct {
	ds            domain.Datastore
	productClient product.ClientInterface
}

func NewCreateOrderUseCase(ds domain.Datastore, productClient product.ClientInterface) CreateOrderUseCase {
	return &createOrderUseCase{
		ds:            ds,
		productClient: productClient,
	}
}

func (c *createOrderUseCase) Run(ctx context.Context, input dto.OrderInputCreate) (*entities.Order, error) {
	_, err := c.ds.GetClientByID(ctx, input.ClientID)
	if err != nil {
		return nil, err
	}

	order := &entities.Order{
		ClientID: input.ClientID,
		Status:   "pending",
	}

	for _, item := range input.Items {
		productResponse, errClient := c.productClient.GetProduct(ctx, item.SKU)
		if errClient != nil {
			return nil, errClient
		}

		orderItem := entities.OrderItem{
			SKU:      item.SKU,
			Quantity: item.Quantity,
			Price:    productResponse.Price,
		}
		order.OrderItems = append(order.OrderItems, orderItem)
	}

	order.CalculateTotal()

	err = c.ds.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
