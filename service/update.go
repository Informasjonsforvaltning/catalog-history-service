package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/Informasjonsforvaltning/catalog-history-service/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateService interface {
	StoreUpdate(ctx context.Context, bytes []byte)
}

type UpdateServiceImp struct {
	ConceptsRepository repository.ConceptsRepositoryImp
}

func InitService() *UpdateServiceImp {
	service := UpdateServiceImp{
		ConceptsRepository: *repository.InitRepository(),
	}
	return &service
}

func (service *UpdateServiceImp) StoreConceptUpdate(ctx context.Context, bytes []byte, conceptId string) (*string, error) {
	var update model.UpdatePayload
	err := json.Unmarshal(bytes, &update)
	logrus.Info("Unmarshalled update")
	if err != nil {
		logrus.Error("Unable to unmarshal concept update")
		return nil, err
	}
	err = update.Validate()
	if err != nil {
		logrus.Error("Concept update is not valid")
		return nil, err
	}
	var updateDbo = model.Update{
		ID:         uuid.New().String(),
		ResourceId: conceptId,
		DateTime:   time.Now(),
		Person:     update.Person,
		Operations: update.Operations,
	}
	err = service.ConceptsRepository.StoreConcept(ctx, updateDbo)
	if err != nil {
		logrus.Error("Could not store concept update" + err.Error())
		return nil, err
	}
	return &updateDbo.ID, nil
}

func (service *UpdateServiceImp) GetConceptUpdates(ctx context.Context, conceptId string, page int, size int, sortBy string, sortOrder string) (*model.Updates, int) {
	query := bson.D{}
	query = append(query, bson.E{Key: "resourceId", Value: conceptId})

	// Set default sort by column to "datetime"
	sortByCol := "datetime"
	switch sortBy {
	case "name":
		sortByCol = "person.name"
	case "email":
		sortByCol = "person.email"
	}

	// Map sortOrder string to integer value
	var sortOrderInt int
	switch sortOrder {
	case "asc":
		sortOrderInt = 1
	case "desc":
		sortOrderInt = -1
	default:
		sortOrderInt = -1 // default to descending order
	}

	databaseUpdates, err := service.ConceptsRepository.GetConceptUpdates(ctx, query, page, size, sortByCol, sortOrderInt)
	if err != nil {
		logrus.Error("Get concept updates failed")
		return nil, http.StatusInternalServerError
	}

	if databaseUpdates == nil {
		logrus.Error("No concept updates found")
		return &model.Updates{Updates: []*model.Update{}}, http.StatusOK
	} else {
		logrus.Info("Returning concept updates")
		return &model.Updates{Updates: databaseUpdates}, http.StatusOK
	}
}

// function to get a update from database
func (service *UpdateServiceImp) GetConceptUpdate(ctx context.Context, conceptId string, updateId string) (*model.Update, int) {
	conceptUpdate, err := service.ConceptsRepository.GetConceptUpdate(ctx, conceptId, updateId)
	if err != nil {
		logrus.Error("Unable to get concept update")
		return nil, http.StatusInternalServerError
	} else if conceptUpdate == nil {
		logrus.Error("Concept update not found")
		return nil, http.StatusNotFound
	} else {
		return conceptUpdate, http.StatusOK
	}
}
