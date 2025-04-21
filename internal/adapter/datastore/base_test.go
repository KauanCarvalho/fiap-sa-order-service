package datastore_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/di"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/domain"

	"gorm.io/gorm"
)

var (
	ctx   context.Context
	cfg   *config.Config
	sqlDB *gorm.DB
	ds    domain.Datastore
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	cfg = config.Load()

	var err error
	sqlDB, err = di.NewDatabaseConnectionPool(cfg)
	if err != nil {
		log.Fatalf("error when creating database connection pool: %v", err)
	}

	ds = datastore.NewDatastore(sqlDB)

	os.Exit(m.Run())
}
