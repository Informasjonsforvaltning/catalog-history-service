package service

import (
	"context"
	"encoding/json"

	"github.com/Informasjonsforvaltning/catalog-history-service/repository"
	"github.com/gin-gonic/gin"
)

type JsonPatchService struct {
	repository *repository.JsonPatchRepository
}

func InitService() *JsonPatchService {
	service := JsonPatchService{repository.InitRepository()}
	return &service
}

// ApplyJSONPatch applies a JSON Patch to a mock JSON document stored in a MongoDB collection
func (s *JsonPatchService) ApplyJSONPatch(c *gin.Context) {
	// create a context
	ctx := context.Background()

	// read the JSON Patch document from the request body
	var patch []byte
	err := json.NewDecoder(c.Request.Body).Decode(&patch)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// apply the JSON Patch using the repository
	err = s.repository.ApplyJSONPatch(ctx, patch)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}
