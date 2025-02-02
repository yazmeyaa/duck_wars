package command_handlers

import (
	"log/slog"

	"github.com/yazmeyaa/duck_wars/internal/game_data"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	inventory_service "github.com/yazmeyaa/duck_wars/internal/service/inventory"
	items_service "github.com/yazmeyaa/duck_wars/internal/service/items"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
	"github.com/yazmeyaa/duck_wars/pkg/bot_commands_handler"
	"github.com/yazmeyaa/duck_wars/pkg/sheduler"
)

func NewCommandsHandler(
	gameData *game_data.GameData,
	sheduler *sheduler.Sheduler,
	userService user_service.UserService,
	characterService character_service.CharacterService,
	inventoryService inventory_service.InventoryService,
	itemsService items_service.ItemsService,
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
		HandleAdventureCommand(sheduler, gameData, userService, characterService))

	return handler
}
