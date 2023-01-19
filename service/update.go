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

func (service *UpdateServiceImp) StoreUpdate(ctx context.Context, bytes []byte) error {
	var update model.Update
	err := json.Unmarshal(bytes, &update)
	logrus.Info("Called StoreUpdate")
	if err != nil {
		logrus.Info("Marshal failed")
		return err
	}
	err = update.Validate()
	if err != nil {
		logrus.Info("Validation failed")
		return err
	}
	logrus.Info("Validated update")
	return service.ConceptsRepository.StoreConcept(ctx, update)
}
