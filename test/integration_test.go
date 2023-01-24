package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Informasjonsforvaltning/catalog-history-service/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllHandler(t *testing.T) {
	// Create a new Gin engine
	r := gin.Default()
	r.GET("/concepts", handlers.GetAllHandler())

	// Create a new request and response recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/concepts", nil)
	r.ServeHTTP(w, req)

	// Assert the response code and response body
	assert.Equal(t, http.StatusOK, w.Code)
	var concepts []interface{}
	json.Unmarshal(w.Body.Bytes(), &concepts)
	assert.NotEmpty(t, concepts)
}

func TestGetUpdateHandler(t *testing.T) {
	// Create a new Gin engine
	r := gin.Default()
	r.GET("/concepts/:conceptId", handlers.GetUpdateHandler())

	// Create a new request and response recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/concepts/123", nil)
	r.ServeHTTP(w, req)

	// Assert the response code and response body
	assert.Equal(t, http.StatusOK, w.Code)
	var concept interface{}
	json.Unmarshal(w.Body.Bytes(), &concept)
	assert.NotNil(t, concept)
}
