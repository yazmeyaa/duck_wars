package command_handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	characterService "github.com/yazmeyaa/ton_game/internal/service/character"
	itemsService "github.com/yazmeyaa/ton_game/internal/service/items"
	userService "github.com/yazmeyaa/ton_game/internal/service/user"
	"github.com/yazmeyaa/ton_game/pkg/bot_commands_handler"
)

func validateCraftCommand(argsStr string, gd *game_data.GameData) bool {
	args := strings.Split(argsStr, " ")
	if len(args) < 1 {
		return false
	}

	recipeID := args[0]
	if recipeID == "" {
		return false
	}

	_, exist := gd.Crafts[recipeID]

	return exist
}

func extractRecipeID(argsStr string) string {
	args := strings.Split(argsStr, " ")

	return args[0]
}

func processErr(bot *tgbotapi.BotAPI, chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, err.Error())
	_, err = bot.Send(msg)
	if err != nil {
		return
	}
}

func CraftItemCommand(gd *game_data.GameData, itemsService itemsService.ItemsService, characterService characterService.CharacterService, userService userService.UserService) bot_commands_handler.HandleBotCommand {
	return func(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		argsStr := update.Message.CommandArguments()

		userID := update.FromChat().ID
		if !validateCraftCommand(argsStr, gd) {
			msg := tgbotapi.NewMessage(userID, "Not valid command usage")
			_, err := bot.Send(msg)
			if err != nil {
				return
			}
			return
		}

		recipeID := extractRecipeID(argsStr)

		user, err := userService.GetByTelegramID(ctx, userID)
		if err != nil {
			processErr(bot, userID, err)
			return
		}
		char, err := characterService.GetByUserID(user.ID)

		if err != nil {
			processErr(bot, userID, err)
			return
		}
		err = itemsService.Craft(ctx, char.ID, recipeID)
		if err != nil {
			processErr(bot, userID, err)
			return
		}

		recipe := gd.Crafts[recipeID]
		item := gd.Items[recipe.ItemID]
		msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Ok! You crafted %dx %s!", recipe.Quantity, item.Name))
		go func() {
			_, err := bot.Send(msg)
			if err != nil {
				return
			}
		}()
	}
}
