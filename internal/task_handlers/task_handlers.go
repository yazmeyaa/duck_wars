package task_handlers

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/duck_wars/internal/game_data"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	inventory_service "github.com/yazmeyaa/duck_wars/internal/service/inventory"
	"github.com/yazmeyaa/duck_wars/pkg/models"
	"github.com/yazmeyaa/duck_wars/pkg/worker"
)

func FailedTask(t models.Task, reason string, logger *slog.Logger) worker.TaskResult {
	logger.Warn("Task failed", "taskID", t.ID, "reason", reason)
	return worker.TaskResult{
		ID:     t.ID,
		Result: models.STATUS_FAILDED,
	}
}

func SuccessTask(t models.Task) worker.TaskResult {
	return worker.TaskResult{
		ID:     t.ID,
		Result: models.STATUS_COMPLETED,
	}
}

func NewTaskHandlersMap(logger *slog.Logger, gd *game_data.GameData, bot *tgbotapi.BotAPI, charService character_service.CharacterService, inventoryService inventory_service.InventoryService) map[string]worker.TaskHandler {
	handlers := map[string]worker.TaskHandler{
		"finish_adventure": finishAdventureHandler(bot, logger, gd, charService, inventoryService),
	}

	return handlers
}
