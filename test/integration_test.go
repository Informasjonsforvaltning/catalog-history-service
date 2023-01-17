package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Informasjonsforvaltning/catalog-history-service/config"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/Informasjonsforvaltning/catalog-history-service/repository"
)

func TestMain(m *testing.M) {
	MongoContainerRunner(m)
}

func TestPingRoute(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReadyRoute(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ready", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateUpdate(t *testing.T) {
	// Initialize mock data
	person := model.Person{
		ID:    "123",
		Email: "example@example.com",
		Name:  "John Doe",
	}
	operations := []model.JsonPatchOperation{
		{
			Op:    "add",
			Path:  "/email",
			Value: "example@example.com",
		},
	}
	update := model.Update{
		Person:     person,
		DateTime:   time.Now(),
		Operations: operations,
	}
	// Convert mock data to JSON
	updateJSON, err := json.Marshal(update)
	if err != nil {
		t.Error(err)
	}
	// Create a new request to the createUpdate function
	req, err := http.NewRequest("POST", "/concepts/{conceptsId}", bytes.NewBuffer(updateJSON))
	if err != nil {
		t.Error(err)
	}
	// Create a new ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(repository.CreateUpdate)
	// Serve the request to the handler
	handler.ServeHTTP(rr, req)
	// Check if the status code is 201 (created)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}
