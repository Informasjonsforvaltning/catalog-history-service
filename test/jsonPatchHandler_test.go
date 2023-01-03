package test

import (
	"context"
	"testing"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// import the Book struct from jsonPatchHandler.go
type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

func TestApplyJSONPatch(t *testing.T) {
	// create a client to connect to the MongoDB server
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://172.19.0.2:27017"))
	if err != nil {
		t.Fatal(err)
	}

	// create a context for the database operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// call the ApplyJSONPatch function
	err = handlers.ApplyJSONPatch(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	// retrieve the updated document from the collection
	collection := client.Database("test").Collection("books")
	var result Book
	err = collection.FindOne(ctx, bson.M{"title": "Updated Title"}).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	// assert that the "title" field has been updated
	if result.Title != "Updated Title" {
		t.Errorf("Expected title to be 'Updated Title', got '%s'", result.Title)
	}

	// disconnect the client
	client.Disconnect(ctx)
}
