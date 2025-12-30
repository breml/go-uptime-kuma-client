package kuma

import (
	"encoding/json"
	"fmt"
)

// structToMap converts a struct to a map using JSON marshaling and unmarshaling.
// This allows structs to be converted to their map representation.
func structToMap(v any) (map[string]any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal struct: %w", err)
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal to map: %w", err)
	}

	return result, nil
}

// convertToStruct converts data from one type to another using JSON marshaling.
// This is useful for converting between different struct types.
func convertToStruct(src any, dst any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal source: %w", err)
	}

	err = json.Unmarshal(data, dst)
	if err != nil {
		return fmt.Errorf("unmarshal to struct: %w", err)
	}

	return nil
}
