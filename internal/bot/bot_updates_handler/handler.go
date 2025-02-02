package bot_updates_handler

import (
	"errors"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/duck_wars/internal/bot/bot_updates_handler/command_handlers"
	"github.com/yazmeyaa/duck_wars/internal/game_data"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	inventory_service "github.com/yazmeyaa/duck_wars/internal/service/inventory"
	items_service "github.com/yazmeyaa/duck_wars/internal/service/items"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
	"github.com/yazmeyaa/duck_wars/pkg/bot_commands_handler"
	"github.com/yazmeyaa/duck_wars/pkg/sheduler"
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
			h.log.Error("Recieved unexpected command.", "command", cmd)
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
	sheduler *sheduler.Sheduler,
	userService user_service.UserService,
	characterService character_service.CharacterService,
	inventoryService inventory_service.InventoryService,
	itemsService items_service.ItemsService,
	logger *slog.Logger,
) *BotUpdateHandler {
	commandsHandler := command_handlers.NewCommandsHandler(
		gd,
		sheduler,
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
