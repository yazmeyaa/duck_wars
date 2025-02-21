package command_handlers

import (
	"log/slog"

	"github.com/yazmeyaa/ton_game/internal/game_data"
	characterService "github.com/yazmeyaa/ton_game/internal/service/character"
	inventoryService "github.com/yazmeyaa/ton_game/internal/service/inventory"
	itemsService "github.com/yazmeyaa/ton_game/internal/service/items"
	userService "github.com/yazmeyaa/ton_game/internal/service/user"
	"github.com/yazmeyaa/ton_game/pkg/bot_commands_handler"
	"github.com/yazmeyaa/ton_game/pkg/scheduler"
)

func NewCommandsHandler(
	gameData *game_data.GameData,
	scheduler *scheduler.Scheduler,
	userService userService.UserService,
	characterService characterService.CharacterService,
	inventoryService inventoryService.InventoryService,
	itemsService itemsService.ItemsService,
	logger *slog.Logger,
) *bot_commands_handler.Handler {
	handler := bot_commands_handler.New()

	handler.RegisterCommand(
		"start",
		HandleStartCommand(userService, characterService, logger),
	)
	handler.RegisterCommand(
		"inventory",
		HandleInventoryCommand(gameData, userService, characterService, inventoryService, logger),
	)
	handler.RegisterCommand("craft",
		CraftItemCommand(gameData, itemsService, characterService, userService),
	)
	handler.RegisterCommand("adventure",
		HandleAdventureCommand(scheduler, gameData, userService, characterService))

	return handler
}
