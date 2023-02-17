package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Informasjonsforvaltning/catalog-history-service/config"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/stretchr/testify/assert"
)

func TestGetConcepts(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/concepts/123456789", nil)
	router.ServeHTTP(w, req)
req, _ := http.NewRequest("GET", "/concepts/123456789/updates", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var actualResponse []model.Update
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)

	assert.Nil(t, err)
	assert.True(t, len(actualResponse) > 0)
}
