package kuma

import (
	"encoding/json"
	"fmt"
)

func structToMap(v any) (map[string]any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal struct: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal to map: %v", err)
	}

	return result, nil
}

func convertToStruct(src any, dst any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("marshal source: %v", err)
	}

	err = json.Unmarshal(data, dst)
	if err != nil {
		return fmt.Errorf("unmarshal to struct: %v", err)
	}

	return nil
}
