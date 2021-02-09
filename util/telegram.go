package util

import (
	"github.com/DictumMortuum/servus/config"
	tb "gopkg.in/tucnak/telebot.v2"
)

type telegramRecipient struct {
	User string
}

func (r telegramRecipient) Recipient() string {
	return r.User
}

func TelegramMessage(message string) error {
	settings := tb.Settings{
		Token: config.App.Telegram.Token,
	}

	bot, err := tb.NewBot(settings)
	if err != nil {
		return err
	}

	for _, user := range config.App.Telegram.Users {
		_, err := bot.Send(telegramRecipient{user}, message)
		if err != nil {
			return err
		}
	}

	return nil
}
