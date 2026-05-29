package notification

import (
	"fmt"
)

// JiraServiceManagement represents a Jira Service Management notification.
type JiraServiceManagement struct {
	Base
	JiraServiceManagementDetails
}

// JiraServiceManagementDetails contains Jira Service Management-specific notification configuration.
type JiraServiceManagementDetails struct {
	CloudID  string `json:"jsmCloudId"`
	Email    string `json:"jsmEmail"`
	APIToken string `json:"jsmApiToken"`
	Priority int    `json:"jsmPriority"`
}

// Type returns the notification type.
func (j JiraServiceManagement) Type() string {
	return j.JiraServiceManagementDetails.Type()
}

// Type returns the notification type.
func (JiraServiceManagementDetails) Type() string {
	return "JiraServiceManagement"
}

// String returns a string representation of the notification.
func (j JiraServiceManagement) String() string {
	masked := j.JiraServiceManagementDetails
	if masked.APIToken != "" {
		masked.APIToken = "***"
	}

	return fmt.Sprintf(
		"%s, %s",
		formatNotification(j.Base, false),
		formatNotification(masked, true),
	)
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (j *JiraServiceManagement) UnmarshalJSON(data []byte) error {
	detail := JiraServiceManagementDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*j = JiraServiceManagement{
		Base:                         base,
		JiraServiceManagementDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (j JiraServiceManagement) MarshalJSON() ([]byte, error) {
	return marshalJSON(j.Base, &j.JiraServiceManagementDetails)
}
