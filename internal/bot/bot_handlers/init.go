package bot_handlers

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
)

type Handler struct {
	log              *slog.Logger
	userService      user_service.UserService
	characterService *character_service.CharacterService
}

func (h *Handler) handleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		h.handleStartCommand(bot, update)
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

func New(logger *slog.Logger, userService user_service.UserService, characterService *character_service.CharacterService) *Handler {
	return &Handler{
		log:              logger,
		userService:      userService,
		characterService: characterService,
	}
}
