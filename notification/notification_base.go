package notification

import (
	"encoding/json"
	"fmt"
)

type Notification interface {
	GetID() int64
	Type() string
	As(any) error
}

type Base struct {
	ID            int64  `json:"id,omitzero"`
	Name          string `json:"name"`
	IsActive      bool   `json:"active"`
	IsDefault     bool   `json:"isDefault"`
	ApplyExisting bool   `json:"applyExisting"`
	UserID        int64  `json:"userId"`

	typeFromConfigStr string
	configStr         string
	raw               []byte
}

func (b Base) String() string {
	return formatNotification(b, true)
}

func (b *Base) UnmarshalJSON(data []byte) error {
	raw := struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		IsActive  bool   `json:"active"`
		UserID    int64  `json:"userId"`
		IsDefault bool   `json:"isDefault"`
		ConfigStr string `json:"config"`
	}{}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	config := map[string]any{}

	err = json.Unmarshal([]byte(raw.ConfigStr), &config)
	if err != nil {
		return err
	}

	notificationTypeAny, ok := config["type"]
	if !ok {
		return fmt.Errorf(`invalid notification, attribute "type" missing`)
	}

	notificationType, ok := notificationTypeAny.(string)
	if !ok {
		return fmt.Errorf(`invalid notification, attribute "type" is not string`)
	}

	applyExistingAny, ok := config["applyExisting"]
	if !ok {
		return fmt.Errorf(`invalid notification, attribute "applyExisting" missing`)
	}

	applyExisting, ok := applyExistingAny.(bool)
	if !ok {
		return fmt.Errorf(`invalid notification, attribute "applyExisting" is not bool`)
	}

	*b = Base{
		ID:            raw.ID,
		Name:          raw.Name,
		IsActive:      raw.IsActive,
		UserID:        raw.UserID,
		IsDefault:     raw.IsDefault,
		ApplyExisting: applyExisting,

		typeFromConfigStr: notificationType,
		configStr:         raw.ConfigStr,
		raw:               data,
	}

	return nil
}

func (b Base) GetID() int64 {
	return b.ID
}

func (b Base) Type() string {
	return b.typeFromConfigStr
}

func (b Base) As(target any) error {
	if b.raw == nil {
		return fmt.Errorf("not unmarshaled notification, cannot convert to %T", target)
	}

	err := json.Unmarshal(b.raw, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal notification to %T: %w", target, err)
	}

	return nil
}
