package repository

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateRepository interface {
	StoreUpdate(ctx context.Context, update model.Update) error
}

type UpdateRepositoryImpl struct {
	collection *mongo.Collection
}
type UpdateService struct {
	UpdateRepository UpdateRepository
}

var updateRepository *UpdateRepositoryImpl

func InitUpdateRepository() *UpdateRepositoryImpl {
	if updateRepository == nil {
		updateRepository = &UpdateRepositoryImpl{collection: connection.MongoCollection()}
	}
	return updateRepository
}

func (r *UpdateRepositoryImpl) createUpdate(ctx context.Context, update model.Update) error {
	_, err := r.collection.InsertOne(ctx, update, nil)
	return err
}

func CreateUpdate(w http.ResponseWriter, r *http.Request) {
	var update model.Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//validating the update struct
	if update.Person.ID == "" || update.Person.Email == "" || update.Person.Name == "" || update.DateTime.IsZero() {
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}
	err = InitUpdateRepository().createUpdate(r.Context(), update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
