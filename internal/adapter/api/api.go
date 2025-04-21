package api

import (
	"fmt"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/api/handler"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/api/middleware"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/domain"
	"go.uber.org/zap"

	docs "github.com/KauanCarvalho/fiap-sa-order-service/swagger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Server struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) Run() {
	// Stores.
	ds := newStores(s.db)

	// Web server.
	r := GenerateRouter(s.cfg, ds)

	err := r.Run(fmt.Sprintf(":%s", s.cfg.Port))
	if err != nil {
		zap.L().Fatal("Failed to start server",
			zap.String("port", s.cfg.Port),
			zap.Error(err),
		)
	}
}

func newStores(db *gorm.DB) domain.Datastore {
	return datastore.NewDatastore(db)
}

func GenerateRouter(cfg *config.Config, ds domain.Datastore) *gin.Engine {
	r := gin.New()

	setupMiddlewares(r, cfg)
	registerRoutes(r, cfg, ds)

	return r
}

func setupMiddlewares(r *gin.Engine, cfg *config.Config) {
	r.Use(
		middleware.Logger(),
		middleware.Recovery(),
	)

	r.Use(requestid.New(
		requestid.WithGenerator(func() string {
			return cfg.AppName + "-" + uuid.New().String()
		}),
	))
}

func registerRoutes(r *gin.Engine, cfg *config.Config, ds domain.Datastore) {
	healthCheckHandler := handler.NewHealthCheckHandler(ds)

	r.GET("/healthcheck", healthCheckHandler.Ping)

	if cfg.IsDevelopment() {
		docs.SwaggerInfo.BasePath = ""
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
