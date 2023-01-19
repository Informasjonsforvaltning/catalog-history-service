package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/handlers"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/stretchr/testify/assert"
)

func TestUpdateJsonPatch(t *testing.T) {
	original := model.Update{
		Person: model.Person{
			ID:    "123",
			Email: "email",
			Name:  "name",
		},
		DateTime: time.Now(),
		Operations: []model.JsonPatchOperation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: "Jane",
			},
			{
				Op:   "remove",
				Path: "/height",
			},
		},
	}
	patch := []byte(`[
		{"op": "replace", "path": "/name", "value": "John"}
	]`)

	originalBytes, _ := json.Marshal(original)
	patchedBytes, _ := jsonpatch.MergePatch(originalBytes, patch)

	// create request
	req, err := http.NewRequest("PATCH", "/", bytes.NewBuffer(patch))
	if err != nil {
		t.Fatal(err)
	}

	// create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.UpdateJsonPatch)
	handler.ServeHTTP(rr, req)

	// check response code
	assert.Equal(t, http.StatusOK, rr.Code)

	// check response body
	assert.JSONEq(t, string(patchedBytes), rr.Body.String())
}
