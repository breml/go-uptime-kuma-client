package notification

import (
	"fmt"
)

// SMTP ...
type SMTP struct {
	Base
	SMTPDetails
}

// SMTPDetails ...
type SMTPDetails struct {
	Host                 string `json:"smtpHost"`
	Port                 int    `json:"smtpPort"`
	Secure               bool   `json:"smtpSecure"`
	IgnoreTLSError       bool   `json:"smtpIgnoreTLSError"`
	DkimDomain           string `json:"smtpDkimDomain"`
	DkimKeySelector      string `json:"smtpDkimKeySelector"`
	DkimPrivateKey       string `json:"smtpDkimPrivateKey"`
	DkimHashAlgo         string `json:"smtpDkimHashAlgo"`
	DkimHeaderFieldNames string `json:"smtpDkimheaderFieldNames"`
	DkimSkipFields       string `json:"smtpDkimskipFields"`
	Username             string `json:"smtpUsername"`
	Password             string `json:"smtpPassword"`
	From                 string `json:"smtpFrom"`
	CC                   string `json:"smtpCC"`
	BCC                  string `json:"smtpBCC"`
	To                   string `json:"smtpTo"`
	CustomSubject        string `json:"customSubject"`
	CustomBody           string `json:"customBody"`
	HTMLBody             bool   `json:"htmlBody"`
}

// Type ...
func (s SMTP) Type() string {
	return s.SMTPDetails.Type()
}

// Type ...
func (n SMTPDetails) Type() string {
	return "smtp"
}

// String ...
func (s SMTP) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(s.Base, false), formatNotification(s.SMTPDetails, true))
}

// UnmarshalJSON ...
func (s *SMTP) UnmarshalJSON(data []byte) error {
	detail := SMTPDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*s = SMTP{
		Base:        base,
		SMTPDetails: detail,
	}

	return nil
}

// MarshalJSON ...
func (s SMTP) MarshalJSON() ([]byte, error) {
	return marshalJSON(s.Base, &s.SMTPDetails)
}
