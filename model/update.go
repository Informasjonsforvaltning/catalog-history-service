package model

import (
	"fmt"
	"time"
)

type UpdateDbo struct {
	ID         string               `json:"id" bson:"id"`
	ResourceId string               `json:"resourceId" bson:"resourceId"`
	Person     Person               `json:"person" bson:"person"`
	DateTime   time.Time            `json:"datetime" bson:"datetime"`
	Operations []JsonPatchOperation `json:"operations" bson:"operations"`
}

type UpdateDto struct {
	Person     Person               `json:"person" bson:"person"`
	Operations []JsonPatchOperation `json:"operations" bson:"operations"`
}

func (update UpdateDto) Validate() error {
	if 0 < len(update.Operations) {
		return nil
	}
	return fmt.Errorf("Update is not valid")
}

type Person struct {
	ID    string `json:"id" bson:"id"`
	Email string `json:"email" bson:"email"`
	Name  string `json:"name" bson:"name"`
}

type JsonPatchOperation struct {
	Op    string `json:"op" bson:"op"`
	Path  string `json:"path" bson:"path"`
	Value string `json:"value" bson:"value"`
}

type UpdateMeta struct {
	ID         string    `json:"id" bson:"id"`
	ResourceId string    `json:"resourceId" bson:"resourceId"`
	DateTime   time.Time `json:"datetime" bson:"datetime"`
	Person     Person    `json:"person" bson:"person"`
}

type UpdateDiff struct {
	ResourceId string               `json:"resourceId"`
	Operations []JsonPatchOperation `json:"operations"`
}
