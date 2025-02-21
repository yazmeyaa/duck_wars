package bot_commands_handler

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandleBotCommand func(bot *tgbotapi.BotAPI, update *tgbotapi.Update)
type commandHandlersMap map[string]HandleBotCommand

var (
	ErrHandlerNotFound = errors.New("no registered handlers for provided command")
)

type Handler struct {
	handlers commandHandlersMap
}

func (c *Handler) RegisterCommand(command string, handler HandleBotCommand) {
	if _, exist := c.handlers[command]; exist {
		panic("BotCommandHandler: command already registered")
	}

	c.handlers[command] = handler
}

func (c *Handler) Handle(command string, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	handler, exist := c.handlers[command]
	if !exist {
		return ErrHandlerNotFound
	}

	handler(bot, update)
	return nil
}

func New() *Handler {
	return &Handler{
		handlers: make(commandHandlersMap),
	}
}
