package settings_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/go-uptime-kuma-client/settings"
)

func TestSettings_MarshalJSON(t *testing.T) {
	s := settings.Settings{
		ServerTimezone:      "Europe/Berlin",
		KeepDataPeriodDays:  180,
		CheckUpdate:         true,
		SearchEngineIndex:   false,
		EntryPage:           "dashboard",
		NSCD:                true,
		TLSExpiryNotifyDays: []int{7, 14, 21},
		TrustProxy:          false,
		PrimaryBaseURL:      "https://uptime.example.com",
		SteamAPIKey:         "secret-key",
		ChromeExecutable:    "/usr/bin/chromium",
	}

	data, err := json.Marshal(s)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	require.Equal(t, "Europe/Berlin", result["serverTimezone"])
	require.InEpsilon(t, float64(180), result["keepDataPeriodDays"], 0)
	require.Equal(t, true, result["checkUpdate"])
	require.Equal(t, false, result["searchEngineIndex"])
	require.Equal(t, "dashboard", result["entryPage"])
	require.Equal(t, true, result["nscd"])
	require.Len(t, result["tlsExpiryNotifyDays"], 3)
	require.Equal(t, false, result["trustProxy"])
	require.Equal(t, "https://uptime.example.com", result["primaryBaseURL"])
	require.Equal(t, "secret-key", result["steamAPIKey"])
	require.Equal(t, "/usr/bin/chromium", result["chromeExecutable"])
}

func TestSettings_UnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"serverTimezone": "Europe/Berlin",
		"keepDataPeriodDays": 7,
		"checkUpdate": true,
		"searchEngineIndex": false,
		"entryPage": "dashboard",
		"nscd": true,
		"tlsExpiryNotifyDays": [7, 14, 21],
		"trustProxy": false,
		"primaryBaseURL": "https://uptime.example.com",
		"steamAPIKey": "my-api-key",
		"chromeExecutable": "/usr/bin/chromium"
	}`

	var s settings.Settings
	err := json.Unmarshal([]byte(jsonStr), &s)
	require.NoError(t, err)

	require.Equal(t, "Europe/Berlin", s.ServerTimezone)
	require.Equal(t, 7, s.KeepDataPeriodDays)
	require.True(t, s.CheckUpdate)
	require.False(t, s.SearchEngineIndex)
	require.Equal(t, "dashboard", s.EntryPage)
	require.True(t, s.NSCD)
	require.Equal(t, []int{7, 14, 21}, s.TLSExpiryNotifyDays)
	require.False(t, s.TrustProxy)
	require.Equal(t, "https://uptime.example.com", s.PrimaryBaseURL)
	require.Equal(t, "my-api-key", s.SteamAPIKey)
	require.Equal(t, "/usr/bin/chromium", s.ChromeExecutable)
}

func TestSettings_RoundTrip(t *testing.T) {
	original := settings.Settings{
		ServerTimezone:      "America/New_York",
		KeepDataPeriodDays:  30,
		CheckUpdate:         false,
		SearchEngineIndex:   true,
		EntryPage:           "statusPage-default",
		NSCD:                false,
		TLSExpiryNotifyDays: []int{7, 14},
		TrustProxy:          true,
		PrimaryBaseURL:      "https://status.example.com",
		SteamAPIKey:         "",
		ChromeExecutable:    "",
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var roundTripped settings.Settings
	err = json.Unmarshal(data, &roundTripped)
	require.NoError(t, err)

	require.Equal(t, original, roundTripped)
}

func TestSettings_String(t *testing.T) {
	s := settings.Settings{
		ServerTimezone:      "UTC",
		KeepDataPeriodDays:  180,
		CheckUpdate:         true,
		EntryPage:           "dashboard",
		TLSExpiryNotifyDays: []int{7, 14, 21},
		SteamAPIKey:         "secret",
	}

	str := s.String()
	require.Contains(t, str, `serverTimezone: "UTC"`)
	require.Contains(t, str, "keepDataPeriodDays: 180")
	require.Contains(t, str, `entryPage: "dashboard"`)
	require.Contains(t, str, `steamAPIKey: "***"`)
	require.NotContains(t, str, "secret")
}

func TestSettings_String_EmptySteamAPIKey(t *testing.T) {
	s := settings.Settings{
		SteamAPIKey: "",
	}

	str := s.String()
	require.Contains(t, str, `steamAPIKey: ""`)
}
