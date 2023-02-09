package service

import (
	"encoding/json"

	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	jsonpatch "github.com/evanphx/json-patch"
)

func jsonPatchOperationsToByteArray(ops []model.JsonPatchOperation) ([]byte, error) {
	jsonBytes, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func applyPatchesToResource(original []byte, databaseUpdates model.UpdateDbo) ([]byte, error) {

	patchJSON, err := jsonPatchOperationsToByteArray(databaseUpdates.Operations)
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return nil, err
	}

	return patch.Apply(original)
}
