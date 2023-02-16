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

func (service *UpdateServiceImp) GetConceptUpdates(ctx context.Context, conceptId string) (*model.Updates, int) {
	query := bson.D{}
	query = append(query, bson.E{Key: "resourceId", Value: conceptId})
	databaseUpdates, err := service.ConceptsRepository.GetConceptUpdates(ctx, query)
	if err != nil {
		logrus.Error("Get concept updates failed")
		return nil, http.StatusInternalServerError
	}

	if databaseUpdates == nil {
		return &model.Updates{Updates: []*model.Update{}}, http.StatusOK
	} else {
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
