package notification

import (
	"fmt"
)

// GoogleSheets represents a Google Sheets notification provider.
// Google Sheets logs notifications into a spreadsheet via a deployed
// Google Apps Script web app.
type GoogleSheets struct {
	Base
	GoogleSheetsDetails
}

// GoogleSheetsDetails contains the Google Sheets-specific configuration.
type GoogleSheetsDetails struct {
	// Note: the upstream API uses lowercase "Url" (not "URL") in the JSON key,
	// unlike the sibling GoogleChat type. Do not "fix" this to maintain wire
	// compatibility with the Uptime Kuma server.
	WebhookURL string `json:"googleSheetsWebhookUrl"`
}

// Type returns the notification type identifier for Google Sheets.
func (g GoogleSheets) Type() string {
	return g.GoogleSheetsDetails.Type()
}

// Type returns the notification type identifier for Google Sheets details.
func (GoogleSheetsDetails) Type() string {
	return "GoogleSheets"
}

// String returns a human-readable representation of the Google Sheets notification.
func (g GoogleSheets) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(g.Base, false), formatNotification(g.GoogleSheetsDetails, true))
}

// UnmarshalJSON deserializes a Google Sheets notification from JSON.
func (g *GoogleSheets) UnmarshalJSON(data []byte) error {
	detail := GoogleSheetsDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*g = GoogleSheets{
		Base:                base,
		GoogleSheetsDetails: detail,
	}

	return nil
}

// MarshalJSON serializes a Google Sheets notification to JSON.
func (g GoogleSheets) MarshalJSON() ([]byte, error) {
	return marshalJSON(g.Base, g.GoogleSheetsDetails)
}
