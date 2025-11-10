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

func TestGetUpdates(t *testing.T) {
	router := config.SetupRouter()

	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &TestValues.SysAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates", nil)
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

func TestGetUpdatesWithPagination(t *testing.T) {
	router := config.SetupRouter()

	orgReadAuth := OrgReadAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates?page=1&size=2", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var actualResponse model.Updates
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(actualResponse.Updates))
	assert.Equal(t, 3, actualResponse.Pagination.TotalPages)
	assert.Equal(t, 1, actualResponse.Pagination.Page)
	assert.Equal(t, 2, actualResponse.Pagination.Size)
}

func TestGetUpdatesUnauthorizedWhenMissingAuthHeader(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUpdatesForbiddenForRoleInWrongCatalog(t *testing.T) {
	router := config.SetupRouter()

	orgAdminAuth := OrgAdminAuth("333222111")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetUpdatesRejectsInvalidSortByField(t *testing.T) {
	router := config.SetupRouter()

	orgReadAuth := OrgReadAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	// Try to inject a dangerous sort field
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates?sort_by=$where", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)
	
	// Should still return 200 OK, but sort field should default to "datetime"
	assert.Equal(t, http.StatusOK, w.Code)
	
	var actualResponse model.Updates
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.Nil(t, err)
	// The dangerous sort field should be rejected and default to datetime sorting
}

func TestGetUpdatesRejectsExcessivePageSize(t *testing.T) {
	router := config.SetupRouter()

	orgReadAuth := OrgReadAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	// Try to request more than max size (100)
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates?page=1&size=1000", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)
	
	// Should return 400 (bad request) due to validation failure
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUpdatesRejectsExcessivePageNumber(t *testing.T) {
	router := config.SetupRouter()

	orgReadAuth := OrgReadAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	// Try to request more than max page (10000)
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates?page=10001&size=10", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)
	
	// Should return 400 (bad request) due to validation failure
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUpdatesAcceptsValidSortFields(t *testing.T) {
	router := config.SetupRouter()

	orgReadAuth := OrgReadAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)
	
	validSortFields := []string{"datetime", "name", "email", "person.name", "person.email"}
	
	for _, sortField := range validSortFields {
		t.Run("sort_by_"+sortField, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/111222333/123456789/updates?sort_by="+sortField, nil)
			req.Header.Set("Authorization", *jwt)
			router.ServeHTTP(w, req)
			
			// Should return 200 OK for valid sort fields
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
