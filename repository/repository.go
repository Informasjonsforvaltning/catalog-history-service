package repository

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
)

type JsonPatchRepository struct {
	collection *mongo.Collection
}

var jsonPatchRepository *JsonPatchRepository

func InitRepository() *JsonPatchRepository {
	if jsonPatchRepository == nil {
		jsonPatchRepository = &JsonPatchRepository{collection: connection.MongoCollection()}
	}
	return jsonPatchRepository
}

// ApplyJSONPatch applies a JSON Patch to a mock JSON document stored in a MongoDB collection
func (r *JsonPatchRepository) ApplyJSONPatch(ctx context.Context, patch []byte) error {
	// create a mock JSON document
	book := model.Book{
		Title:  "Mock Book",
		Author: "John Doe",
		Year:   2020,
	}

	// insert the mock JSON document into the collection
	_, err := r.collection.InsertOne(ctx, book)
	if err != nil {
		return err
	}

	// apply the JSON Patch document to the mock JSON document
	var updatedBook model.Book
	err = json.Unmarshal(patch, &updatedBook)
	if err != nil {
		return err
	}

	// update the stored document in the collection with the patched document
	_, err = r.collection.UpdateOne(ctx, bson.M{"title": "Mock Book"}, bson.M{"$set": updatedBook})
	if err != nil {
		return err
	}

	return nil
}
