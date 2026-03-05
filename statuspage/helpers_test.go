package statuspage_test

import (
	"testing"

	"github.com/breml/go-uptime-kuma-client/internal/ptr"
	"github.com/breml/go-uptime-kuma-client/statuspage"
)

func TestThemeHelpers(t *testing.T) {
	tests := []struct {
		name     string
		theme    string
		wantFunc func() string
	}{
		{
			name:     "light theme",
			theme:    "light",
			wantFunc: statuspage.ThemeLight,
		},
		{
			name:     "dark theme",
			theme:    "dark",
			wantFunc: statuspage.ThemeDark,
		},
		{
			name:     "auto theme",
			theme:    "auto",
			wantFunc: statuspage.ThemeAuto,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wantFunc(); got != tt.theme {
				t.Errorf("theme helper returned %q, want %q", got, tt.theme)
			}
		})
	}
}

func TestValidTheme(t *testing.T) {
	tests := []struct {
		name  string
		theme string
		want  bool
	}{
		{
			name:  "light is valid",
			theme: "light",
			want:  true,
		},
		{
			name:  "dark is valid",
			theme: "dark",
			want:  true,
		},
		{
			name:  "auto is valid",
			theme: "auto",
			want:  true,
		},
		{
			name:  "invalid theme",
			theme: "invalid",
			want:  false,
		},
		{
			name:  "empty string is invalid",
			theme: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := statuspage.ValidTheme(tt.theme); got != tt.want {
				t.Errorf("ValidTheme(%q) = %v, want %v", tt.theme, got, tt.want)
			}
		})
	}
}

func TestStyleHelpers(t *testing.T) {
	tests := []struct {
		name     string
		style    string
		wantFunc func() string
	}{
		{
			name:     "info style",
			style:    "info",
			wantFunc: statuspage.StyleInfo,
		},
		{
			name:     "warning style",
			style:    "warning",
			wantFunc: statuspage.StyleWarning,
		},
		{
			name:     "danger style",
			style:    "danger",
			wantFunc: statuspage.StyleDanger,
		},
		{
			name:     "primary style",
			style:    "primary",
			wantFunc: statuspage.StylePrimary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wantFunc(); got != tt.style {
				t.Errorf("style helper returned %q, want %q", got, tt.style)
			}
		})
	}
}

func TestValidIncidentStyle(t *testing.T) {
	tests := []struct {
		name  string
		style string
		want  bool
	}{
		{
			name:  "info is valid",
			style: "info",
			want:  true,
		},
		{
			name:  "warning is valid",
			style: "warning",
			want:  true,
		},
		{
			name:  "danger is valid",
			style: "danger",
			want:  true,
		},
		{
			name:  "primary is valid",
			style: "primary",
			want:  true,
		},
		{
			name:  "invalid style",
			style: "invalid",
			want:  false,
		},
		{
			name:  "empty string is invalid",
			style: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := statuspage.ValidIncidentStyle(tt.style); got != tt.want {
				t.Errorf("ValidIncidentStyle(%q) = %v, want %v", tt.style, got, tt.want)
			}
		})
	}
}

func TestAnalyticsTypeHelpers(t *testing.T) {
	tests := []struct {
		name          string
		analyticsType string
		wantFunc      func() string
	}{
		{
			name:          "google analytics type",
			analyticsType: "google",
			wantFunc:      statuspage.AnalyticsTypeGoogle,
		},
		{
			name:          "umami analytics type",
			analyticsType: "umami",
			wantFunc:      statuspage.AnalyticsTypeUmami,
		},
		{
			name:          "plausible analytics type",
			analyticsType: "plausible",
			wantFunc:      statuspage.AnalyticsTypePlausible,
		},
		{
			name:          "matomo analytics type",
			analyticsType: "matomo",
			wantFunc:      statuspage.AnalyticsTypeMatomo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.wantFunc(); got != tt.analyticsType {
				t.Errorf("analytics type helper returned %q, want %q", got, tt.analyticsType)
			}
		})
	}
}

func TestValidAnalyticsType(t *testing.T) {
	tests := []struct {
		name          string
		analyticsType *string
		want          bool
	}{
		{
			name:          "nil is valid",
			analyticsType: nil,
			want:          true,
		},
		{
			name:          "google is valid",
			analyticsType: ptr.To("google"),
			want:          true,
		},
		{
			name:          "umami is valid",
			analyticsType: ptr.To("umami"),
			want:          true,
		},
		{
			name:          "plausible is valid",
			analyticsType: ptr.To("plausible"),
			want:          true,
		},
		{
			name:          "matomo is valid",
			analyticsType: ptr.To("matomo"),
			want:          true,
		},
		{
			name:          "invalid analytics type",
			analyticsType: ptr.To("invalid"),
			want:          false,
		},
		{
			name:          "empty string is invalid",
			analyticsType: ptr.To(""),
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := statuspage.ValidAnalyticsType(tt.analyticsType); got != tt.want {
				t.Errorf("ValidAnalyticsType(%v) = %v, want %v", tt.analyticsType, got, tt.want)
			}
		})
	}
}
