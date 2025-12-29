package statuspage

// StatusPage represents a public status page.
type StatusPage struct {
	ID                    int64         `json:"id"`
	Slug                  string        `json:"slug"`
	Title                 string        `json:"title"`
	Description           string        `json:"description"`
	Icon                  string        `json:"icon"`
	Theme                 string        `json:"theme"`
	Published             bool          `json:"published"`
	ShowTags              bool          `json:"showTags"`
	DomainNameList        []string      `json:"domainNameList"`
	GoogleAnalyticsID     string        `json:"googleAnalyticsId"`
	CustomCSS             string        `json:"customCSS"`
	FooterText            string        `json:"footerText"`
	ShowPoweredBy         bool          `json:"showPoweredBy"`
	ShowCertificateExpiry bool          `json:"showCertificateExpiry"`
	PublicGroupList       []PublicGroup `json:"publicGroupList"`
}

// PublicGroup represents a group of monitors in a public status page.
type PublicGroup struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Weight      int             `json:"weight"`
	MonitorList []PublicMonitor `json:"monitorList"`
}

// PublicMonitor represents a monitor in a public status page.
type PublicMonitor struct {
	ID      int64 `json:"id"`
	SendURL *bool `json:"sendUrl,omitempty"`
}
