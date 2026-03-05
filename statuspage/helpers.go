package statuspage

const (
	themeLight = "light"
	themeDark  = "dark"
	themeAuto  = "auto"

	styleInfo    = "info"
	styleWarning = "warning"
	styleDanger  = "danger"
	stylePrimary = "primary"

	analyticsTypeGoogle    = "google"
	analyticsTypeUmami     = "umami"
	analyticsTypePlausible = "plausible"
	analyticsTypeMatomo    = "matomo"
)

// ThemeLight returns the "light" theme identifier.
func ThemeLight() string {
	return themeLight
}

// ThemeDark returns the "dark" theme identifier.
func ThemeDark() string {
	return themeDark
}

// ThemeAuto returns the "auto" theme identifier.
func ThemeAuto() string {
	return themeAuto
}

// ValidTheme checks if the provided theme is valid.
func ValidTheme(theme string) bool {
	return theme == themeLight || theme == themeDark || theme == themeAuto
}

// StyleInfo returns the "info" incident style identifier.
func StyleInfo() string {
	return styleInfo
}

// StyleWarning returns the "warning" incident style identifier.
func StyleWarning() string {
	return styleWarning
}

// StyleDanger returns the "danger" incident style identifier.
func StyleDanger() string {
	return styleDanger
}

// StylePrimary returns the "primary" incident style identifier.
func StylePrimary() string {
	return stylePrimary
}

// ValidIncidentStyle checks if the provided incident style is valid.
func ValidIncidentStyle(style string) bool {
	return style == styleInfo || style == styleWarning || style == styleDanger || style == stylePrimary
}

// AnalyticsTypeGoogle returns the "google" analytics type identifier.
func AnalyticsTypeGoogle() string {
	return analyticsTypeGoogle
}

// AnalyticsTypeUmami returns the "umami" analytics type identifier.
func AnalyticsTypeUmami() string {
	return analyticsTypeUmami
}

// AnalyticsTypePlausible returns the "plausible" analytics type identifier.
func AnalyticsTypePlausible() string {
	return analyticsTypePlausible
}

// AnalyticsTypeMatomo returns the "matomo" analytics type identifier.
func AnalyticsTypeMatomo() string {
	return analyticsTypeMatomo
}

// ValidAnalyticsType checks if the provided analytics type is valid.
// A nil value is considered valid, as Uptime Kuma accepts null for this field.
func ValidAnalyticsType(analyticsType *string) bool {
	if analyticsType == nil {
		return true
	}

	at := *analyticsType
	return at == analyticsTypeGoogle || at == analyticsTypeUmami ||
		at == analyticsTypePlausible ||
		at == analyticsTypeMatomo
}
