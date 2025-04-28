package mappers

import (
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
)

func ToOrderDTO(order entities.Order) dto.OrderOutput {
	orderItemDTOs := make([]dto.OrderItemOutput, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItemDTOs[i] = ToOrderItemDTO(item)
	}

	return dto.OrderOutput{
		ID:        order.ID,
		ClientID:  order.ClientID,
		Status:    order.Status,
		Price:     order.Price,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		Items:     orderItemDTOs,
	}
}
