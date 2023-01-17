package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/Informasjonsforvaltning/catalog-history-service/service"
)

func NewUpdateHandler(us service.UpdateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
