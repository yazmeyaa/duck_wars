package telegram_bot

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/duck_wars/internal/config"
)

var (
	ErrTokenNotProvided = errors.New("telegram bot token not priveded")
)

func New(config *config.Config) (*tgbotapi.BotAPI, error) {
	if config.TelegramBot.Token == "" {
		return nil, ErrTokenNotProvided
	}
	bot, err := tgbotapi.NewBotAPI(config.TelegramBot.Token)
	if err != nil {
		return nil, err
	}

	return bot, nil
}
