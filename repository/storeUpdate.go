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
type UpdateService struct {
	UpdateRepository UpdateRepository
}

var updateRepository *UpdateRepositoryImpl

func InitUpdateRepository() *UpdateRepositoryImpl {
	if updateRepository == nil {
		updateRepository = &UpdateRepositoryImpl{collection: connection.MongoCollection()}
	}
	return updateRepository
}

func NewUpdateService(r UpdateRepository) *UpdateService {
	return &UpdateService{
		UpdateRepository: r,
	}
}

func (s *UpdateService) StoreUpdate(ctx context.Context, update model.Update) error {
	return s.UpdateRepository.StoreUpdate(ctx, update)
}
