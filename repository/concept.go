package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
)

type ConceptsRepository interface {
	StoreConcept(ctx context.Context, concept model.Concept) error
}

// conceptsRepository is a struct that holds a reference to a MongoDB collection
type ConceptsRepositoryImp struct {
	collection *mongo.Collection
}

var conceptsRepository *ConceptsRepositoryImp

func InitRepository() *ConceptsRepositoryImp {
	if conceptsRepository == nil {
		conceptsRepository = &ConceptsRepositoryImp{collection: connection.MongoCollection()}
	}
	return conceptsRepository
}

func (r *ConceptsRepositoryImp) StoreConcept(ctx context.Context, update model.Update) error {
	_, err := r.collection.InsertOne(ctx, update, nil)
	return err
}
