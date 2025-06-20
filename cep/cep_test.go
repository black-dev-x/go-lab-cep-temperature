package cep

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/black-dev-x/go-lab-cloud-run/test"
)

const validPath = "/ws/88888888/json/"
const invalidPath = "/ws/invalid/json/"
const notFoundPath = "/ws/11111111/json/"

type Interceptor struct {
	Transport http.RoundTripper
}

type CEPResponse struct {
	Localidade string `json:"localidade"`
}

func TestGetValidCEP(t *testing.T) {
	test.When(validPath).ReturnStatusCode(http.StatusOK).ReturnBody(CEPResponse{Localidade: "Golândia"}).Execute()

	cep, err := Get("88888888")
	fmt.Printf("CEP: %+v, Error: %v\n", cep, err)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cep.Localidade != "Golândia" {
		t.Errorf("Expected localidade to be 'Golândia', got '%s'", cep.Localidade)
	}
}

func TestGetInvalidCep(t *testing.T) {
	test.When(invalidPath).ReturnStatusCode(http.StatusBadRequest).Execute()

	cep, err := Get("invalid")
	fmt.Printf("CEP: %+v, Error: %v\n", cep, err)
	if err == nil || err.Error() != Invalid {
		t.Fatalf("Expected error '%s', got '%v'", Invalid, err)
	}
	if cep != nil {
		t.Errorf("Expected CEP to be nil, got %+v", cep)
	}
}
