package service

import (
	"context"
	"encoding/json"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// JsonPatchService is a struct that holds a reference to a MongoDB collection
type JsonPatchService struct {
	collection *mongo.Collection
}

// InitService creates a new JsonPatchService and returns a pointer to it
func InitService() *JsonPatchService {
	service := JsonPatchService{collection: connection.MongoCollection()}
	return &service
}

// ApplyJSONPatch applies a JSON Patch to a mock JSON document stored in a MongoDB collection
func (service *JsonPatchService) ApplyJSONPatch(ctx context.Context, client *mongo.Client) error {
	// create a JSON Patch document to update the "title" field of the mock JSON document
	patch, err := json.Marshal([]struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value string `json:"value"`
	}{{
		Op:    "replace",
		Path:  "/begreper",
		Value: "Updated Begrep",
	}})
	if err != nil {
		return err
	}

	// apply the JSON Patch document to the mock JSON document
	var updatedBegrep model.Begrep
	err = json.Unmarshal(patch, &updatedBegrep)
	if err != nil {
		return err
	}

	// update the stored document in the collection with the patched document
	_, err = service.collection.UpdateOne(ctx, bson.M{"def": "someNewDef"}, bson.M{"$set": updatedBegrep})
	if err != nil {
		return err
	}

	return nil
}
