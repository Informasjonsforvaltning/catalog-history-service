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
	var updateDbo = model.UpdateDbo{
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

func (service *UpdateServiceImp) GetConceptUpdates(ctx context.Context, conceptId string, pageSize int, pageNumber int) (*[]*model.UpdateMeta, int) {
	query := bson.D{}
	query = append(query, bson.E{Key: "resourceId", Value: conceptId})
	query = append(query, bson.E{Key: "datetime", Value: bson.M{"$lt": time.Now()}})

	// Calculate the number of results to skip based on the page size and number.
	skip := pageSize * (pageNumber - 1)

	databaseUpdates, err := service.ConceptsRepository.GetConceptUpdates(ctx, query, skip, pageSize)
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
func (service *UpdateServiceImp) GetConceptUpdate(ctx context.Context, conceptId string, updateId string) (*model.UpdateDbo, int) {
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

// function to get a diff from database
func (service *UpdateServiceImp) GetConceptUpdateDiff(ctx context.Context, conceptId string, updateId string) (*model.UpdateDiff, int) {
	conceptUpdate, err := service.ConceptsRepository.GetConceptUpdate(ctx, conceptId, updateId)
	if err != nil {
		logrus.Error("Unable to get concept update")
		return nil, http.StatusInternalServerError
	}
	if conceptUpdate == nil {
		logrus.Error("Concept update not found")
		return nil, http.StatusNotFound
	}
	query := bson.D{}
	query = append(query, bson.E{Key: "resourceId", Value: conceptId})
	query = append(query, bson.E{Key: "datetime", Value: bson.M{"$lt": conceptUpdate.DateTime}})
	databaseUpdates, listErr := service.ConceptsRepository.GetConceptUpdates(ctx, query)
	if listErr != nil {
		logrus.Error("Get concept updates failed")
		return nil, http.StatusInternalServerError
	}
	concept, err := service.BuildResourceFromPatches(databaseUpdates)
	if err != nil {
		logrus.Error("Unable to build concept from patches")
		logrus.Error(err)
		return nil, http.StatusInternalServerError

	}
	return &model.UpdateDiff{
		ResourceId: conceptUpdate.ResourceId,
		Operations: conceptUpdate.Operations,
		Resource:   string(concept),
	}, http.StatusOK
}

func BuildResourceFromPatches(databaseUpdates []*model.UpdateDbo) ([]byte, error) {
	resource := []byte("{}")
	var err error

	for _, update := range databaseUpdates {
		if update == nil {
			logrus.Warning("Update is nil")
		} else {
			logrus.Info("Update: " + update.ID)
			resource, err = applyPatchesToResource(resource, *update)
			if err != nil {
				return nil, err
			}
		}
	}
	return resource, nil

}
