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
			ID:    "123",
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
	req, _ := http.NewRequest("POST", "/concepts/123456789/updates", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	//assert.NotNil(t, location)
	assert.Equal(t, http.StatusCreated, w.Code)

	var newUpdate model.Update
	json.Unmarshal(w.Body.Bytes(), &newUpdate)
	assert.NotNil(t, newUpdate)
}
