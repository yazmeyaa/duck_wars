package bot_updates_handler

import (
	"errors"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/ton_game/internal/bot/bot_updates_handler/command_handlers"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	characterService "github.com/yazmeyaa/ton_game/internal/service/character"
	inventoryService "github.com/yazmeyaa/ton_game/internal/service/inventory"
	itemsService "github.com/yazmeyaa/ton_game/internal/service/items"
	userService "github.com/yazmeyaa/ton_game/internal/service/user"
	"github.com/yazmeyaa/ton_game/pkg/bot_commands_handler"
	"github.com/yazmeyaa/ton_game/pkg/scheduler"
)

type BotUpdateHandler struct {
	commandsHandler *bot_commands_handler.Handler
	log             *slog.Logger
}

func (h *BotUpdateHandler) handleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	cmd := update.Message.Command()
	err := h.commandsHandler.Handle(cmd, bot, update)

	if err != nil {
		if errors.Is(err, bot_commands_handler.ErrHandlerNotFound) {
			h.log.Error("Received unexpected command.", "command", cmd)
			return
		}
		h.log.Error("Error while handling bot command", "error", err.Error())
		return
	}
}

func (h *BotUpdateHandler) HandleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	switch {
	case update.Message != nil:
		if update.Message.IsCommand() {
			h.handleCommand(bot, update)
			break
		}
	}
}

func New(
	gd *game_data.GameData,
	scheduler *scheduler.Scheduler,
	userService userService.UserService,
	characterService characterService.CharacterService,
	inventoryService inventoryService.InventoryService,
	itemsService itemsService.ItemsService,
	logger *slog.Logger,
) *BotUpdateHandler {
	commandsHandler := command_handlers.NewCommandsHandler(
		gd,
		scheduler,
		userService,
		characterService,
		inventoryService,
		itemsService,
		logger,
	)
	return &BotUpdateHandler{
		commandsHandler: commandsHandler,
		log:             logger,
	}
}
