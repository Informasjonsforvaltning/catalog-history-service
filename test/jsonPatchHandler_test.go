package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/Informasjonsforvaltning/catalog-history-service/repository"
	"github.com/Informasjonsforvaltning/catalog-history-service/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestApplyJSONPatch(t *testing.T) {
	// create a new gin router and context
	router := gin.Default()
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	// create a new service and repository
	service := service.InitService()
	repository := repository.InitRepository()

	// create a mock JSON document
	book := model.Book{
		Title:  "Mock Book",
		Author: "John Doe",
		Year:   2020,
	}

	// insert the mock JSON document into the collection
	_, err := repository.collection.InsertOne(context.Background(), book)
	assert.NoError(t, err)

	// create a JSON Patch document to update the "title" field of the mock JSON document
	patch, err := json.Marshal([]struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}{{
		Op:    "replace",
		Path:  "/title",
		Value: "Updated Title",
	}})
	assert.NoError(t, err)

	// create a request to apply the JSON Patch
	req, err := http.NewRequest("PATCH", "/", strings.NewReader(string(patch)))
	assert.NoError(t, err)

	// apply the JSON Patch
	router.PATCH("/", service.ApplyJSONPatch)
	router.ServeHTTP(ctx.Writer, req)

	// check that the status code is 200
	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	// check that the mock JSON document has been updated
	var updatedBook model.Book
	err = repository.collection.FindOne(context.Background(), bson.M{"title": "Updated Title"}).Decode(&updatedBook)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedBook.Title)
}
