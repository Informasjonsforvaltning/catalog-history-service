package model

import (
	"fmt"
	"time"
)

type UpdateDbo struct {
	ID         string               `json:"id"`
	ResourceId string               `json:"resourceId"`
	Person     Person               `json:"person"`
	DateTime   time.Time            `json:"datetime"`
	Operations []JsonPatchOperation `json:"operations"`
}

type UpdateDto struct {
	Person     Person               `json:"person"`
	Operations []JsonPatchOperation `json:"operations"`
}

func (update UpdateDto) Validate() error {
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

type UpdateMeta struct {
	ID         string    `json:"id"`
	ResourceId string    `json:"resourceId"`
	DateTime   time.Time `json:"datetime"`
	Person     Person    `json:"person"`
}

type UpdateDiff struct {
	ResourceId string               `json:"resourceId"`
	Operations []JsonPatchOperation `json:"operations"`
}
