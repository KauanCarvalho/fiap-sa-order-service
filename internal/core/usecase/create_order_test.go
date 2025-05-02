package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/payment"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/product"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type productClientStub struct {
	GetProductFunc func(ctx context.Context, sku string) (*product.Response, error)
}

func (p *productClientStub) GetProduct(ctx context.Context, sku string) (*product.Response, error) {
	return p.GetProductFunc(ctx, sku)
}

type paymentClientStub struct {
	AuthorizePaymentFunc func(ctx context.Context, amount float64, externalReference, paymentMethod string) (*payment.Response, error)
}

func (p *paymentClientStub) AuthorizePayment(
	ctx context.Context,
	amount float64,
	externalReference, paymentMethod string,
) (*payment.Response, error) {
	return p.AuthorizePaymentFunc(ctx, amount, externalReference, paymentMethod)
}

func TestCreateOrderUseCase_Run(t *testing.T) {
	prepareTestDatabase()

	t.Run("successfully creates an order", func(t *testing.T) {
		productClient := &productClientStub{
			GetProductFunc: func(_ context.Context, _ string) (*product.Response, error) {
				return &product.Response{Price: 10.0}, nil
			},
		}

		paymentClient := &paymentClientStub{
			AuthorizePaymentFunc: func(_ context.Context, _ float64, _ string, _ string) (*payment.Response, error) {
				return &payment.Response{
					Amount:            30.0,
					Status:            "pending",
					ExternalReference: "order-123",
				}, nil
			},
		}

		uc := usecase.NewCreateOrderUseCase(ds, productClient, paymentClient)

		input := dto.OrderInputCreate{
			ClientID: 1,
			Items: []dto.OrderItemInputCreate{
				{SKU: "product-1", Quantity: 2},
				{SKU: "product-2", Quantity: 1},
			},
		}

		order, err := uc.Run(ctx, input)

		require.NoError(t, err)
		require.NotNil(t, order)
		assert.Equal(t, input.ClientID, order.ClientID)
		assert.Len(t, order.OrderItems, 2)
		assert.InEpsilon(t, 30.0, order.Price, 0.01)
	})

	t.Run("error when client not found", func(t *testing.T) {
		productClient := &productClientStub{}
		paymentClient := &paymentClientStub{}

		uc := usecase.NewCreateOrderUseCase(ds, productClient, paymentClient)

		input := dto.OrderInputCreate{
			ClientID: 9999,
			Items: []dto.OrderItemInputCreate{
				{SKU: "product-1", Quantity: 1},
			},
		}

		order, err := uc.Run(ctx, input)

		require.Error(t, err)
		assert.Nil(t, order)
	})

	t.Run("error when product client returns error", func(t *testing.T) {
		productClient := &productClientStub{
			GetProductFunc: func(_ context.Context, _ string) (*product.Response, error) {
				return nil, errors.New("product not found")
			},
		}

		paymentClient := &paymentClientStub{}

		uc := usecase.NewCreateOrderUseCase(ds, productClient, paymentClient)

		input := dto.OrderInputCreate{
			ClientID: 1,
			Items: []dto.OrderItemInputCreate{
				{SKU: "invalid-sku", Quantity: 1},
			},
		}

		order, err := uc.Run(ctx, input)

		require.Error(t, err)
		assert.Nil(t, order)
	})

	t.Run("error when payment client returns error", func(t *testing.T) {
		productClient := &productClientStub{
			GetProductFunc: func(_ context.Context, _ string) (*product.Response, error) {
				return &product.Response{Price: 10.0}, nil
			},
		}

		paymentClient := &paymentClientStub{
			AuthorizePaymentFunc: func(_ context.Context, _ float64, _ string, _ string) (*payment.Response, error) {
				return nil, errors.New("payment authorization failed")
			},
		}

		uc := usecase.NewCreateOrderUseCase(ds, productClient, paymentClient)

		input := dto.OrderInputCreate{
			ClientID: 1,
			Items: []dto.OrderItemInputCreate{
				{SKU: "product-1", Quantity: 2},
			},
		}

		order, err := uc.Run(ctx, input)

		require.Error(t, err)
		assert.Nil(t, order)
	})
}
