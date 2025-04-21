package handler_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/api"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/di"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/domain"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ctx       context.Context
	cfg       *config.Config
	sqlDB     *gorm.DB
	ds        domain.Datastore
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

	ds = datastore.NewDatastore(sqlDB)

	ginEngine = api.GenerateRouter(cfg, ds)

	os.Exit(m.Run())
}
