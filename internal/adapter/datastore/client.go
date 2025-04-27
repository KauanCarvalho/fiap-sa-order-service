package datastore

import (
	"context"
	"errors"
	"strings"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	internalErrors "github.com/KauanCarvalho/fiap-sa-order-service/internal/shared/errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func (d *datastore) CreateClient(ctx context.Context, client *entities.Client) error {
	if err := d.db.WithContext(ctx).Create(client).Error; err != nil {
		if isDuplicateCPFError(err) {
			return ErrExistingRecord
		}
		return internalErrors.NewInternalError("failed to create client", err)
	}
	return nil
}

func (d *datastore) GetClientByID(ctx context.Context, id uint) (*entities.Client, error) {
	var client entities.Client

	if err := d.db.WithContext(ctx).Where("id = ?", id).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, internalErrors.NewInternalError("failed to get client by ID", err)
	}
	return &client, nil
}

func (d *datastore) GetClientByCpf(ctx context.Context, cpf string) (*entities.Client, error) {
	var client entities.Client

	if err := d.db.WithContext(ctx).Where("cpf = ?", cpf).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, internalErrors.NewInternalError("failed to get client by CPF", err)
	}
	return &client, nil
}

func isDuplicateCPFError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "cpf")
	}
	return false
}
