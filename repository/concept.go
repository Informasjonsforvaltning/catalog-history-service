package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/mongodb"
	"github.com/Informasjonsforvaltning/catalog-history-service/logging"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/sirupsen/logrus"
)

type ConceptsRepository interface {
	StoreConcept(ctx context.Context, update model.Update) error
	GetConceptUpdates(ctx context.Context, query bson.D) ([]model.Update, error)
	GetConceptUpdate(ctx context.Context, id string) (*model.Update, error)
}

// conceptsRepository is a struct that holds a reference to a MongoDB collection
type ConceptsRepositoryImp struct {
	collection *mongo.Collection
}

var conceptsRepository *ConceptsRepositoryImp

func InitRepository() *ConceptsRepositoryImp {
	if conceptsRepository == nil {
		conceptsRepository = &ConceptsRepositoryImp{collection: mongodb.Collection()}
	}
	return conceptsRepository
}

func (r *ConceptsRepositoryImp) StoreConcept(ctx context.Context, update model.Update) error {
	_, err := r.collection.InsertOne(ctx, update, nil)
	return err
}

func (r ConceptsRepositoryImp) GetConceptUpdates(ctx context.Context, query bson.D, page int, size int, sortBy string, sortOrder int) ([]model.Update, error) {
	skip := (page - 1) * size

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(size))

	// Sort the results by the specified field and order
	sort := bson.D{{Key: sortBy, Value: sortOrder}}
	opts.SetSort(sort)

	var updates []model.Update

	current, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		logging.LogAndPrintError(err)
		return updates, err
	}
	defer func(current *mongo.Cursor, ctx context.Context) {
		err := current.Close(ctx)
		if err != nil {
			logging.LogAndPrintError(err)
		}
	}(current, ctx)

	for current.Next(ctx) {
		var update model.Update
		err := bson.Unmarshal(current.Current, &update)
		if err != nil {
			logging.LogAndPrintError(err)
			return updates, err
		}
		updates = append(updates, update)
	}
	if err := current.Err(); err != nil {
		logging.LogAndPrintError(err)
		return updates, err
	}
	return updates, nil
}

func (r ConceptsRepositoryImp) GetConceptUpdate(ctx context.Context, conceptId string, updateId string) (*model.Update, error) {
	filter := bson.D{{Key: "_id", Value: updateId}, {Key: "resourceId", Value: conceptId}}

	bytes, err := r.collection.FindOne(ctx, filter).DecodeBytes()
	logrus.Info("Starting to get concept update from database")
	if err == mongo.ErrNoDocuments {
		logrus.Error("concept update not found in db")
		logging.LogAndPrintError(err)
		return nil, nil
	}
	if err != nil {
		logrus.Errorf("error when getting concept from db: %s", err)
		logging.LogAndPrintError(err)
		return nil, err
	}

	var update model.Update
	unmarshalError := bson.Unmarshal(bytes, &update)
	if unmarshalError != nil {
		logrus.Errorf("error when unmarshalling concept from db: %s", unmarshalError)
		logging.LogAndPrintError(unmarshalError)
		return nil, unmarshalError
	}

	return &update, nil
}
