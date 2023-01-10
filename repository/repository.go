package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
)

// conceptsRepository is a struct that holds a reference to a MongoDB collection
type ConceptsRepository struct {
	collection *mongo.Collection
}

var conceptsRepository *ConceptsRepository

func InitRepository() *ConceptsRepository {
	if conceptsRepository == nil {
		conceptsRepository = &ConceptsRepository{collection: connection.MongoCollection()}
	}
	return conceptsRepository
}

func (repository *ConceptsRepository) InsertConcept(ctx context.Context, concept model.Concept) (string, error) {
	// Insert the concept document into the collection
	result, err := repository.collection.InsertOne(ctx, concept)
	if err != nil {
		return "", err
	}

	// Return the inserted document's ID
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
