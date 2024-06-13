package service

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/logging"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/Informasjonsforvaltning/catalog-history-service/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateService interface {
	StoreUpdate(ctx context.Context, bytes []byte, catalogId string, resourceId string)
	GetUpdates(ctx context.Context, catalogId string, resourceId string, page int, size int, sortBy string, sortOrder string)
	GetUpdate(ctx context.Context, catalogId string, resourceId string, updateId string)
}

type UpdateServiceImpl struct {
	UpdateRepository repository.UpdateRepository
}

func InitUpdateService() *UpdateServiceImpl {
	service := UpdateServiceImpl{
		UpdateRepository: *repository.InitRepository(),
	}
	return &service
}

func (service UpdateServiceImpl) StoreUpdate(ctx context.Context, bytes []byte, catalogId string, resourceId string) (string, error) {
	var update model.UpdatePayload
	err := json.Unmarshal(bytes, &update)
	logrus.Info("Unmarshalled update")
	if err != nil {
		logrus.Error("Unable to unmarshal update")
		logging.LogAndPrintError(err)
		return "", err
	}
	err = update.Validate()
	if err != nil {
		logrus.Error("update is not valid")
		logging.LogAndPrintError(err)
		return "", err
	}
	var updateDbo = model.Update{
		ID:         uuid.New().String(),
		CatalogId:  catalogId,
		ResourceId: resourceId,
		DateTime:   time.Now(),
		Person:     update.Person,
		Operations: update.Operations,
	}
	err = service.UpdateRepository.StoreUpdate(ctx, updateDbo)
	if err != nil {
		logrus.Error("Could not store update")
		logging.LogAndPrintError(err)
		return "", err
	}
	return updateDbo.ID, nil
}

func (service UpdateServiceImpl) GetUpdates(ctx context.Context, catalogId string, resourceId string, page int, size int, sortBy string, sortOrder string) (model.Updates, int) {
	query := bson.D{}
	query = append(query, bson.E{Key: "catalogId", Value: catalogId})
	query = append(query, bson.E{Key: "resourceId", Value: resourceId})

	// Set default sort by column to "datetime"
	sortByCol := "datetime"
	switch sortBy {
	case "name":
		sortByCol = "person.name"
	case "email":
		sortByCol = "person.email"
	}

	// Map sortOrder string to integer value
	sortOrderInt := -1
	if sortOrder == "asc" {
		sortOrderInt = 1
	}

	databaseUpdates, count, err := service.UpdateRepository.GetUpdates(ctx, query, page, size, sortByCol, sortOrderInt)
	if err != nil {
		logrus.Error("Get updates failed")
		logging.LogAndPrintError(err)
		return model.Updates{}, http.StatusInternalServerError
	}

	if databaseUpdates == nil {
		logrus.Error("No updates found")
		pagination := model.Pagination{TotalPages: 0, Page: page, Size: size}
		return model.Updates{Updates: []model.Update{}, Pagination: pagination}, http.StatusOK
	} else {
		logrus.Debug("Returning updates")
		totalPages := int(math.Ceil(float64(count) / float64(size)))
		pagination := model.Pagination{TotalPages: totalPages, Page: page, Size: size}
		return model.Updates{Updates: databaseUpdates, Pagination: pagination}, http.StatusOK
	}
}

func (service UpdateServiceImpl) GetUpdate(ctx context.Context, catalogId string, resourceId string, updateId string) (*model.Update, int) {
	update, err := service.UpdateRepository.GetUpdate(ctx, catalogId, resourceId, updateId)
	if err != nil {
		logrus.Error("Unable to get update")
		logging.LogAndPrintError(err)
		return nil, http.StatusInternalServerError
	} else if update == nil {
		logrus.Error("update not found")
		return nil, http.StatusNotFound
	} else {
		return update, http.StatusOK
	}
}
