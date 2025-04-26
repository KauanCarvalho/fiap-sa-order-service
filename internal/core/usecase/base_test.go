package usecase_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/di"

	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/gorm"
)

var (
	ctx      context.Context
	cfg      *config.Config
	sqlDB    *gorm.DB
	fixtures *testfixtures.Loader
	ds       domain.Datastore
	cc       usecase.CreateClientUseCase
	gc       usecase.GetClientUseCase
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	cfg = config.Load()

	var err error
	sqlDB, err = di.NewDatabaseConnectionPool(cfg)
	if err != nil {
		log.Fatalf("error when creating database connection pool: %v", err)
	}

	db, dbErr := sqlDB.DB()
	if dbErr != nil {
		log.Fatalf("error when getting database connection: %v", dbErr)
	}

	fixtures, err = di.SetupFixtures(db, "../../../testdata/fixtures")
	if err != nil {
		log.Fatalf("error when initializing fixtures: %v", err)
	}

	ds = datastore.NewDatastore(sqlDB)
	cc = usecase.NewCreateClientUseCase(ds)
	gc = usecase.NewGetClientUseCase(ds)

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatalf("error when loading fixtures: %v", err)
	}
}
