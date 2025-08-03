package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/model"
)

func TestNotificationConfig_JSON(t *testing.T) {
	notification := model.Notification{
		ID:        1,
		Name:      "Ntfy",
		Active:    true,
		UserID:    1,
		IsDefault: false,
		Config: model.NotificationConfig{
			"isDefault":                false,
			"name":                     "Ntfy",
			"ntfyAuthenticationMethod": "none",
			"ntfyPriority":             5.0,
			"ntfyserverurl":            "https://ntfy.sh",
			"ntfytopic":                "ntfytopic",
			"type":                     "ntfy",
		},
	}

	data, err := json.Marshal(notification)
	if err != nil {
		t.Fatalf("no error expected, got: %v", err)
	}

	got := model.Notification{}
	err = json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("no error expected, got: %v", err)
	}

	require.Equal(t, notification, got)
}
