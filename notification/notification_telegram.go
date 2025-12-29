package notification

import (
	"fmt"
)

// Telegram ...
type Telegram struct {
	Base
	TelegramDetails
}

// TelegramDetails ...
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

// Type ...
func (t Telegram) Type() string {
	return t.TelegramDetails.Type()
}

// Type ...
func (n TelegramDetails) Type() string {
	return "telegram"
}

// String ...
func (t Telegram) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TelegramDetails, true))
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (t Telegram) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, &t.TelegramDetails)
}
