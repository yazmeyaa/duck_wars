package command_handlers

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	characterService "github.com/yazmeyaa/ton_game/internal/service/character"
	userService "github.com/yazmeyaa/ton_game/internal/service/user"
	"github.com/yazmeyaa/ton_game/pkg/bot_commands_handler"
	"github.com/yazmeyaa/ton_game/pkg/models"
	"github.com/yazmeyaa/ton_game/pkg/scheduler"
)

func validateAdventureArguments(argsStr string) bool {
	args := strings.Split(argsStr, " ")
	if len(args) < 1 {
		return false
	}

	if args[0] == "" {
		return false
	}

	return true
}

func extractAdventureID(argsStr string) string {
	adventure := strings.Split(argsStr, " ")[0]

	return adventure
}

func HandleAdventureCommand(scheduler *scheduler.Scheduler, gd *game_data.GameData, userService userService.UserService, characterService characterService.CharacterService) bot_commands_handler.HandleBotCommand {
	return func(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		chatID := update.FromChat().ID
		cmdArgs := update.Message.CommandArguments()
		if !validateAdventureArguments(cmdArgs) {
			processErr(bot, chatID, errors.New("wrong arguments provided"))
			return
		}
		user, err := userService.GetByTelegramID(ctx, chatID)
		if err != nil {
			processErr(bot, chatID, err)
			return
		}

		char, err := characterService.GetByUserID(user.ID)
		if err != nil {
			processErr(bot, chatID, err)
			return
		}

		adventureID := extractAdventureID(cmdArgs)
		adventure, exist := gd.Adventures[adventureID]
		if !exist {
			processErr(bot, chatID, errors.New("adventure not exist"))
			return
		}
		payload := game_data.AdventureTaskPayload{
			AdventureID: adventureID,
			UserID:      chatID,
			CharacterID: char.ID,
		}

		payloadBytes, err := json.Marshal(&payload)
		if err != nil {
			processErr(bot, chatID, err)
			return
		}
		task := models.Task{
			Command:     "finish_adventure",
			Payload:     string(payloadBytes),
			ScheduledAt: time.Now().Add(time.Second * time.Duration(adventure.Duration)),
		}
		_, err = scheduler.ScheduleTask(task)
		if err != nil {
			processErr(bot, chatID, err)
			return
		}
	}
}
