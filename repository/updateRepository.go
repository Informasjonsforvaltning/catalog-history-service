package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/mongodb"
	"github.com/Informasjonsforvaltning/catalog-history-service/logging"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/sirupsen/logrus"
)

type UpdateRepository interface {
	StoreUpdate(ctx context.Context, update model.Update) error
	GetUpdates(ctx context.Context, query bson.D, page int, size int, sortBy string, sortOrder int) ([]model.Update, int64, error)
	GetUpdate(ctx context.Context, catalogId string, resourceId string, updateId string) (*model.Update, error)
}

type UpdateRepositoryImpl struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var updateRepository *UpdateRepositoryImpl

func InitRepository() *UpdateRepositoryImpl {
	if updateRepository == nil {
		client := mongodb.MongoClient()
		updateRepository = &UpdateRepositoryImpl{client: client, collection: mongodb.Collection(client)}
	}
	return updateRepository
}

func (r UpdateRepositoryImpl) StoreUpdate(ctx context.Context, update model.Update) error {
	return r.client.UseSession(ctx, func(sctx mongo.SessionContext) error {
		err := sctx.StartTransaction(options.Transaction().
			SetReadConcern(readconcern.Snapshot()).
			SetWriteConcern(writeconcern.Majority()),
		)

		if err != nil {
			return err
		}

		_, err = r.collection.InsertOne(ctx, update, nil)
		if err != nil {
			sctx.AbortTransaction(sctx)
			return err
		} else {
			return nil
		}
	})
}

func (r UpdateRepositoryImpl) GetUpdates(ctx context.Context, query bson.D, page int, size int, sortBy string, sortOrder int) ([]model.Update, int64, error) {
	// Defense-in-depth: Validate pagination and sort field even though they should be validated in service layer
	// This ensures CodeQL and other static analysis tools can see the validation guard
	// Primary validation happens in service layer
	validatedPage, validatedSize, err := ValidatePagination(page, size)
	if err != nil {
		return nil, 0, err
	}
	page = validatedPage
	size = validatedSize
	
	// Validate sort field to prevent injection
	validatedSortBy := ValidateSortField(sortBy)
	
	skip := page * size
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(size))

	// Build sort using validated field name
	// The MongoDB driver safely handles this as values are treated as literals, not code
	sort := bson.D{{Key: validatedSortBy, Value: sortOrder}}
	opts.SetSort(sort)

	var updates []model.Update

	current, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		logging.LogAndPrintError(err)
		return updates, 0, err
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
			return updates, 0, err
		}
		updates = append(updates, update)
	}
	if err := current.Err(); err != nil {
		logging.LogAndPrintError(err)
		return updates, 0, err
	}

	count, err := r.collection.CountDocuments(ctx, query, options.Count())
	if err != nil {
		logging.LogAndPrintError(err)
		return nil, 0, err
	}

	return updates, count, nil
}

func (r UpdateRepositoryImpl) GetUpdate(ctx context.Context, catalogId string, resourceId string, updateId string) (*model.Update, error) {
	logrus.Info("Starting to get update from database")
	
	// Defense-in-depth: Validate input parameters even though they should be validated in service layer
	// This ensures CodeQL and other static analysis tools can see the validation guard
	// Primary validation happens in service layer
	if err := ValidateID(catalogId, "catalogId"); err != nil {
		logrus.Errorf("Invalid catalogId: %v", err)
		return nil, err
	}
	if err := ValidateID(resourceId, "resourceId"); err != nil {
		logrus.Errorf("Invalid resourceId: %v", err)
		return nil, err
	}
	if err := ValidateID(updateId, "updateId"); err != nil {
		logrus.Errorf("Invalid updateId: %v", err)
		return nil, err
	}
	
	// Build query using bson.D - values are safely escaped by the MongoDB driver
	// This prevents NoSQL injection as values are treated as literals, not code
	filter := bson.D{{Key: "_id", Value: updateId}, {Key: "catalogId", Value: catalogId}, {Key: "resourceId", Value: resourceId}}

	bytes, err := r.collection.FindOne(ctx, filter).DecodeBytes()
	if err == mongo.ErrNoDocuments {
		logrus.Error("update not found in db")
		logging.LogAndPrintError(err)
		return nil, nil
	}
	if err != nil {
		logrus.Errorf("error when getting update from db: %s", err)
		logging.LogAndPrintError(err)
		return nil, err
	}

	var update model.Update
	unmarshalError := bson.Unmarshal(bytes, &update)
	if unmarshalError != nil {
		logrus.Errorf("error when unmarshalling update from db: %s", unmarshalError)
		logging.LogAndPrintError(unmarshalError)
		return nil, unmarshalError
	}

	return &update, nil
}
