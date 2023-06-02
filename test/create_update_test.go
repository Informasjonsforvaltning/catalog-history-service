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

	wGet := httptest.NewRecorder()
	reqGet, _ := http.NewRequest("GET", w.Header().Get("Location"), nil)
	router.ServeHTTP(wGet, reqGet)
	assert.Equal(t, http.StatusOK, wGet.Code)

	var created model.Update
	err := json.Unmarshal(wGet.Body.Bytes(), &created)
	assert.Nil(t, err)
	assert.Equal(t, "123", created.Person.ID)
	assert.Equal(t, 2, len(created.Operations))
}
