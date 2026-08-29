package notification

import (
	"fmt"
)

// WxPusher represents a WxPusher notification provider.
// WxPusher delivers notifications through the standalone WxPusher app via
// the simple push (极简推送) endpoint.
type WxPusher struct {
	Base
	WxPusherDetails
}

// WxPusherDetails contains the configuration fields for WxPusher notifications.
type WxPusherDetails struct {
	// SPT is the WxPusher simple push token, for example SPT_xxxxxxxxxxxx.
	// Required: the server rejects the notification at send time if no token
	// remains. Multiple tokens are separated by commas, the server trims each
	// one, discards empty entries and delivers to all of them, batching at
	// most 10 tokens per request.
	SPT string `json:"wxpusherSPT"`
}

// Type returns the notification type identifier for WxPusher.
func (w WxPusher) Type() string {
	return w.WxPusherDetails.Type()
}

// Type returns the notification type identifier for WxPusherDetails.
func (WxPusherDetails) Type() string {
	return "WxPusher"
}

// String returns a string representation of the WxPusher notification.
func (w WxPusher) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(w.Base, false), formatNotification(w.WxPusherDetails, true))
}

// UnmarshalJSON unmarshals JSON data into a WxPusher notification.
func (w *WxPusher) UnmarshalJSON(data []byte) error {
	detail := WxPusherDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*w = WxPusher{
		Base:            base,
		WxPusherDetails: detail,
	}

	return nil
}

// MarshalJSON marshals the WxPusher notification into JSON.
func (w WxPusher) MarshalJSON() ([]byte, error) {
	return marshalJSON(w.Base, w.WxPusherDetails)
}
