package notification

import (
	"fmt"
)

// Matrix ...
type Matrix struct {
	Base
	MatrixDetails
}

// MatrixDetails ...
type MatrixDetails struct {
	HomeserverURL  string `json:"matrixHomeserverUrl"`
	InternalRoomID string `json:"matrixInternalRoomId"`
	AccessToken    string `json:"matrixAccessToken"`
}

// Type ...
func (m Matrix) Type() string {
	return m.MatrixDetails.Type()
}

// Type ...
func (n MatrixDetails) Type() string {
	return "matrix"
}

// String ...
func (m Matrix) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(m.Base, false), formatNotification(m.MatrixDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (m Matrix) MarshalJSON() ([]byte, error) {
	return marshalJSON(m.Base, &m.MatrixDetails)
}
