package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Informasjonsforvaltning/catalog-history-service/config"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateUpdate(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	toBeCreated := model.Update{
		Person: model.Person{
			ID:    "123456789",
			Email: "example@example.com",
			Name:  "Ola Nordmann",
		},
		Operations: []model.JsonPatchOperation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: "Test Navn",
			},
			{
				Op:   "remove",
				Path: "/height",
			},
		},
	}

	body, _ := json.Marshal(toBeCreated)
	req, _ := http.NewRequest("POST", "/concept/123456789", bytes.NewReader(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var actualResponse model.Update
	json.Unmarshal(w.Body.Bytes(), &actualResponse)

	toBeCreated.ID = actualResponse.ID
	assert.Equal(t, toBeCreated, actualResponse)
}
