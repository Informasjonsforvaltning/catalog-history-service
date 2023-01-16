package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/model"
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
type UpdateService struct {
}

type JsonPatchOperation struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func StoreUpdate(w http.ResponseWriter, r *http.Request) {
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
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Update stored"))
}

func NewUpdateService() *UpdateService {
	return &UpdateService{}
}
