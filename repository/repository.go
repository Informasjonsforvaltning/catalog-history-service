package repository

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
)

// JsonPatchRepository is a struct that holds a reference to a MongoDB collection
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
