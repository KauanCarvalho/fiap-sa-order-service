package bdd_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/cucumber/godog"
)

var orderID string

func givenAnExistingOrderWithStatus(status string) error {
	var order entities.Order
	err := sqlDB.Where("status = ?", status).First(&order).Error
	if err != nil {
		return fmt.Errorf("failed to find order with status %q: %w", status, err)
	}

	orderID = strconv.FormatUint(uint64(order.ID), 10)
	return nil
}

func iPATCHToReady() error {
	endpoint := fmt.Sprintf("/api/v1/admin/orders/%s/ready", orderID)
	req, err := http.NewRequest(http.MethodPatch, endpoint, nil)
	if err != nil {
		return err
	}
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	return nil
}

func iPATCHToInvalidStatus() error {
	endpoint := fmt.Sprintf("/api/v1/admin/orders/%s/invalid-status", orderID)
	req, err := http.NewRequest(http.MethodPatch, endpoint, nil)
	if err != nil {
		return err
	}
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	return nil
}

func iPATCHToNonExistingOrder() error {
	endpoint := "/api/v1/admin/orders/99999/ready"
	req, err := http.NewRequest(http.MethodPatch, endpoint, nil)
	if err != nil {
		return err
	}
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	return nil
}

func iPATCHToInvalidID() error {
	endpoint := "/api/v1/admin/orders/invalid-id/ready"
	req, err := http.NewRequest(http.MethodPatch, endpoint, nil)
	if err != nil {
		return err
	}
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	return nil
}

func theOrderStatusInDatabaseShouldBe(expected string) error {
	var order entities.Order
	err := sqlDB.First(&order, orderID).Error
	if err != nil {
		return fmt.Errorf("failed to fetch order from database: %w", err)
	}
	if order.Status != expected {
		return fmt.Errorf("expected order status %q, got %q", expected, order.Status)
	}
	return nil
}

func InitializeScenarioAdminOrder(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ResetAndLoadFixtures()
		return ctx, nil
	})

	ctx.Step(`^an existing order with status "([^"]*)"$`, givenAnExistingOrderWithStatus)
	ctx.Step(`^I PATCH "/api/v1/admin/orders/\{orderID\}/ready"$`, iPATCHToReady)
	ctx.Step(`^I PATCH "/api/v1/admin/orders/invalid-id/ready"$`, iPATCHToInvalidID)
	ctx.Step(`^I PATCH "/api/v1/admin/orders/\{orderID\}/invalid-status"$`, iPATCHToInvalidStatus)
	ctx.Step(`^I PATCH "/api/v1/admin/orders/99999/ready"$`, iPATCHToNonExistingOrder)
	ctx.Step(`^I PATCH "/api/v1/admin/orders/1/invalid-status"$`, iPATCHToInvalidStatus)
	ctx.Step(`^the order status in the database should be "([^"]*)"$`, theOrderStatusInDatabaseShouldBe)
	ctx.Step(`^the response status should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response should contain "([^"]*)"$`, theResponseShouldContain)
}
