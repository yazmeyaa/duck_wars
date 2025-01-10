package task_handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"strings"

	"github.com/yazmeyaa/duck_wars/internal/game_data"
	"github.com/yazmeyaa/duck_wars/internal/models"
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

type RewardItem struct {
	ItemID   string
	Quantity int
}

func calculateRewards(gd *game_data.GameData, adventureID string) ([]RewardItem, bool) {
	adventure, exists := gd.Adventures[adventureID]
	if !exists {
		return nil, false
	}

	rewards, exists := gd.AdventureRewards[adventure.AdventureID]
	if !exists {
		return nil, false
	}

	rewardItems := make([]RewardItem, 0)
	for _, reward := range rewards.Rewards {
		if rand.Float32() <= reward.Chance {
			q := rand.Intn(reward.Quantity.Max-reward.Quantity.Min) + reward.Quantity.Min
			rewardItems = append(rewardItems, RewardItem{ItemID: reward.ItemID, Quantity: q})
		}
	}
	return rewardItems, true
}

func FinishAdventureHandler(logger *slog.Logger, gd *game_data.GameData) worker.TaskHandler {
	return func(ctx context.Context, t models.Task) worker.TaskResult {
		var payload game_data.AdventureTaskPayload
		taskPayloadBytes := strings.NewReader(t.Payload)

		if err := json.NewDecoder(taskPayloadBytes).Decode(&payload); err != nil {
			return FailedTask(t, "Failed to decode task payload", logger)
		}

		adventure, exists := gd.Adventures[payload.AdventureID]
		if !exists {
			return FailedTask(t, "Provided adventure in task is not exists in game data files", logger)
		}

		isSuccess := rand.Float32() <= adventure.SuccessChance

		var rewardItems []RewardItem

		if isSuccess {
			items, ok := calculateRewards(gd, payload.AdventureID)
			rewardItems = items
			if !ok {
				return FailedTask(t, "Failed to calculate rewards", logger)
			}
		}

		logger.Info("Items to reward", "len", len(rewardItems))

		return SuccessTask(t)
	}
}

func NewTaskHandlersMap(logger *slog.Logger, gd *game_data.GameData) map[string]worker.TaskHandler {
	handlers := map[string]worker.TaskHandler{
		"finish_adventure": FinishAdventureHandler(logger, gd),
	}

	return handlers
}
