package ports

import (
	"context"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
)

type ClientRepository interface {
	CreateClient(ctx context.Context, client *entities.Client) error
	GetClientByCpf(ctx context.Context, cpf string) (*entities.Client, error)
	GetClientByID(ctx context.Context, id uint) (*entities.Client, error)
}
