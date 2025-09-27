package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Monitor interface {
	GetID() int64
	Type() string
	As(any) error
}

type Base struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	PathName        string  `json:"pathName"`
	Interval        int64   `json:"interval"`
	RetryInterval   int64   `json:"retryInterval"`
	ResendInterval  int64   `json:"resendInterval"`
	MaxRetries      int64   `json:"maxretries"`
	UpsideDown      bool    `json:"upsideDown"`
	NotificationIDs []int64 `json:"-"`
	IsActive        bool    `json:"active"`

	internalType string
	raw          []byte
}

func (b Base) String() string {
	return formatMonitor(b, true)
}

func (b *Base) UnmarshalJSON(data []byte) error {
	raw := struct {
		ID              int64           `json:"id"`
		Type            string          `json:"type"`
		Name            string          `json:"name"`
		Description     *string         `json:"description"`
		PathName        string          `json:"pathName"`
		Interval        int64           `json:"interval"`
		RetryInterval   int64           `json:"retryInterval"`
		ResendInterval  int64           `json:"resendInterval"`
		MaxRetries      int64           `json:"maxretries"`
		UpsideDown      bool            `json:"upsideDown"`
		NotificationIDs map[string]bool `json:"notificationIDList"`
		IsActive        bool            `json:"active"`
	}{}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	*b = Base{
		ID:             raw.ID,
		Name:           raw.Name,
		Description:    raw.Description,
		PathName:       raw.PathName,
		Interval:       raw.Interval,
		RetryInterval:  raw.RetryInterval,
		ResendInterval: raw.ResendInterval,
		MaxRetries:     raw.MaxRetries,
		UpsideDown:     raw.UpsideDown,
		IsActive:       raw.IsActive,

		internalType: raw.Type,
		raw:          data,
	}

	for id := range orderedByKey(raw.NotificationIDs) {
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return fmt.Errorf("notification ID is not int64: %w", err)
		}

		b.NotificationIDs = append(b.NotificationIDs, i)
	}

	return nil
}

func (b Base) MarshalJSON() ([]byte, error) {
	if b.raw == nil || b.internalType == "" {
		return nil, fmt.Errorf("not unmarshaled notification, unable to marshal")
	}

	raw := map[string]any{}
	err := json.Unmarshal(b.raw, &raw)
	if err != nil {
		return nil, fmt.Errorf("invalid internal state for raw, failed to unmarshal: %w", err)
	}

	raw["id"] = b.ID
	raw["type"] = b.internalType
	raw["name"] = b.Name
	raw["description"] = b.Description
	raw["pathName"] = b.PathName
	raw["interval"] = b.Interval
	raw["retryInterval"] = b.RetryInterval
	raw["resendInterval"] = b.ResendInterval
	raw["maxretries"] = b.MaxRetries
	raw["upsideDown"] = b.UpsideDown
	raw["active"] = b.IsActive

	ids := map[string]bool{}
	for _, id := range b.NotificationIDs {
		ids[strconv.FormatInt(id, 10)] = true
	}

	raw["notificationIDList"] = ids

	return json.Marshal(raw)
}

func (b Base) GetID() int64 {
	return b.ID
}

func (b Base) Type() string {
	return b.internalType
}

func (b Base) As(target any) error {
	if b.raw == nil {
		return fmt.Errorf("not unmarshaled monitor, cannot convert to %T", target)
	}

	err := json.Unmarshal(b.raw, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal monitor to %T: %w", target, err)
	}

	return nil
}
