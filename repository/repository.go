package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
)

type DataSourceRepository struct {
	collection *mongo.Collection
}

var dataSourceRepository *DataSourceRepository

func InitRepository() *DataSourceRepository {
	if dataSourceRepository == nil {
		dataSourceRepository = &DataSourceRepository{collection: connection.MongoCollection()}
	}
	return dataSourceRepository
}

func (r *DataSourceRepository) GetAllDataSources(ctx context.Context) ([]model.DataSource, error) {
	current, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer current.Close(ctx)
	var dataSources []model.DataSource
	for current.Next(ctx) {
		var dataSource model.DataSource
		err := bson.Unmarshal(current.Current, &dataSource)
		if err != nil {
			return nil, err
		}
		dataSources = append(dataSources, dataSource)
	}
	if err := current.Err(); err != nil {
		return nil, err
	}
	return dataSources, nil
}
