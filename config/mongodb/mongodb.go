package mongodb

import (
	"context"
	"github.com/Informasjonsforvaltning/catalog-history-service/logging"
	"strings"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/env"
)

func ConnectionString() string {
	var connectionString strings.Builder
	connectionString.WriteString("mongodb://")
	if env.MongoUsername() != "" {
		connectionString.WriteString(env.MongoUsername())
		connectionString.WriteString(":")
		connectionString.WriteString(env.MongoPassword())
		connectionString.WriteString("@")
	}
	connectionString.WriteString(env.MongoHost())
	connectionString.WriteString("/")
	connectionString.WriteString(env.ConstantValues.MongoDatabase)
	if env.MongoUsername() != "" {
		connectionString.WriteString("?authSource=")
		connectionString.WriteString(env.MongoAuthSource())
	}
	if env.MongoReplicaSet() != "" {
		if strings.Contains(connectionString.String(), "?") {
			connectionString.WriteString("&")
		} else {
			connectionString.WriteString("?")
		}
		connectionString.WriteString("replicaSet=")
		connectionString.WriteString(env.MongoReplicaSet())
	}

	return connectionString.String()
}

func MongoClient() *mongo.Client {
	mongoOptions := options.Client().ApplyURI(ConnectionString())
	client, err := mongo.Connect(context.Background(), mongoOptions)
	if err != nil {
		logrus.Error("mongo client failed")
		logging.LogAndPrintError(err)
	}
	return client
}

func Collection(client *mongo.Client) *mongo.Collection {
	collection := client.Database(env.ConstantValues.MongoDatabase).Collection(env.ConstantValues.MongoCollection)

	return collection
}
