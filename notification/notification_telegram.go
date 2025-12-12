package notification

import (
	"fmt"
)

type Telegram struct {
	Base
	TelegramDetails
}

type TelegramDetails struct {
	BotToken             string `json:"telegramBotToken"`
	ChatID               string `json:"telegramChatID"`
	ServerURL            string `json:"telegramServerUrl"`
	SendSilently         bool   `json:"telegramSendSilently"`
	ProtectContent       bool   `json:"telegramProtectContent"`
	MessageThreadID      string `json:"telegramMessageThreadID"`
	UseTemplate          bool   `json:"telegramUseTemplate"`
	Template             string `json:"telegramTemplate"`
	TemplateParseMode    string `json:"telegramTemplateParseMode"`
}

func (t Telegram) Type() string {
	return t.TelegramDetails.Type()
}

func (n TelegramDetails) Type() string {
	return "telegram"
}

func (t Telegram) String() string {
	return fmt.Sprintf("%s, %s", formatNotification(t.Base, false), formatNotification(t.TelegramDetails, true))
}

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

func (t Telegram) MarshalJSON() ([]byte, error) {
	return marshalJSON(t.Base, t.TelegramDetails)
}
