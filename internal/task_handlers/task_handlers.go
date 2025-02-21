package task_handlers

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	characterservice "github.com/yazmeyaa/ton_game/internal/service/character"
	inventoryservice "github.com/yazmeyaa/ton_game/internal/service/inventory"
	"github.com/yazmeyaa/ton_game/pkg/models"
	"github.com/yazmeyaa/ton_game/pkg/worker"
)

func FailedTask(t models.Task, reason string, logger *slog.Logger) worker.TaskResult {
	logger.Warn("Task failed", "taskID", t.ID, "reason", reason)
	return worker.TaskResult{
		ID:     t.ID,
		Result: models.StatusFailed,
	}
}

func SuccessTask(t models.Task) worker.TaskResult {
	return worker.TaskResult{
		ID:     t.ID,
		Result: models.StatusCompleted,
	}
}

func NewTaskHandlersMap(logger *slog.Logger, gd *game_data.GameData, bot *tgbotapi.BotAPI, charService characterservice.CharacterService, inventoryService inventoryservice.InventoryService) map[string]worker.TaskHandler {
	handlers := map[string]worker.TaskHandler{
		"finish_adventure": finishAdventureHandler(bot, logger, gd, charService, inventoryService),
	}

	return handlers
}
