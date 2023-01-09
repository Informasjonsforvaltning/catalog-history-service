package service

import (
	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"go.mongodb.org/mongo-driver/mongo"
)

// JsonPatchService is a struct that holds a reference to a MongoDB collection
type JsonPatchService struct {
	collection *mongo.Collection
}

// InitService creates a new JsonPatchService and returns a pointer to it
func InitService() *JsonPatchService {
	service := JsonPatchService{collection: connection.MongoCollection()}
	return &service
}
