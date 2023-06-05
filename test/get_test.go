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
}

func TestGetListUnauthorizedWhenMissingAuthHeader(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetListForbiddenForRoleInWrongCatalog(t *testing.T) {
	router := config.SetupRouter()

	orgAdminAuth := OrgAdminAuth("333222111")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates", nil)
	req.Header.Set("Authorization", *jwt)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
