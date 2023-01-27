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
	var update model.UpdateDto
	err := json.Unmarshal(bytes, &update)
	if err != nil {
		logrus.Error("Unable to unmarshal concept update")
		return nil, err
	}
	err = update.Validate()
	if err != nil {
		logrus.Error("Concept update is not valid")
		return nil, err
	}
	var updateDbo = model.UpdateDbo{
		ID:         uuid.New().String(),
		ResourceId: conceptId,
		DateTime:   time.Now(),
		Person:     update.Person,
		Operations: update.Operations,
	}
	err = service.ConceptsRepository.StoreConcept(ctx, updateDbo)
	if err != nil {
		logrus.Error("Could not store concept update")
		return nil, err
	}
	return &updateDbo.ID, nil
}

func (service *UpdateServiceImp) GetConceptUpdates(ctx context.Context, conceptId *string) (*[]*model.UpdateMeta, int) {
	query := bson.D{}
	if conceptId != nil {
		query = append(query, bson.E{Key: "id", Value: conceptId})
	}
	databaseUpdates, err := service.ConceptsRepository.GetConceptUpdates(ctx, query)
	if err != nil {
		logrus.Error("Get concept updates failed")
		return nil, http.StatusInternalServerError
	}

	if databaseUpdates == nil {
		databaseUpdates = []*model.UpdateDbo{}
	}

	var updates []*model.UpdateMeta

	for _, update := range databaseUpdates {
		updates = append(updates, &model.UpdateMeta{
			ID:         update.ID,
			ResourceId: update.ResourceId,
			DateTime:   update.DateTime,
			Person:     update.Person,
		})
	}

	return &updates, http.StatusOK
}

// function to get a update from database
func (service *UpdateServiceImp) GetConceptUpdate(ctx context.Context, conceptId string) (*model.UpdateDbo, int) {
	conceptUpdate, err := service.ConceptsRepository.GetConceptUpdate(ctx, conceptId)
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
