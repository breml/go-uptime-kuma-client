package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/breml/go-uptime-kuma-client/tag"
)

type Monitor interface {
	GetID() int64
	Type() string
	As(any) error
	GetNotificationIDs() []int64
}

type Base struct {
	ID              int64            `json:"id"`
	Name            string           `json:"name"`
	Description     *string          `json:"description"`
	PathName        string           `json:"pathName"`
	Parent          *int64           `json:"parent"`
	ProxyID         *int64           `json:"proxyId"`
	Interval        int64            `json:"interval"`
	RetryInterval   int64            `json:"retryInterval"`
	ResendInterval  int64            `json:"resendInterval"`
	MaxRetries      int64            `json:"maxretries"`
	UpsideDown      bool             `json:"upsideDown"`
	NotificationIDs []int64          `json:"-"`
	Tags            []tag.MonitorTag `json:"tags"`
	IsActive        bool             `json:"active"`

	internalType string
	raw          []byte
}

func (b Base) String() string {
	return formatMonitor(b, true)
}

func (b *Base) UnmarshalJSON(data []byte) error {
	raw := struct {
		ID              int64            `json:"id"`
		Type            string           `json:"type"`
		Name            string           `json:"name"`
		Description     *string          `json:"description"`
		PathName        string           `json:"pathName"`
		Parent          *int64           `json:"parent"`
		ProxyID         *int64           `json:"proxyId"`
		Interval        int64            `json:"interval"`
		RetryInterval   int64            `json:"retryInterval"`
		ResendInterval  int64            `json:"resendInterval"`
		MaxRetries      int64            `json:"maxretries"`
		UpsideDown      bool             `json:"upsideDown"`
		NotificationIDs map[string]bool  `json:"notificationIDList"`
		Tags            []tag.MonitorTag `json:"tags"`
		IsActive        bool             `json:"active"`
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
		Parent:         raw.Parent,
		ProxyID:        raw.ProxyID,
		Interval:       raw.Interval,
		RetryInterval:  raw.RetryInterval,
		ResendInterval: raw.ResendInterval,
		MaxRetries:     raw.MaxRetries,
		UpsideDown:     raw.UpsideDown,
		IsActive:       raw.IsActive,

		internalType: raw.Type,
		raw:          data,
	}

	// Only set Tags if it's not an empty slice to maintain nil semantics
	if len(raw.Tags) > 0 {
		b.Tags = raw.Tags
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

	// Monitor was unmarshaled from JSON, use existing raw data as base.
	err := json.Unmarshal(b.raw, &raw)
	if err != nil {
		return nil, fmt.Errorf("invalid internal state for raw, failed to unmarshal: %w", err)
	}

	raw["id"] = b.ID
	raw["type"] = b.getType()
	raw["name"] = b.Name
	raw["description"] = b.Description
	// Don't set pathName, server generates it.
	// raw["pathName"] = b.PathName
	raw["parent"] = b.Parent
	raw["proxyId"] = b.ProxyID
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

func (b Base) getType() string {
	return b.internalType
}

func (b Base) GetID() int64 {
	return b.ID
}

func (b Base) Type() string {
	return b.internalType
}

func (b Base) GetNotificationIDs() []int64 {
	return b.NotificationIDs
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
