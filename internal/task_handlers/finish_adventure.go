package task_handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	characterservice "github.com/yazmeyaa/ton_game/internal/service/character"
	inventoryservice "github.com/yazmeyaa/ton_game/internal/service/inventory"
	"github.com/yazmeyaa/ton_game/pkg/models"
	"github.com/yazmeyaa/ton_game/pkg/worker"
)

type RewardItem struct {
	ItemID   string
	Quantity int
	Title    string
}

type AdventureResult struct {
	Message string
	Success bool
	Error   error
}

type CharacterService interface {
	AddExperience(characterID string, experience float64) error
}

type InventoryService interface {
	AddItem(ctx context.Context, itemID string, characterID string, quantity uint16) error
}

func randMinMax(min, max int) int {
	return rand.Intn(max-min) + min
}

func randChanceFloat32(chance float32) bool {
	return rand.Float32() <= chance
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
		if randChanceFloat32(reward.Chance) {
			item, exist := gd.Items[reward.ItemID]
			var title = reward.ItemID
			if exist {
				title = item.Name
			}

			q := randMinMax(reward.Quantity.Min, reward.Quantity.Max)
			rewardItems = append(rewardItems, RewardItem{ItemID: reward.ItemID, Quantity: q, Title: title})
		}
	}
	return rewardItems, true
}

func parseAdventureTaskPayload(payload string) (game_data.AdventureTaskPayload, error) {
	var result game_data.AdventureTaskPayload
	taskPayloadBytes := strings.NewReader(payload)

	if err := json.NewDecoder(taskPayloadBytes).Decode(&result); err != nil {
		return game_data.AdventureTaskPayload{}, fmt.Errorf("failed to decode payload: %w", err)
	}

	return result, nil
}

func calculateExperience(rewards game_data.AdventureRewards) int {
	return randMinMax(rewards.Experience.Min, rewards.Experience.Max)
}

func formatMessage(success bool, message string, experience int, items []RewardItem) string {
	var sb strings.Builder
	sb.WriteString(message)
	sb.WriteRune('\n')

	if success {
		sb.WriteString(fmt.Sprintf("You got XP: %d\n\n", experience))

		if len(items) > 0 {
			sb.WriteString("You got some items:\n")
			for _, item := range items {
				sb.WriteString(fmt.Sprintf("%dx %s\n", item.Quantity, item.Title))
			}
		}
	}
	return sb.String()
}

func processSuccessfulAdventure(ctx context.Context, gd *game_data.GameData, payload game_data.AdventureTaskPayload, charService characterservice.CharacterService, inventoryService inventoryservice.InventoryService) (AdventureResult, error) {
	adventure, exists := gd.Adventures[payload.AdventureID]
	if !exists {
		return AdventureResult{Error: fmt.Errorf("provided adventure in task is not exists in game data files")}, fmt.Errorf("adventure not found: %s", payload.AdventureID)
	}
	items, ok := calculateRewards(gd, payload.AdventureID)
	if !ok {
		return AdventureResult{Error: fmt.Errorf("failed to calculate rewards")}, fmt.Errorf("failed to calculate rewards for adventure: %s", payload.AdventureID)
	}
	adventureRewards := gd.AdventureRewards[payload.AdventureID]
	experience := calculateExperience(adventureRewards)
	err := charService.AddExperience(payload.CharacterID, float64(experience))
	if err != nil {
		return AdventureResult{Error: fmt.Errorf("failed to add xp: %w", err)}, fmt.Errorf("error adding experience: %w", err)
	}
	for _, item := range items {
		ctx, cancel := context.WithTimeout(ctx, time.Second*4)
		err := inventoryService.AddItem(ctx, item.ItemID, payload.CharacterID, uint16(item.Quantity))
		cancel()
		if err != nil {
			return AdventureResult{Error: fmt.Errorf("error adding item to user inventory: %w", err)}, fmt.Errorf("error adding item %s to inventory for character %d: %w", item.ItemID, payload.CharacterID, err)
		}
	}

	message := adventure.SuccessMessages[rand.Intn(len(adventure.SuccessMessages))]
	formattedMessage := formatMessage(true, message, experience, items)

	return AdventureResult{
		Message: formattedMessage,
		Success: true,
	}, nil
}

func processFailedAdventure(gd *game_data.GameData, payload game_data.AdventureTaskPayload) (AdventureResult, error) {
	adventure, exists := gd.Adventures[payload.AdventureID]
	if !exists {
		return AdventureResult{Error: fmt.Errorf("provided adventure in task is not exists in game data files")}, fmt.Errorf("adventure not found: %s", payload.AdventureID)
	}
	message := adventure.FailureMessages[rand.Intn(len(adventure.FailureMessages))]
	formattedMessage := formatMessage(false, message, 0, nil)
	return AdventureResult{
		Message: formattedMessage,
		Success: false,
	}, nil
}

func sendMessage(bot *tgbotapi.BotAPI, userID int64, message string) {
	msg := tgbotapi.NewMessage(userID, message)
	go func() {
		_, err := bot.Send(msg)
		if err != nil {
			return
		}
	}()
}

func finishAdventureHandler(bot *tgbotapi.BotAPI, logger *slog.Logger, gd *game_data.GameData, charService characterservice.CharacterService, inventoryService inventoryservice.InventoryService) worker.TaskHandler {
	return func(ctx context.Context, t models.Task) worker.TaskResult {
		payload, err := parseAdventureTaskPayload(t.Payload)

		if err != nil {
			logger.Error("Failed to parse task payload", "error", err.Error())
			return FailedTask(t, "Failed to decode task payload", logger)
		}

		adventure, exists := gd.Adventures[payload.AdventureID]
		if !exists {
			logger.Error("Provided adventure in task is not exists in game data files", "adventure_id", payload.AdventureID)
			return FailedTask(t, "Provided adventure in task is not exists in game data files", logger)
		}

		isSuccess := randChanceFloat32(adventure.SuccessChance)

		var result AdventureResult
		if isSuccess {
			result, err = processSuccessfulAdventure(ctx, gd, payload, charService, inventoryService)
			if err != nil {
				logger.Error("Error processing successful adventure", "error", err.Error())
				return FailedTask(t, "Error processing successful adventure", logger)
			}
		} else {
			result, err = processFailedAdventure(gd, payload)
			if err != nil {
				logger.Error("Error processing failed adventure", "error", err.Error())
				return FailedTask(t, "Error processing failed adventure", logger)
			}
		}

		sendMessage(bot, payload.UserID, result.Message)

		return SuccessTask(t)
	}
}
