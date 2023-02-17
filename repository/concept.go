package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/sirupsen/logrus"
)

type ConceptsRepository interface {
	StoreConcept(ctx context.Context, update model.Update) error
	GetConceptUpdates(ctx context.Context, query bson.D) ([]*model.Update, error)
	GetConceptUpdate(ctx context.Context, id string) (*model.Update, error)
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

func (r ConceptsRepositoryImp) GetConceptUpdates(ctx context.Context, query bson.D) ([]*model.Update, error) {
	current, err := r.collection.Find(ctx, query)
	logrus.Info("Starting GetConceptUpdates")
	if err != nil {
		return nil, err
	}
	defer current.Close(ctx)
	var updates []*model.Update
	for current.Next(ctx) {
		var update model.Update
		err := bson.Unmarshal(current.Current, &update)
		if err != nil {
			return nil, err
		}
		updates = append(updates, &update)
	}
	if err := current.Err(); err != nil {
		return nil, err
	}
	logrus.Info("Finished getting all concept updates from database")
	return updates, nil
}

func (r ConceptsRepositoryImp) GetConceptUpdate(ctx context.Context, conceptId string, updateId string) (*model.Update, error) {
	filter := bson.D{{Key: "id", Value: updateId}, {Key: "resourceId", Value: conceptId}}
	bytes, err := r.collection.FindOne(ctx, filter).DecodeBytes()
	logrus.Info("Starting to get concept update from database")

	if err == mongo.ErrNoDocuments {
		logrus.Error("concept update not found in db")
		return nil, nil
	}
	if err != nil {
		logrus.Errorf("error when getting concept from db: %s", err)
		return nil, err
	}

	var update model.Update
	unmarshalError := bson.Unmarshal(bytes, &update)
	if unmarshalError != nil {
		logrus.Errorf("error when unmarshalling concept from db: %s", unmarshalError)
		return nil, unmarshalError
	}

	return &update, nil
}
