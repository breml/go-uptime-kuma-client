package model

import (
	"encoding/json"
)

// TODO: consider to define marshal and unmarshal on Notification and perform the splitting automatically.
// This would also allow to e.g. move the type attribute from config to the parent level.

type Notification struct {
	ID        int                `json:"id"`
	Name      string             `json:"name"`
	Active    bool               `json:"active"`
	UserID    int                `json:"userId"`
	IsDefault bool               `json:"isDefault"`
	Config    NotificationConfig `json:"config"`
}

func (n Notification) Writeable() map[string]any {
	config := make(map[string]any, len(n.Config)+2)
	for k, v := range n.Config {
		config[k] = v
	}
	config["name"] = n.Name
	config["isDefault"] = n.IsDefault

	return config
}

func (n Notification) Validate() error {
	// TODO: Ensure, that provided attributes in config match the type of the notification
	// and the expected type of the value.
	return nil
}

type NotificationConfig map[string]any

func (n NotificationConfig) MarshalJSON() ([]byte, error) {
	config := map[string]any(n)
	body, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	str := string(body)
	b, err := json.Marshal(str)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (n *NotificationConfig) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	config := map[string]any{}
	err = json.Unmarshal([]byte(str), &config)
	if err != nil {
		return err
	}

	*n = config

	return nil
}
