package notification

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Notification is the interface that all notification types must implement.
type Notification interface {
	// GetID returns the notification's unique identifier.
	GetID() int64
	// Type returns the notification's type name.
	Type() string
	// As converts the notification to the given target type.
	As(any) error
}

// Base contains the common fields for all notification types.
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

// UnmarshalJSON unmarshals a notification from JSON data.
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
		return fmt.Errorf("unmarshal notification base: %w", err)
	}

	config := map[string]any{}

	err = json.Unmarshal([]byte(raw.ConfigStr), &config)
	if err != nil {
		return fmt.Errorf("unmarshal notification config: %w", err)
	}

	notificationTypeAny, ok := config["type"]
	if !ok {
		return errors.New(`invalid notification, attribute "type" missing`)
	}

	notificationType, ok := notificationTypeAny.(string)
	if !ok {
		return errors.New(`invalid notification, attribute "type" is not string`)
	}

	var applyExisting bool
	applyExistingAny, ok := config["applyExisting"]
	if ok {
		applyExisting, ok = applyExistingAny.(bool)
		if !ok {
			return errors.New(`invalid notification, attribute "applyExisting" is not bool`)
		}
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

// MarshalJSON marshals a notification to JSON data.
func (b Base) MarshalJSON() ([]byte, error) {
	if b.configStr == "" {
		return nil, errors.New("not unmarshaled notification, unable to marshal")
	}

	genericDetails := GenericDetails{}
	err := json.Unmarshal([]byte(b.configStr), &genericDetails)
	if err != nil {
		return nil, fmt.Errorf("invalid internal state for configStr, failed to unmarshal: %w", err)
	}

	return marshalJSON(b, genericDetails)
}

// GetID returns the notification's unique identifier.
func (b Base) GetID() int64 {
	return b.ID
}

// Type returns the notification's type name.
func (b Base) Type() string {
	return b.typeFromConfigStr
}

// As converts the notification to the given target type.
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
