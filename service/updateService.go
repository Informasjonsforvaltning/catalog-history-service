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
	GetConceptUpdates(ctx context.Context, catalogId string, page int, size int, sortBy string, sortOrder string)
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
	// Validate input parameters to prevent NoSQL injection
	if err := repository.ValidateID(catalogId, "catalogId"); err != nil {
		logrus.Errorf("Invalid catalogId: %v", err)
		logging.LogAndPrintError(err)
		return "", err
	}
	if err := repository.ValidateID(resourceId, "resourceId"); err != nil {
		logrus.Errorf("Invalid resourceId: %v", err)
		logging.LogAndPrintError(err)
		return "", err
	}

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
	// Validate input parameters to prevent NoSQL injection and DoS attacks
	// Validation errors are client errors, so return 400 Bad Request
	if err := repository.ValidateID(catalogId, "catalogId"); err != nil {
		logrus.Errorf("Invalid catalogId: %v", err)
		logging.LogAndPrintError(err)
		return model.Updates{}, http.StatusBadRequest
	}
	if err := repository.ValidateID(resourceId, "resourceId"); err != nil {
		logrus.Errorf("Invalid resourceId: %v", err)
		logging.LogAndPrintError(err)
		return model.Updates{}, http.StatusBadRequest
	}

	// Validate pagination parameters
	validatedPage, validatedSize, err := repository.ValidatePagination(page, size)
	if err != nil {
		logrus.Errorf("Invalid pagination parameters: %v", err)
		logging.LogAndPrintError(err)
		return model.Updates{}, http.StatusBadRequest
	}

	// Validate and whitelist sort field to prevent injection attacks
	// Only specific fields are allowed to prevent malicious field names
	sortByCol := "datetime"
	switch sortBy {
	case "name":
		sortByCol = "person.name"
	case "email":
		sortByCol = "person.email"
	case "datetime":
		sortByCol = "datetime"
	// Any other value defaults to "datetime" for safety
	}

	// Map sortOrder string to integer value
	sortOrderInt := -1
	if sortOrder == "asc" {
		sortOrderInt = 1
	}

	// Build query using bson.D - values are safely escaped by the MongoDB driver
	// This prevents NoSQL injection as values are treated as literals, not code
	query := bson.D{}
	query = append(query, bson.E{Key: "catalogId", Value: catalogId})
	query = append(query, bson.E{Key: "resourceId", Value: resourceId})

	databaseUpdates, count, err := service.UpdateRepository.GetUpdates(ctx, query, validatedPage, validatedSize, sortByCol, sortOrderInt)
	if err != nil {
		logrus.Error("Get updates failed")
		logging.LogAndPrintError(err)
		return model.Updates{}, http.StatusInternalServerError
	}

	if databaseUpdates == nil {
		logrus.Error("No updates found")
		pagination := model.Pagination{TotalPages: 0, Page: validatedPage, Size: validatedSize}
		return model.Updates{Updates: []model.Update{}, Pagination: pagination}, http.StatusOK
	} else {
		logrus.Debug("Returning updates")
		// Use validated size for accurate pagination calculation
		totalPages := int(math.Ceil(float64(count) / float64(validatedSize)))
		pagination := model.Pagination{TotalPages: totalPages, Page: validatedPage, Size: validatedSize}
		return model.Updates{Updates: databaseUpdates, Pagination: pagination}, http.StatusOK
	}
}

func (service UpdateServiceImpl) GetUpdate(ctx context.Context, catalogId string, resourceId string, updateId string) (*model.Update, int) {
	// Validate input parameters to prevent NoSQL injection
	// Validation errors are client errors, so return 400 Bad Request
	if err := repository.ValidateID(catalogId, "catalogId"); err != nil {
		logrus.Errorf("Invalid catalogId: %v", err)
		logging.LogAndPrintError(err)
		return nil, http.StatusBadRequest
	}
	if err := repository.ValidateID(resourceId, "resourceId"); err != nil {
		logrus.Errorf("Invalid resourceId: %v", err)
		logging.LogAndPrintError(err)
		return nil, http.StatusBadRequest
	}
	if err := repository.ValidateID(updateId, "updateId"); err != nil {
		logrus.Errorf("Invalid updateId: %v", err)
		logging.LogAndPrintError(err)
		return nil, http.StatusBadRequest
	}

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

func (service UpdateServiceImpl) GetConceptUpdates(ctx context.Context, catalogId string, page int, size int, sortBy string, sortOrder string) (model.Updates, int) {
	// Validate catalogId to prevent NoSQL injection
	if err := repository.ValidateID(catalogId, "catalogId"); err != nil {
		logrus.Errorf("Invalid catalogId: %v", err)
		logging.LogAndPrintError(err)
		return model.Updates{}, http.StatusBadRequest
	}

	// Validate pagination parameters
	validatedPage, validatedSize, err := repository.ValidatePagination(page, size)
	if err != nil {
		logrus.Errorf("Invalid pagination parameters: %v", err)
		logging.LogAndPrintError(err)
		return model.Updates{}, http.StatusBadRequest
	}

	// Validate and whitelist sort field to prevent injection attacks
	sortByCol := "datetime"
	switch sortBy {
	case "name":
		sortByCol = "person.name"
	case "email":
		sortByCol = "person.email"
	case "datetime":
		sortByCol = "datetime"
	}

	// Map sortOrder string to integer value
	sortOrderInt := -1
	if sortOrder == "asc" {
		sortOrderInt = 1
	}

	// Build query with catalogId only (no resourceId filter)
	query := bson.D{{Key: "catalogId", Value: catalogId}}

	databaseUpdates, count, err := service.UpdateRepository.GetUpdates(ctx, query, validatedPage, validatedSize, sortByCol, sortOrderInt)
	if err != nil {
		logrus.Error("Get concept updates failed")
		logging.LogAndPrintError(err)
		return model.Updates{}, http.StatusInternalServerError
	}

	if databaseUpdates == nil {
		logrus.Debug("No concept updates found")
		pagination := model.Pagination{TotalPages: 0, Page: validatedPage, Size: validatedSize}
		return model.Updates{Updates: []model.Update{}, Pagination: pagination}, http.StatusOK
	} else {
		logrus.Debug("Returning concept updates")
		totalPages := int(math.Ceil(float64(count) / float64(validatedSize)))
		pagination := model.Pagination{TotalPages: totalPages, Page: validatedPage, Size: validatedSize}
		return model.Updates{Updates: databaseUpdates, Pagination: pagination}, http.StatusOK
	}
}
