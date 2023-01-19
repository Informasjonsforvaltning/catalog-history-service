package model

import (
	"fmt"
	"time"
)

type Update struct {
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
