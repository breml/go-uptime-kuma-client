package notification

import (
	"fmt"
)

// LunaSea represents a LunaSea notification provider.
// LunaSea is a self-hosted controller for various services.
type LunaSea struct {
	Base
	LunaSeaDetails
}

// LunaSeaDetails contains the configuration fields for LunaSea notifications.
type LunaSeaDetails struct {
	// Target is the target type: "user" or "device".
	Target string `json:"lunaseaTarget"`
	// LunaSeaUserID is the LunaSea user ID (when Target is "user").
	LunaSeaUserID string `json:"lunaseaUserID"`
	// Device is the LunaSea device ID (when Target is "device").
	Device string `json:"lunaseaDevice"`
}

// Type returns the notification type identifier for LunaSea.
func (l LunaSea) Type() string {
	return l.LunaSeaDetails.Type()
}

// Type returns the notification type identifier for LunaSeaDetails.
func (n LunaSeaDetails) Type() string {
	return "lunasea"
}

// String returns a string representation of the LunaSea notification.
func (l LunaSea) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(l.Base, false), formatNotification(l.LunaSeaDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a LunaSea notification.
func (l *LunaSea) UnmarshalJSON(data []byte) error {
	detail := LunaSeaDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*l = LunaSea{
		Base:           base,
		LunaSeaDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the LunaSea notification into JSON.
func (l LunaSea) MarshalJSON() ([]byte, error) {
	return marshalJSON(l.Base, l.LunaSeaDetails)
}
