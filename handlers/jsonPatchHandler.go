package handlers

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// define a struct to hold the mock JSON document
type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// ApplyJSONPatch applies a JSON Patch to a mock JSON document stored in a MongoDB collection
func ApplyJSONPatch(ctx context.Context, client *mongo.Client) error {
	// create a collection to store the mock JSON document
	collection := client.Database("test").Collection("books")

	// create a mock JSON document
	book := Book{
		Title:  "Mock Book",
		Author: "John Doe",
		Year:   2020,
	}

	// insert the mock JSON document into the collection
	_, err := collection.InsertOne(ctx, book)
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}

	// apply the JSON Patch document to the mock JSON document
	var updatedBook Book
	err = json.Unmarshal(patch, &updatedBook)
	if err != nil {
		return err
	}

	// update the stored document in the collection with the patched document
	_, err = collection.UpdateOne(ctx, bson.M{"title": "Mock Book"}, bson.M{"$set": updatedBook})
	if err != nil {
		return err
	}

	return nil
}
