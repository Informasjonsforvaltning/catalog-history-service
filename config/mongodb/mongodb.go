package mongodb

import (
	"context"
	"github.com/Informasjonsforvaltning/catalog-history-service/logging"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/env"
)

func ConnectionString() string {
	authParams := env.ConstantValues.MongoAuthParams
	dbName := env.ConstantValues.MongoDatabase
	host := env.MongoHost()
	password := env.MongoPassword()
	user := env.MongoUsername()

	connectionString := "mongodb://" + user + ":" + password + "@" + host + "/" + dbName + "?" + authParams

	return connectionString
}

func Collection() *mongo.Collection {
	mongoOptions := options.Client().ApplyURI(ConnectionString())
	client, err := mongo.Connect(context.Background(), mongoOptions)
	if err != nil {
		logrus.Error("mongo client failed")
		logging.LogAndPrintError(err)
	}
	if err != nil {
		logrus.Error("mongo client mongodb failed")
		logging.LogAndPrintError(err)
	}
	collection := client.Database(env.ConstantValues.MongoDatabase).Collection(env.ConstantValues.MongoCollection)

	return collection
}
