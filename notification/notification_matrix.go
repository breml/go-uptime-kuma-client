package notification

import (
	"fmt"
)

type Matrix struct {
	Base
	MatrixDetails
}

type MatrixDetails struct {
	HomeserverURL  string `json:"matrixHomeserverUrl"`
	InternalRoomID string `json:"matrixInternalRoomId"`
	AccessToken    string `json:"matrixAccessToken"`
}

func (m Matrix) Type() string {
	return m.MatrixDetails.Type()
}

func (n MatrixDetails) Type() string {
	return "matrix"
}

func (m Matrix) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(m.Base, false), formatNotification(m.MatrixDetails, true))
}

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

func (m Matrix) MarshalJSON() ([]byte, error) {
	return marshalJSON(m.Base, m.MatrixDetails)
}
