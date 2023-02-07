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

	assert.Equal(t, http.StatusOK, w.Code)

	var actualResponse []model.UpdateMeta
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)

	assert.Nil(t, err)
	assert.True(t, len(actualResponse) > 0)
}

// test the func compareJSON using the response from the getConcepts test
func TestCompareJSON(t *testing.T) {
	router := config.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/concepts/123456789", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var actualResponse []model.UpdateMeta
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)

	assert.Nil(t, err)
	assert.True(t, len(actualResponse) > 0)

	var expectedResponse = []model.UpdateMeta{
		{
			ID:         "123456789",
			ResourceId: "123456789",
			DateTime:   "2020-01-01T00:00:00Z",
			Person:     "123456789",
			Operations: []model.Operation{
				{
					Op:    "replace",
					Path:  "/name",
					Value: "Johnny",
				},
			},
		},
	}

	assert.True(t, compareJSON(w.Body.String(), toJSON(expectedResponse)))
}

func compareJSON(a, b string) bool {
	// return Equal([]byte(a), []byte(b))

	var objA, objB map[string]interface{}
	json.Unmarshal([]byte(a), &objA)
	json.Unmarshal([]byte(b), &objB)

	// fmt.Printf("Comparing %#v\nagainst %#v\n", objA, objB)
	return reflect.DeepEqual(objA, objB)
}

func applyPatch(doc, patch string) (string, error) {
	obj, err := DecodePatch([]byte(patch))

	if err != nil {
		panic(err)
	}

	out, err := obj.Apply([]byte(doc))

	if err != nil {
		return "", err
	}

	return string(out), nil
}
