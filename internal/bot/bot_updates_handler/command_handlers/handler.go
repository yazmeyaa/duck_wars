package command_handlers

import (
	"log/slog"

	"github.com/yazmeyaa/duck_wars/internal/game_data"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	"github.com/yazmeyaa/duck_wars/internal/service/inventory_service"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
	"github.com/yazmeyaa/duck_wars/pkg/bot_commands_handler"
)

func NewCommandsHandler(
	gameData *game_data.GameData,
	userService user_service.UserService,
	characterService *character_service.CharacterService,
	inventoryService *inventory_service.InventoryService,
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

	return handler
}
