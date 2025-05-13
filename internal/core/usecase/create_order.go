package usecase

import (
	"context"
	"strconv"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/payment"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/product"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/dto"

	"gorm.io/gorm"
)

type CreateOrderUseCase interface {
	Run(ctx context.Context, input dto.OrderInputCreate) (*entities.Order, error)
}

type createOrderUseCase struct {
	ds            domain.Datastore
	productClient product.ClientInterface
	paymentClient payment.ClientInterface
}

func NewCreateOrderUseCase(
	ds domain.Datastore,
	productClient product.ClientInterface,
	paymentClient payment.ClientInterface,
) CreateOrderUseCase {
	return &createOrderUseCase{
		ds:            ds,
		productClient: productClient,
		paymentClient: paymentClient,
	}
}

func (c *createOrderUseCase) Run(ctx context.Context, input dto.OrderInputCreate) (*entities.Order, error) {
	var client *entities.Client
	var err error
	switch {
	case input.ClientID > 0:
		client, err = c.ds.GetClientByID(ctx, input.ClientID)
		if err != nil {
			return nil, err
		}
	case input.CognitoID != "":
		client, err = c.ds.GetClientByCognitoID(ctx, input.CognitoID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, gorm.ErrRecordNotFound
	}

	var createdOrder *entities.Order
	err = c.ds.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order := &entities.Order{
			ClientID: client.ID,
			Status:   "pending",
		}

		for _, item := range input.Items {
			productResponse, errClient := c.productClient.GetProduct(ctx, item.SKU)
			if errClient != nil {
				return errClient
			}
			orderItem := entities.OrderItem{
				SKU:      item.SKU,
				Quantity: item.Quantity,
				Price:    productResponse.Price,
			}
			order.OrderItems = append(order.OrderItems, orderItem)
		}

		order.CalculateTotal()

		if err = c.ds.CreateOrderTx(ctx, tx, order); err != nil {
			return err
		}

		externalRef := strconv.FormatUint(uint64(order.ID), 10)
		paymentResponse, errPayment := c.paymentClient.AuthorizePayment(ctx, order.Price, externalRef, "pix")
		if errPayment != nil {
			return errPayment
		}

		order.Payment.Status = paymentResponse.Status
		order.Payment.QRCode = paymentResponse.QRCode
		order.Payment.PaymentMethod = paymentResponse.PaymentMethod
		createdOrder = order

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}
