package repository

import (
	"context"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateRepository interface {
	StoreUpdate(ctx context.Context, update model.Update) error
}

type UpdateRepositoryImpl struct {
	collection *mongo.Collection
}

var updateRepository *UpdateRepositoryImpl

func InitUpdateRepository() *UpdateRepositoryImpl {
	if updateRepository == nil {
		updateRepository = &UpdateRepositoryImpl{collection: connection.MongoCollection()}
	}
	return updateRepository
}

func (r *UpdateRepositoryImpl) StoreUpdate(ctx context.Context, update model.Update) error {
	_, err := r.collection.InsertOne(ctx, update)
	return err
}
