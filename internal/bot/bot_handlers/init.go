package bot_handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/duck_wars/internal/game_data"
	"github.com/yazmeyaa/duck_wars/internal/models"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
	"github.com/yazmeyaa/duck_wars/pkg/sheduler"
)

type Handler struct {
	log              *slog.Logger
	userService      user_service.UserService
	characterService *character_service.CharacterService
	sheduler         *sheduler.Sheduler
	gameData         *game_data.GameData
}

func (h *Handler) handleForest(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	fn := "handleForest"

	userID := update.FromChat().ID
	adventure, exist := h.gameData.Adventures["forest"]
	if !exist {
		h.log.Error("Adventure is nil", "fn", fn)
		msg := tgbotapi.NewMessage(userID, "Error while execute request. :(")
		bot.Send(msg)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	user, err := h.userService.GetByTelegramID(ctx, userID)
	if err != nil {
		h.log.Error("Cannot get user.", "fn", fn)
		msg := tgbotapi.NewMessage(userID, "Error while execute request. :(")
		bot.Send(msg)
		return
	}

	char, err := h.characterService.GetByUserID(user.ID)
	if err != nil {
		h.log.Error("Cannot retrieve user character.", "fn", fn, "error", err.Error())
		msg := tgbotapi.NewMessage(userID, "Error while execute request. :(")
		bot.Send(msg)
		return
	}

	adventureID := adventure.AdventureID
	payload := game_data.AdventureTaskPayload{
		UserID:      userID,
		CharacterID: char.ID,
		AdventureID: adventureID,
	}

	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		h.log.Error("JSON marshal error.", "fn", fn, "error", err.Error())
		msg := tgbotapi.NewMessage(userID, "Error while execute request. :(")
		bot.Send(msg)
		return
	}

	task := models.Task{
		Command:    "finish_adventure",
		Payload:    string(payloadBytes),
		SheduledAt: time.Now().Add(time.Second * time.Duration(adventure.Duration)),
	}

	if _, err := h.sheduler.SheduleTask(task); err != nil {
		h.log.Error("Error while shedule task.", "fn", fn, "error", err.Error())
		msg := tgbotapi.NewMessage(userID, "Error while execute request. :(")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(userID, "Ok, you go to the forest! :)")
	bot.Send(msg)
}

func (h *Handler) handleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		h.handleStartCommand(bot, update)
	case "forest":
		h.handleForest(bot, update)
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

func New(logger *slog.Logger, gd *game_data.GameData, userService user_service.UserService, characterService *character_service.CharacterService, sheduler *sheduler.Sheduler) *Handler {
	return &Handler{
		log:              logger,
		userService:      userService,
		characterService: characterService,
		sheduler:         sheduler,
		gameData:         gd,
	}
}
