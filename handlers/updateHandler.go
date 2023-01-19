package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/model"

	jsonpatch "github.com/evanphx/json-patch"
)

func UpdateJsonPatch(w http.ResponseWriter, r *http.Request) {
	original := &model.Update{
		Person: model.Person{
			ID:    "123",
			Email: "emaill",
			Name:  "name",
		},
		DateTime: time.Now(),
		Operations: []model.JsonPatchOperation{
			{
				Op:    "replace",
				Path:  "/name",
				Value: "Jane",
			},
			{
				Op:   "remove",
				Path: "/height",
			},
		},
	}

	originalBytes, err := json.Marshal(original)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(originalBytes))
	request, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(request))
	patchedJson, err := jsonpatch.MergePatch(originalBytes, request)
	fmt.Println(string(patchedJson))

	w.Write(patchedJson)
	w.Header().Set("Content-Type", "application/json")
}
