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
	toBeCreated := model.UpdatePayload{
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
	req, _ := http.NewRequest("POST", "/concepts/123456789", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)
	location, _ := w.Result().Location()

	//assert.NotNil(t, location)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest("GET", location.Path, nil)

	var newUpdate model.Update
	json.Unmarshal(w.Body.Bytes(), &newUpdate)
	assert.NotNil(t, newUpdate)
}
