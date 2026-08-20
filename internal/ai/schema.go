package ai

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
)

func GenerateSchema[T any]() (map[string]any, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	var value T

	schema := reflector.Reflect(value)

	data, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to marshal schema: %w",
			err,
		)
	}

	var rawSchema map[string]json.RawMessage

	if err := json.Unmarshal(
		data,
		&rawSchema,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to decode schema: %w",
			err,
		)
	}

	result := make(
		map[string]any,
		len(rawSchema),
	)

	for key, value := range rawSchema {
		result[key] = value
	}

	return result, nil
}