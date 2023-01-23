package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/Informasjonsforvaltning/catalog-history-service/repository"
	"github.com/sirupsen/logrus"
)

type Update struct {
	Person     Person               `json:"person"`
	DateTime   time.Time            `json:"datetime"`
	Operations []JsonPatchOperation `json:"operations"`
}

type Person struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
type JsonPatchOperation struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}
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
