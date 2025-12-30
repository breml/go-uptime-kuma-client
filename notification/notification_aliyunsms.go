package notification

import (
	"fmt"
)

// AliyunSMS is a notification provider for Aliyun SMS service.
type AliyunSMS struct {
	Base
	AliyunSMSDetails
}

// AliyunSMSDetails contains the configuration fields for Aliyun SMS notifications.
type AliyunSMSDetails struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	PhoneNumber     string `json:"phonenumber"`
	SignName        string `json:"signName"`
	TemplateCode    string `json:"templateCode"`
}

// Type returns the notification type identifier.
func (a AliyunSMS) Type() string {
	return a.AliyunSMSDetails.Type()
}

// Type returns the notification type identifier for AliyunSMS.
func (AliyunSMSDetails) Type() string {
	return "AliyunSMS"
}

// String returns a string representation of the AliyunSMS notification.
func (a AliyunSMS) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(a.Base, false), formatNotification(a.AliyunSMSDetails, true))
}

// UnmarshalJSON unmarshals JSON data into an AliyunSMS notification.
func (a *AliyunSMS) UnmarshalJSON(data []byte) error {
	detail := AliyunSMSDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*a = AliyunSMS{
		Base:             base,
		AliyunSMSDetails: detail,
	}

	return nil
}

// MarshalJSON marshals an AliyunSMS notification into JSON data.
func (a AliyunSMS) MarshalJSON() ([]byte, error) {
	return marshalJSON(a.Base, a.AliyunSMSDetails)
}
