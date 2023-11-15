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

func TestPaginationAndSorting(t *testing.T) {
	router := config.SetupRouter()

	orgReadAuth := OrgReadAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgReadAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates?page=1&size=2&sort_by=datetime&sort_order=desc", nil)
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

	// Check that updates are returned in descending order by date
	assert.True(t, actualResponse.Updates[0].DateTime.After(actualResponse.Updates[1].DateTime))
}

func TestPaginationAndSortingTwo(t *testing.T) {
	router := config.SetupRouter()

	orgWriteAuth := OrgWriteAuth("111222333")
	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &orgWriteAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates?page=1&size=2&sort_by=name&sort_order=asc", nil)
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

	// Check that updates are returned in ascending order by person name
	assert.True(t, actualResponse.Updates[0].Person.Name <= actualResponse.Updates[1].Person.Name)
}
