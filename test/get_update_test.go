package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/config"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/stretchr/testify/assert"
)

func TestGetUpdate(t *testing.T) {
	router := config.SetupRouter()

	orgReadAuth := OrgAdminAuth("123456789")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/123456789/112/updates/113", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	expectedResponse := model.Update{
		ID:         "113",
		CatalogId:  "123456789",
		ResourceId: "112",
		Person: model.Person{
			ID:    "110",
			Email: "example@example.com",
			Name:  "Doe Doe",
		},
		DateTime: time.Date(2019, 1, 4, 0, 0, 0, 0, time.UTC),
		Operations: []model.JsonPatchOperation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: "Bob",
			},
		},
	}

	var actualResponse model.Update
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestGetUpdateUnauthorizedWhenMissingAuthHeader(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/123456789/112/updates/113", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUpdateForbiddenForRoleInWrongCatalog(t *testing.T) {
	router := config.SetupRouter()

	orgAdminAuth := OrgAdminAuth("333222111")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/123456789/112/updates/113", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetUpdateNotFoundForWrongResourceId(t *testing.T) {
	router := config.SetupRouter()

	orgAdminAuth := OrgAdminAuth("123456789")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/123456789/123/updates/113", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUpdateNotFoundForWrongCatalogId(t *testing.T) {
	router := config.SetupRouter()

	orgAdminAuth := OrgAdminAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/112/updates/113", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUpdateRejectsInvalidCatalogIdWithDollarSign(t *testing.T) {
	router := config.SetupRouter()

	orgAdminAuth := OrgAdminAuth("123$456")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/123$456/112/updates/113", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	// Should return 400 (bad request) due to validation failure in service layer
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUpdateRejectsInvalidResourceIdWithCurlyBraces(t *testing.T) {
	router := config.SetupRouter()

	orgAdminAuth := OrgAdminAuth("123456789")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/123456789/{112}/updates/113", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	// Should return 400 (bad request) due to validation failure in service layer
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUpdateRejectsInvalidUpdateIdWithBrackets(t *testing.T) {
	router := config.SetupRouter()

	orgAdminAuth := OrgAdminAuth("123456789")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/123456789/112/updates/[113]", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	// Should return 400 (bad request) due to validation failure in service layer
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUpdateRejectsEmptyCatalogId(t *testing.T) {
	router := config.SetupRouter()

	orgAdminAuth := OrgAdminAuth("123456789")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "//112/updates/113", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	// Should return 404 (not found) as empty path parameter is handled by router
	// or 400 if validation catches it
	assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusBadRequest)
}
