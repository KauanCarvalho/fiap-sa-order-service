package bdd_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/api"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/payment"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/product"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/di"
	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/go-testfixtures/testfixtures/v3"
	"go.uber.org/zap"
)

var (
	engine   *gin.Engine
	cfg      *config.Config
	recorder *httptest.ResponseRecorder
	request  *http.Request
	sqlDB    *gorm.DB
	ds       domain.Datastore
	cc       usecase.CreateClientUseCase
	gc       usecase.GetClientUseCase
	co       usecase.CreateOrderUseCase
	uo       usecase.UpdateOrderUseCase
	gpo      usecase.GetPaginatedOrdersUseCase
	fixtures *testfixtures.Loader
	bodyData map[string]string
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	logger := zap.NewNop()
	zap.ReplaceGlobals(logger)

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

	fixtures, err = di.SetupFixtures(db, "../../../testdata/fixtures")
	if err != nil {
		log.Fatalf("error when initializing fixtures: %v", err)
	}

	productClient := product.NewClient(*cfg)
	paymentClient := payment.NewClient(*cfg)

	ds = datastore.NewDatastore(sqlDB)
	cc = usecase.NewCreateClientUseCase(ds)
	gc = usecase.NewGetClientUseCase(ds)
	co = usecase.NewCreateOrderUseCase(ds, productClient, paymentClient)
	uo = usecase.NewUpdateOrderUseCase(ds)
	gpo = usecase.NewGetPaginatedOrdersUseCase(ds)

	engine = api.GenerateRouter(cfg, ds, cc, gc, co, uo, gpo)

	code := m.Run()

	os.Exit(code)
}

func resetState() {
	recorder = httptest.NewRecorder()
	bodyData = make(map[string]string)
}

func loadTestFixtures() {
	if err := fixtures.Load(); err != nil {
		log.Fatalf("error when loading fixtures: %v", err)
	}
}

func ResetAndLoadFixtures() {
	resetState()
	loadTestFixtures()
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name: "Features",
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			InitializeScenarioClient(sc)
			InitializeScenarioAdminOrder(sc)
			InitializeScenarioCheckout(sc)
		},
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"./features"},
		},
	}

	if suite.Run() != 0 {
		t.Fatal("tests failed")
	}
}

func setupTestRouter(productServiceURL, paymentServiceURL string) *gin.Engine {
	mockCfg := *cfg
	mockCfg.ProductServiceURL = productServiceURL
	mockCfg.PaymentServiceURL = paymentServiceURL

	mockProductClient := product.NewClient(mockCfg)
	mockPaymentClient := payment.NewClient(mockCfg)

	mockCo := usecase.NewCreateOrderUseCase(ds, mockProductClient, mockPaymentClient)

	return api.GenerateRouter(&mockCfg, ds, cc, gc, mockCo, uo, gpo)
}
