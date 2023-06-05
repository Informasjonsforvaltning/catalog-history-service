package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	orgAdminAuth := OrgAdminAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	body, _ := json.Marshal(toBeCreated)
	req, _ := http.NewRequest("POST", "/111222333/123456789/updates", bytes.NewBuffer(body))
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	wGet := httptest.NewRecorder()
	reqGet, _ := http.NewRequest("GET", w.Header().Get("Location"), nil)
	reqGet.Header.Set("Authorization", *jwt)
	router.ServeHTTP(wGet, reqGet)
	assert.Equal(t, http.StatusOK, wGet.Code)

	var created model.Update
	err := json.Unmarshal(wGet.Body.Bytes(), &created)
	assert.Nil(t, err)
	assert.Equal(t, "123", created.Person.ID)
	assert.Equal(t, 2, len(created.Operations))
}

func TestCreateUnauthorizedWhenMissingAuthHeader(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	toBeCreated := model.UpdatePayload{}

	body, _ := json.Marshal(toBeCreated)
	req, _ := http.NewRequest("POST", "/111222333/123456789/updates", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateForbiddenForReadRole(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	toBeCreated := model.UpdatePayload{}

	body, _ := json.Marshal(toBeCreated)
	orgReadAuth := OrgReadAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)
	req, _ := http.NewRequest("POST", "/111222333/123456789/updates", bytes.NewBuffer(body))
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreateForbiddenForRoleInWrongCatalog(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	toBeCreated := model.UpdatePayload{}

	body, _ := json.Marshal(toBeCreated)
	orgAdminAuth := OrgAdminAuth("333222111")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	req, _ := http.NewRequest("POST", "/111222333/123456789/updates", bytes.NewBuffer(body))
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
