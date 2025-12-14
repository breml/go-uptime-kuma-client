package notification

import (
	"fmt"
)

type FortySixElks struct {
	Base
	FortySixElksDetails
}

type FortySixElksDetails struct {
	Username   string `json:"elksUsername"`
	AuthToken  string `json:"elksAuthToken"`
	FromNumber string `json:"elksFromNumber"`
	ToNumber   string `json:"elksToNumber"`
}

func (f FortySixElks) Type() string {
	return f.FortySixElksDetails.Type()
}

func (n FortySixElksDetails) Type() string {
	return "46elks"
}

func (f FortySixElks) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(f.Base, false), formatNotification(f.FortySixElksDetails, true))
}

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

func (f FortySixElks) MarshalJSON() ([]byte, error) {
	return marshalJSON(f.Base, f.FortySixElksDetails)
}
