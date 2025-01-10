package task_handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/duck_wars/internal/game_data"
	"github.com/yazmeyaa/duck_wars/internal/models"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	"github.com/yazmeyaa/duck_wars/pkg/worker"
)

type RewardItem struct {
	ItemID   string
	Quantity int
	Title    string
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
			item, exist := gd.Items[reward.ItemID]
			var title string = reward.ItemID
			if exist {
				title = item.Name
			}
			q := rand.Intn(reward.Quantity.Max-reward.Quantity.Min) + reward.Quantity.Min
			rewardItems = append(rewardItems, RewardItem{ItemID: reward.ItemID, Quantity: q, Title: title})
		}
	}
	return rewardItems, true
}

func FinishAdventureHandler(bot *tgbotapi.BotAPI, logger *slog.Logger, gd *game_data.GameData, charService *character_service.CharacterService) worker.TaskHandler {
	return func(ctx context.Context, t models.Task) worker.TaskResult {
		var payload game_data.AdventureTaskPayload
		taskPayloadBytes := strings.NewReader(t.Payload)

		if err := json.NewDecoder(taskPayloadBytes).Decode(&payload); err != nil {
			return FailedTask(t, "Failed to decode task payload", logger)
		}

		logger.Info("Adventure info", "id", payload.AdventureID, "userID", payload.CharacterID)
		adventure, exists := gd.Adventures[payload.AdventureID]

		if !exists {
			return FailedTask(t, "Provided adventure in task is not exists in game data files", logger)
		}

		isSuccess := rand.Float32() <= adventure.SuccessChance

		var rewardItems []RewardItem
		var experience int

		if isSuccess {
			items, ok := calculateRewards(gd, payload.AdventureID)
			rewardItems = items
			if !ok {
				return FailedTask(t, "Failed to calculate rewards", logger)
			}

			xp := gd.AdventureRewards[payload.AdventureID].Experience
			experience = xp.Min + rand.Intn(xp.Max-xp.Min)

			err := charService.AddExperience(payload.CharacterID, float64(experience))
			if err != nil {
				logger.Error("Add XP error", "error", err.Error())
				return FailedTask(t, "Failed to add xp", logger)
			}
		}

		var sb strings.Builder

		if isSuccess {
			sb.WriteString(adventure.SuccessMessages[rand.Intn(len(adventure.SuccessMessages)-1)])
			sb.WriteRune('\n')
			sb.WriteString(fmt.Sprintf("You got XP: %d\n\n", experience))

			if len(rewardItems) > 0 {
				sb.WriteString("You got some items:\n")
				for _, item := range rewardItems {
					sb.WriteString(fmt.Sprintf("%dx %s\n", item.Quantity, item.Title))
				}
			}
		} else {
			sb.WriteString(adventure.FailureMessages[rand.Intn(len(adventure.FailureMessages)-1)])
		}

		msg := tgbotapi.NewMessage(payload.UserID, sb.String())
		bot.Send(msg)

		logger.Info("Items to reward", "len", len(rewardItems))

		return SuccessTask(t)
	}
}
