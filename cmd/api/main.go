package main

import (
	"log"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/api"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/di"

	"go.uber.org/zap"
)

// @title          Fiap SA Order Service
// @version        0.0.1
// @description    Rest API for Fiap SA Order Service
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host    fa-sa-order-service
// @schemes http
// @accept  json
// @produce json
// main is the entry point of the application.
func main() {
	cfg := config.Load()

	db, err := di.NewDatabaseConnectionPool(cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot initialize zap logger: %v", err)
	}
	defer logger.Sync() //nolint:errcheck // It is not necessary to check for errors at this moment.

	zap.ReplaceGlobals(logger.With(zap.String("app", cfg.AppName), zap.String("env", cfg.AppEnv)))

	apiServer := api.NewServer(cfg, db)
	apiServer.Run()
}
