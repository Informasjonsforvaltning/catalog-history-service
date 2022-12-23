package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/config"
	"github.com/Informasjonsforvaltning/catalog-history-service/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config.LoggerSetup()

	router := config.SetupRouter()
	router.Run(":9091")

	// create a client to connect to the MongoDB server
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// create a context for the database operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// call the function from jsonPatchHandler.go
	err = handlers.ApplyJSONPatch(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("JSON patch applied successfully!")
}
