package settings

// Settings represents the general server settings of Uptime Kuma.
type Settings struct {
	ServerTimezone      string `json:"serverTimezone"`
	KeepDataPeriodDays  int    `json:"keepDataPeriodDays"`
	CheckUpdate         bool   `json:"checkUpdate"`
	SearchEngineIndex   bool   `json:"searchEngineIndex"`
	EntryPage           string `json:"entryPage"`
	NSCD                bool   `json:"nscd"`
	TLSExpiryNotifyDays []int  `json:"tlsExpiryNotifyDays"`
	TrustProxy          bool   `json:"trustProxy"`
	PrimaryBaseURL      string `json:"primaryBaseURL"`
	SteamAPIKey         string `json:"steamAPIKey"`
	ChromeExecutable    string `json:"chromeExecutable"`
}
