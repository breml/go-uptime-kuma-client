package notification

import (
	"fmt"
)

// Telegram represents a telegram notification.
type Telegram struct {
	Base
	TelegramDetails
}

// TelegramDetails contains telegram-specific notification configuration.
type TelegramDetails struct {
	BotToken          string `json:"telegramBotToken"`
	ChatID            string `json:"telegramChatID"`
	ServerURL         string `json:"telegramServerUrl"`
	SendSilently      bool   `json:"telegramSendSilently"`
	ProtectContent    bool   `json:"telegramProtectContent"`
	MessageThreadID   string `json:"telegramMessageThreadID"`
	UseTemplate       bool   `json:"telegramUseTemplate"`
	Template          string `json:"telegramTemplate"`
	TemplateParseMode string `json:"telegramTemplateParseMode"`
}

// Type returns the notification type.
func (t Telegram) Type() string {
	return t.TelegramDetails.Type()
}

// Type returns the notification type.
func (TelegramDetails) Type() string {
	return "telegram"
}

// String returns a string representation of the notification.
func (t Telegram) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TelegramDetails, true))
}

// UnmarshalJSON unmarshals a JSON byte slice into a notification.
func (t *Telegram) UnmarshalJSON(data []byte) error {
	detail := TelegramDetails{}
	base, err := unmarshalTo(data, &detail)
	if err != nil {
		return err
	}

	*t = Telegram{
		Base:            base,
		TelegramDetails: detail,
	}

	return nil
}

// MarshalJSON marshals a notification into a JSON byte slice.
func (t Telegram) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, &t.TelegramDetails)
}
