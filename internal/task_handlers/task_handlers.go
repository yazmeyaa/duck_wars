package task_handlers

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/duck_wars/internal/game_data"
	"github.com/yazmeyaa/duck_wars/internal/models"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
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

func NewTaskHandlersMap(logger *slog.Logger, gd *game_data.GameData, bot *tgbotapi.BotAPI, charService *character_service.CharacterService) map[string]worker.TaskHandler {
	handlers := map[string]worker.TaskHandler{
		"finish_adventure": FinishAdventureHandler(bot, logger, gd, charService),
	}

	return handlers
}
