package model

import (
	"fmt"
	"time"
)

type Update struct {
	ID         string               `json:"id"`
	ResourceId string               `json:"resourceId"`
	Person     Person               `json:"person"`
	DateTime   time.Time            `json:"datetime"`
	Operations []JsonPatchOperation `json:"operations"`
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
	Updates []*Update `json:"updates"`
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
