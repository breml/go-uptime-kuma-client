package notification

import (
	"encoding/json"
	"fmt"
)

// Webhook represents a webhook notification provider.
// It supports multiple content types ('json', 'form-data', 'custom') and custom headers.
type Webhook struct {
	Base
	WebhookDetails
}

// WebhookDetails contains the webhook-specific configuration.
type WebhookDetails struct {
	WebhookURL               string                   `json:"webhookURL"`
	WebhookContentType       string                   `json:"webhookContentType"`
	WebhookCustomBody        string                   `json:"webhookCustomBody,omitempty"`
	WebhookAdditionalHeaders WebhookAdditionalHeaders `json:"webhookAdditionalHeaders,omitempty"`
}

// WebhookAdditionalHeaders is a map of HTTP headers that marshals to/from a JSON string.
// The Uptime Kuma server expects the headers as a JSON-encoded string, not a nested object.
type WebhookAdditionalHeaders map[string]string

// Type returns the notification type identifier for webhooks.
func (w Webhook) Type() string {
	return w.WebhookDetails.Type()
}

// Type returns the notification type identifier for webhook details.
func (d WebhookDetails) Type() string {
	return "webhook"
}

// String returns a human-readable representation of the webhook notification.
func (w Webhook) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(w.Base, false), formatNotification(w.WebhookDetails, true))
}

// UnmarshalJSON deserializes a webhook notification from JSON.
func (w *Webhook) UnmarshalJSON(data []byte) error {
	detail := WebhookDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*w = Webhook{
		Base:           base,
		WebhookDetails: detail,
	}

	return nil
}

// MarshalJSON serializes a webhook notification to JSON.
func (w Webhook) MarshalJSON() ([]byte, error) {
	return marshalJSON(w.Base, w.WebhookDetails)
}

// MarshalJSON serializes the headers map to a JSON string.
// Example: {"Authorization": "Bearer token"} becomes "{\"Authorization\":\"Bearer token\"}".
func (h WebhookAdditionalHeaders) MarshalJSON() ([]byte, error) {
	if h == nil {
		return []byte("null"), nil
	}

	// First marshal the map to JSON bytes
	mapJSON, err := json.Marshal(map[string]string(h))
	if err != nil {
		return nil, fmt.Errorf("marshal webhook headers map: %w", err)
	}

	// Then marshal the JSON bytes as a string
	data, err := json.Marshal(string(mapJSON))
	if err != nil {
		return nil, fmt.Errorf("marshal webhook headers: %w", err)
	}
	return data, nil
}

// UnmarshalJSON deserializes a JSON string into the headers map.
// Example: "{\"Authorization\":\"Bearer token\"}" becomes {"Authorization": "Bearer token"}.
func (h *WebhookAdditionalHeaders) UnmarshalJSON(data []byte) error {
	// First check if it's null or empty string
	if string(data) == "null" || string(data) == `""` {
		*h = nil
		return nil
	}

	// Unmarshal the outer JSON string
	var jsonStr string
	err := json.Unmarshal(data, &jsonStr)
	if err != nil {
		return fmt.Errorf("unmarshal webhook headers outer json: %w", err)
	}

	// If the string is empty after unmarshaling, treat as nil
	if jsonStr == "" {
		*h = nil
		return nil
	}

	// Unmarshal the inner JSON object
	var headers map[string]string
	err = json.Unmarshal([]byte(jsonStr), &headers)
	if err != nil {
		return fmt.Errorf("unmarshal webhook headers inner json: %w", err)
	}

	*h = WebhookAdditionalHeaders(headers)
	return nil
}
