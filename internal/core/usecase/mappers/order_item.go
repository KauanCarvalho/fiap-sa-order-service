package mappers

import (
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
)

func ToOrderItemDTO(orderItem entities.OrderItem) dto.OrderItemOutput {
	return dto.OrderItemOutput{
		SKU:      orderItem.SKU,
		Quantity: orderItem.Quantity,
		Price:    orderItem.Price,
	}
}
