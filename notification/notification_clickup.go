package notification

import (
	"fmt"
)

// ClickUp represents a ClickUp notification provider.
// ClickUp is a project management platform; the notification is posted as a
// Markdown message into one of its chat channels.
type ClickUp struct {
	Base
	ClickUpDetails
}

// ClickUpDetails contains the configuration fields for ClickUp notifications.
type ClickUpDetails struct {
	// Token is the personal access token, sent verbatim as the Authorization
	// header, without a Bearer prefix.
	Token string `json:"clickupToken"`
	// WorkspaceID is the id of the workspace holding the chat channel, the
	// segment of the workspace URL right after app.clickup.com/.
	WorkspaceID string `json:"clickupWorkspaceId"`
	// ChannelID is the id of the chat channel the message is posted to, the
	// segment of the channel URL right after /chat/r/.
	ChannelID string `json:"clickupChannelId"`
	// DisableURL suppresses the monitor address in the message body.
	//
	// Uptime Kuma renders this as a checkbox, which writes the key only once it
	// has been toggled. Upstream reads the flag as a plain falsy check, so an
	// absent key and an explicit false behave alike and the pointer buys no
	// semantics, only round trip fidelity: it keeps an untouched setting
	// distinguishable from an explicit false, so an edit does not send back a
	// value the server never stored.
	DisableURL *bool `json:"clickupDisableUrl,omitempty"`
}

// Type returns the notification type identifier for ClickUp.
func (c ClickUp) Type() string {
	return c.ClickUpDetails.Type()
}

// Type returns the notification type identifier for ClickUpDetails.
func (ClickUpDetails) Type() string {
	return "ClickUp"
}

// String returns a string representation of the ClickUp notification.
func (c ClickUp) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(c.Base, false), formatNotification(c.ClickUpDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a ClickUp notification.
func (c *ClickUp) UnmarshalJSON(data []byte) error {
	detail := ClickUpDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*c = ClickUp{
		Base:           base,
		ClickUpDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the ClickUp notification into JSON.
func (c ClickUp) MarshalJSON() ([]byte, error) {
	return marshalJSON(c.Base, c.ClickUpDetails)
}
