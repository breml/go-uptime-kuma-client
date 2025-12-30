package notification

import (
	"fmt"
)

// Matrix represents a matrix notification.
type Matrix struct {
	Base
	MatrixDetails
}

// MatrixDetails contains matrix-specific notification configuration.
type MatrixDetails struct {
	HomeserverURL  string `json:"matrixHomeserverUrl"`
	InternalRoomID string `json:"matrixInternalRoomId"`
	AccessToken    string `json:"matrixAccessToken"`
}

// Type returns the notification type.
func (m Matrix) Type() string {
	return m.MatrixDetails.Type()
}

// Type returns the notification type.
func (MatrixDetails) Type() string {
	return "matrix"
}

// String returns a string representation of the notification.
func (m Matrix) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(m.Base, false), formatNotification(m.MatrixDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (m *Matrix) UnmarshalJSON(data []byte) error {
	detail := MatrixDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*m = Matrix{
		Base:          base,
		MatrixDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (m Matrix) MarshalJSON() ([]byte, error) {
	return marshalJSON(m.Base, &m.MatrixDetails)
}
