package model

import (
	"fmt"
	"time"
)

type Update struct {
	ID         string               `bson:"_id" json:"id"`
	CatalogId  string               `bson:"catalogId" json:"catalogId"`
	ResourceId string               `bson:"resourceId" json:"resourceId"`
	Person     Person               `bson:"person" json:"person"`
	DateTime   time.Time            `bson:"datetime" json:"datetime"`
	Operations []JsonPatchOperation `bson:"operations" json:"operations"`
}

type UpdatePayload struct {
	Person     Person               `json:"person"`
	Operations []JsonPatchOperation `json:"operations"`
}

func (update UpdatePayload) Validate() error {
	if 0 < len(update.Operations) {
		return nil
	}
	return fmt.Errorf("Update is not valid")
}

type Updates struct {
	Updates []Update `json:"updates"`
}

type Person struct {
	ID    string `bson:"id" json:"id"`
	Email string `bson:"email" json:"email"`
	Name  string `bson:"name" json:"name"`
}

type JsonPatchOperation struct {
	Op    string `bson:"op" json:"op"`
	Path  string `bson:"path" json:"path"`
	Value any    `bson:"value" json:"value"`
}
