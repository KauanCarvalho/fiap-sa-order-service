package handler_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/api"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/product"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/di"
	"github.com/go-testfixtures/testfixtures/v3"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ctx       context.Context
	cfg       *config.Config
	sqlDB     *gorm.DB
	fixtures  *testfixtures.Loader
	ds        domain.Datastore
	cc        usecase.CreateClientUseCase
	gc        usecase.GetClientUseCase
	co        usecase.CreateOrderUseCase
	uo        usecase.UpdateOrderUseCase
	gpo       usecase.GetPaginatedOrdersUseCase
	ginEngine *gin.Engine
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	logger := zap.NewNop()
	zap.ReplaceGlobals(logger)

	ctx = context.Background()
	cfg = config.Load()

	var err error
	sqlDB, err = di.NewDatabaseConnectionPool(cfg)
	if err != nil {
		log.Fatalf("error when initializing database connection: %v", err)
	}

	db, dbErr := sqlDB.DB()
	if dbErr != nil {
		log.Fatalf("error when getting database connection: %v", dbErr)
	}

	fixtures, err = di.SetupFixtures(db, "../../../../testdata/fixtures")
	if err != nil {
		log.Fatalf("error when initializing fixtures: %v", err)
	}

	productClient := product.NewClient(*cfg)

	ds = datastore.NewDatastore(sqlDB)
	cc = usecase.NewCreateClientUseCase(ds)
	gc = usecase.NewGetClientUseCase(ds)
	co = usecase.NewCreateOrderUseCase(ds, productClient)
	uo = usecase.NewUpdateOrderUseCase(ds)
	gpo = usecase.NewGetPaginatedOrdersUseCase(ds)
	ginEngine = api.GenerateRouter(cfg, ds, cc, gc, co, uo, gpo)

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatalf("error when loading fixtures: %v", err)
	}
}

func setupTestRouter(productServiceURL string) *gin.Engine {
	mockCfg := *cfg
	mockCfg.ProductServiceURL = productServiceURL

	mockProductClient := product.NewClient(mockCfg)
	mockCo := usecase.NewCreateOrderUseCase(ds, mockProductClient)

	return api.GenerateRouter(&mockCfg, ds, cc, gc, mockCo, uo, gpo)
}
