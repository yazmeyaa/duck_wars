package bot_handlers

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
)

type Handler struct {
	log         *slog.Logger
	userService user_service.UserService
}

func (h *Handler) handleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	switch update.Message.Command() {
	default:
		fmt.Printf("%d\n", update.Message.From.ID)
		h.log.Info(update.Message.Command())
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Hi!")

		bot.Send(msg)
	}
}

func (h *Handler) HandleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	switch {
	case update.Message != nil:
		if update.Message.IsCommand() {
			h.handleCommand(bot, update)
			break
		}
	}
}

func New(logger *slog.Logger, userService user_service.UserService) *Handler {
	return &Handler{
		log:         logger,
		userService: userService,
	}
}
