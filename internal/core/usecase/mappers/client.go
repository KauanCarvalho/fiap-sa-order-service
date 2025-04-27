package mappers

import (
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
)

func ToClientDTO(client entities.Client) dto.ClientOutput {
	return dto.ClientOutput{
		ID:        client.ID,
		Name:      client.Name,
		CPF:       client.CPF,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}
