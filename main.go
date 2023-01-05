package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Informasjonsforvaltning/catalog-history-service/config"
	jsonpatch "github.com/evanphx/json-patch"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Begrep struct {
	ID   string `bson:"_id"`
	Term string `bson:"term"`
	Def  string `bson:"def"`
}

func main() {
	config.LoggerSetup()

	router := config.SetupRouter()
	router.Run(":8080")
	router.Run(":9091")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	begreperCollection := client.Database("mydb").Collection("begreper")

	// Find a specific document in the "begreper" collection
	var begrep Begrep
	err = begreperCollection.FindOne(context.TODO(), bson.M{"_id": "123"}).Decode(&begrep)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the Begreper struct to a JSON object
	jsonObject, err := bson.MarshalExtJSON(begrep, false, false)
	if err != nil {
		log.Fatal(err)
	}

	// Apply a JSON patch to the JSON object
	patch, err := jsonpatch.DecodePatch([]byte(`[{"op": "replace", "path": "/def", "value": "new def"}]`))
	if err != nil {
		log.Fatal(err)
	}
	modified, err := patch.Apply(jsonObject)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(modified))

	// Update the document in the "begreper" collection with the modified JSON object
	_, err = begreperCollection.UpdateOne(context.TODO(), bson.M{"_id": "123"}, bson.M{"$set": modified})
	if err != nil {
		log.Fatal(err)
	}
}
