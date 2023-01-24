package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/Informasjonsforvaltning/catalog-history-service/repository"
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

func (service *UpdateServiceImp) StoreConceptUpdate(ctx context.Context, bytes []byte) error {
	var update model.Update
	err := json.Unmarshal(bytes, &update)
	if err != nil {
		logrus.Error("Unable to unmarshal concept update")
		return err
	}
	err = update.Validate()
	if err != nil {
		logrus.Error("Concept update is not valid")
		return err
	}
	return service.ConceptsRepository.StoreConcept(ctx, update)
}

func (service *UpdateServiceImp) GetConceptUpdates(ctx context.Context, id *string) (*[]*model.Update, int) {
	query := bson.D{}
	if id != nil {
		query = append(query, bson.E{Key: "id", Value: id})
	}
	conceptUpdates, err := service.ConceptsRepository.GetConceptUpdates(ctx, query)
	if err != nil {
		logrus.Error("Get concept updates failed")
		return nil, http.StatusInternalServerError
	}

	if conceptUpdates == nil {
		conceptUpdates = []*model.Update{}
	}
	return &conceptUpdates, http.StatusOK
}

// function to get a update from database
func (service *UpdateServiceImp) GetConceptUpdate(ctx context.Context, id string) (*model.Update, int) {
	conceptUpdate, err := service.ConceptsRepository.GetConceptUpdate(ctx, id)
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
