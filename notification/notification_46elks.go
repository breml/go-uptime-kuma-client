package notification

import (
	"fmt"
)

// FortySixElks ...
type FortySixElks struct {
	Base
	FortySixElksDetails
}

// FortySixElksDetails ...
type FortySixElksDetails struct {
	Username   string `json:"elksUsername"`
	AuthToken  string `json:"elksAuthToken"`
	FromNumber string `json:"elksFromNumber"`
	ToNumber   string `json:"elksToNumber"`
}

// Type ...
func (f FortySixElks) Type() string {
	return f.FortySixElksDetails.Type()
}

// Type ...
func (FortySixElksDetails) Type() string {
	return "46elks"
}

// String ...
func (f FortySixElks) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(f.Base, false), formatNotification(f.FortySixElksDetails, true))
}

// UnmarshalJSON ...
func (f *FortySixElks) UnmarshalJSON(data []byte) error {
	detail := FortySixElksDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*f = FortySixElks{
		Base:                base,
		FortySixElksDetails: detail,
	}

	return nil
}

// MarshalJSON ...
func (f FortySixElks) MarshalJSON() ([]byte, error) {
	return marshalJSON(f.Base, &f.FortySixElksDetails)
}
