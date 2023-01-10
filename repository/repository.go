package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
)

// BegreperRepository is a struct that holds a reference to a MongoDB collection
type BegreperRepository struct {
	collection *mongo.Collection
}

var begreperRepository *BegreperRepository

func InitRepository() *BegreperRepository {
	if begreperRepository == nil {
		begreperRepository = &BegreperRepository{collection: connection.MongoCollection()}
	}
	return begreperRepository
}

func (repository *BegreperRepository) InsertBegrep(ctx context.Context, begrep model.Begrep) (string, error) {
	// Insert the begrep document into the collection
	result, err := repository.collection.InsertOne(ctx, begrep)
	if err != nil {
		return "", err
	}

	// Return the inserted document's ID
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
