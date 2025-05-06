package bdd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cucumber/godog"
)

func iHaveAClientWithNameAndCpf(name, cpf string) error {
	bodyData["name"] = name
	bodyData["cpf"] = cpf
	return nil
}

func iSendAPOSTRequestTo(endpoint string) error {
	bodyJSON, _ := json.Marshal(bodyData)
	var err error
	request, err = http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(recorder, request)
	return nil
}

func iSendAPOSTRequestToWithNameAndCpf(endpoint, name, cpf string) error {
	body := map[string]string{"name": name, "cpf": cpf}
	bodyJSON, _ := json.Marshal(body)
	request, _ = http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(bodyJSON))
	request.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	return nil
}

func aClientWithCpfAlreadyExists(cpf string) error {
	return iSendAPOSTRequestToWithNameAndCpf("/api/v1/clients", "Existing User", cpf)
}

func aClientWithNameAndCpfExists(name, cpf string) error {
	return iSendAPOSTRequestToWithNameAndCpf("/api/v1/clients", name, cpf)
}

func iSendAGETRequestTo(endpoint string) error {
	var err error
	request, err = http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	return nil
}

func theResponseCodeShouldBe(code int) error {
	if recorder.Code != code {
		return fmt.Errorf("expected status %d, got %d", code, recorder.Code)
	}
	return nil
}

func theResponseShouldContain(values string) error {
	bodyBytes, _ := io.ReadAll(recorder.Body)
	body := string(bodyBytes)
	for _, val := range strings.Split(values, " and ") {
		if !strings.Contains(body, val) {
			return fmt.Errorf("expected body to contain %q", val)
		}
	}
	return nil
}

func theResponseShouldContainTwoValues(val1, val2 string) error {
	bodyBytes, _ := io.ReadAll(recorder.Body)
	body := string(bodyBytes)

	if !strings.Contains(body, val1) {
		return fmt.Errorf("expected body to contain %q", val1)
	}
	if !strings.Contains(body, val2) {
		return fmt.Errorf("expected body to contain %q", val2)
	}

	return nil
}

func InitializeScenarioClient(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ResetAndLoadFixtures()
		return ctx, nil
	})

	ctx.Step(`^I have a client with name "([^"]*)" and cpf "([^"]*)"$`, iHaveAClientWithNameAndCpf)
	ctx.Step(`^I send a POST request to "([^"]*)"$`, iSendAPOSTRequestTo)
	ctx.Step(
		`^I send a POST request to "([^"]*)" with name "([^"]*)" and cpf "([^"]*)"$`,
		iSendAPOSTRequestToWithNameAndCpf,
	)
	ctx.Step(`^I send a GET request to "([^"]*)"$`, iSendAGETRequestTo)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response should contain "([^"]*)"$`, theResponseShouldContain)
	ctx.Step(`^the response should contain "([^"]*)" and "([^"]*)"$`, theResponseShouldContainTwoValues)
	ctx.Step(`^a client with cpf "([^"]*)" already exists$`, aClientWithCpfAlreadyExists)
	ctx.Step(`^a client with name "([^"]*)" and cpf "([^"]*)" exists$`, aClientWithNameAndCpfExists)
}
