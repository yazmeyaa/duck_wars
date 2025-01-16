package command_handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/duck_wars/internal/game_data"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	items_service "github.com/yazmeyaa/duck_wars/internal/service/items"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
	"github.com/yazmeyaa/duck_wars/pkg/bot_commands_handler"
)

func validateCraftCommand(argsStr string, gd *game_data.GameData) bool {
	args := strings.Split(argsStr, " ")
	if len(args) < 1 {
		return false
	}

	reciepeID := args[0]
	if reciepeID == "" {
		return false
	}

	_, exist := gd.Crafts[reciepeID]

	return exist
}

func extractReciepeID(argsStr string) string {
	args := strings.Split(argsStr, " ")

	return args[0]
}

func processErr(bot *tgbotapi.BotAPI, chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, err.Error())
	bot.Send(msg)
}

func CraftItemCommand(gd *game_data.GameData, itemsService items_service.ItemsService, characterService *character_service.CharacterService, userService user_service.UserService) bot_commands_handler.HandleBotCommand {
	return func(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		argsStr := update.Message.CommandArguments()

		userID := update.FromChat().ID
		if !validateCraftCommand(argsStr, gd) {
			msg := tgbotapi.NewMessage(userID, "Not valid command usage")
			bot.Send(msg)
			return
		}

		reciepeID := extractReciepeID(argsStr)

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
		err = itemsService.Craft(ctx, char.ID, reciepeID)
		if err != nil {
			processErr(bot, userID, err)
			return
		}

		reciepe := gd.Crafts[reciepeID]
		item := gd.Items[reciepe.ItemID]
		msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Ok! You crafted %dx %s!", reciepe.Quantity, item.Name))
		go bot.Send(msg)
	}
}
