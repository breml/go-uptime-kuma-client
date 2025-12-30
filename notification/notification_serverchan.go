package notification

import (
	"fmt"
)

// ServerChan represents a ServerChan notification provider.
// ServerChan is a push notification service popular in China for sending notifications.
type ServerChan struct {
	Base
	ServerChanDetails
}

// ServerChanDetails contains the configuration fields for ServerChan notifications.
type ServerChanDetails struct {
	// SendKey is the ServerChan send key for authentication.
	SendKey string `json:"serverChanSendKey"`
}

// Type returns the notification type identifier for ServerChan.
func (s ServerChan) Type() string {
	return s.ServerChanDetails.Type()
}

// Type returns the notification type identifier for ServerChanDetails.
func (ServerChanDetails) Type() string {
	return "ServerChan"
}

// String returns a string representation of the ServerChan notification.
func (s ServerChan) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.ServerChanDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a ServerChan notification.
func (s *ServerChan) UnmarshalJSON(data []byte) error {
	detail := ServerChanDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = ServerChan{
		Base:              base,
		ServerChanDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the ServerChan notification into JSON.
func (s ServerChan) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, s.ServerChanDetails)
}
