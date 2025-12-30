package notification

import (
	"fmt"
)

// YZJ represents a YZJ (云之家/Yunzhijia) notification provider.
// YZJ is a Chinese enterprise communication and collaboration platform for sending notifications.
type YZJ struct {
	Base
	YZJDetails
}

// YZJDetails contains the configuration fields for YZJ notifications.
type YZJDetails struct {
	// WebHookURL is the YZJ webhook URL endpoint.
	WebHookURL string `json:"yzjWebHookUrl"`
	// Token is the YZJ authentication token.
	Token string `json:"yzjToken"`
}

// Type returns the notification type identifier for YZJ.
func (y YZJ) Type() string {
	return y.YZJDetails.Type()
}

// Type returns the notification type identifier for YZJDetails.
func (YZJDetails) Type() string {
	return "YZJ"
}

// String returns a string representation of the YZJ notification.
func (y YZJ) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(y.Base, false), formatNotification(y.YZJDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a YZJ notification.
func (y *YZJ) UnmarshalJSON(data []byte) error {
	detail := YZJDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*y = YZJ{
		Base:       base,
		YZJDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the YZJ notification into JSON.
func (y YZJ) MarshalJSON() ([]byte, error) {
	return marshalJSON(y.Base, y.YZJDetails)
}
