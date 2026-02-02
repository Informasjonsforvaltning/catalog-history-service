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

func TestGetConceptUpdates(t *testing.T) {
	router := config.SetupRouter()

	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &TestValues.SysAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/concepts/updates", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var actualResponse model.Updates
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)

	assert.Nil(t, err)
	assert.True(t, len(actualResponse.Updates) > 0)
	assert.Equal(t, 1, actualResponse.Pagination.TotalPages)
	assert.Equal(t, 1, actualResponse.Pagination.Page)
	assert.Equal(t, 10, actualResponse.Pagination.Size)
}

func TestGetConceptUpdatesWithPagination(t *testing.T) {
	router := config.SetupRouter()

	orgReadAuth := OrgReadAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/concepts/updates?page=1&size=2", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var actualResponse model.Updates
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(actualResponse.Updates))
	// TotalPages depends on test order (TestCreateUpdate may add documents)
	assert.True(t, actualResponse.Pagination.TotalPages >= 2)
	assert.Equal(t, 1, actualResponse.Pagination.Page)
	assert.Equal(t, 2, actualResponse.Pagination.Size)
}

func TestGetConceptUpdatesUnauthorizedWhenMissingAuthHeader(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/concepts/updates", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetConceptUpdatesForbiddenForRoleInWrongCatalog(t *testing.T) {
	router := config.SetupRouter()

	orgAdminAuth := OrgAdminAuth("333222111")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/concepts/updates", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetConceptUpdatesRejectsInvalidCatalogId(t *testing.T) {
	router := config.SetupRouter()

	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &TestValues.SysAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/invalid$id/concepts/updates", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetConceptUpdatesRejectsExcessivePageSize(t *testing.T) {
	router := config.SetupRouter()

	orgReadAuth := OrgReadAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/concepts/updates?page=1&size=1000", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetConceptUpdatesRejectsExcessivePageNumber(t *testing.T) {
	router := config.SetupRouter()

	orgReadAuth := OrgReadAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/concepts/updates?page=10001&size=10", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetConceptUpdatesEmptyResultsForUnknownCatalog(t *testing.T) {
	router := config.SetupRouter()

	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &TestValues.SysAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/999888777/concepts/updates", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var actualResponse model.Updates
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(actualResponse.Updates))
	assert.Equal(t, 0, actualResponse.Pagination.TotalPages)
}

func TestGetConceptUpdatesAcceptsValidSortFields(t *testing.T) {
	router := config.SetupRouter()

	orgReadAuth := OrgReadAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)

	validSortFields := []string{"datetime", "name", "email"}

	for _, sortField := range validSortFields {
		t.Run("sort_by_"+sortField, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/111222333/concepts/updates?sort_by="+sortField, nil)
			req.Header.Set("Authorization", *jwt)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
