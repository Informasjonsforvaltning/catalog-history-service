package model

import (
	"fmt"
	"time"
)

type Update struct {
	ID         string               `json:"id" bson:"id"`
	Person     Person               `json:"person"`
	DateTime   time.Time            `json:"datetime"`
	Operations []JsonPatchOperation `json:"operations"`
}

func (update Update) Validate() error {
	if 0 < len(update.Operations) {
		return nil
	}
	return fmt.Errorf("Update is not valid")
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
