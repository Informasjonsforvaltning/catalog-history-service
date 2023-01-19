package repository

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/connection"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateRepository interface {
	StoreUpdate(ctx context.Context, update model.Update) error
}

type UpdateRepositoryImpl struct {
	collection *mongo.Collection
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if update.Person.ID == "" || update.Person.Email == "" || update.Person.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Person must be set"))
		return
	}
	if len(update.Operations) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Operations must be set"))
		return
	}
	update.DateTime = time.Now()
	// create a mongo session
	session, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// close the session when done
	defer session.Disconnect(context.Background())
	// get the "concepts" collection
	collection := session.Database("catalog-history-service").Collection("concepts")
	// insert the update struct into the collection
	_, err = collection.InsertOne(context.TODO(), update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Update stored"))
}
