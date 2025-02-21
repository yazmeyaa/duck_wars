package command_handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	characterService "github.com/yazmeyaa/ton_game/internal/service/character"
	inventoryService "github.com/yazmeyaa/ton_game/internal/service/inventory"
	userService "github.com/yazmeyaa/ton_game/internal/service/user"
	"github.com/yazmeyaa/ton_game/pkg/bot_commands_handler"
)

func HandleInventoryCommand(gameData *game_data.GameData, userService userService.UserService, characterService characterService.CharacterService, inventoryService inventoryService.InventoryService, logger *slog.Logger) bot_commands_handler.HandleBotCommand {
	return func(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
		tgID := update.FromChat().ID

		getUserCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		u, err := userService.GetByTelegramID(getUserCtx, tgID)

		if err != nil {
			logger.Error("Cannot get user.", "error", err.Error())
			msg := tgbotapi.NewMessage(tgID, "Internal error.")
			bot.Send(msg)
			return
		}

		c, err := characterService.GetByUserID(u.ID)
		if err != nil {
			logger.Error("Cannot get user character.", "error", err.Error())
			msg := tgbotapi.NewMessage(tgID, "Internal error.")
			bot.Send(msg)
			return
		}
		invCtx, cancelInv := context.WithTimeout(context.Background(), time.Second*5)
		defer cancelInv()
		inventory, err := inventoryService.GetInventory(invCtx, c.ID)
		if err != nil {
			logger.Error("Cannot get user inventory.", "error", err.Error())
			msg := tgbotapi.NewMessage(tgID, "Internal error.")
			bot.Send(msg)
			return
		}

		var sb strings.Builder
		gameItems := gameData.Items

		sb.WriteString("Your inventory:\n\n")
		for _, item := range inventory {
			itemID := item.ItemID
			gameItem, exist := gameItems[itemID]
			var itemName string
			if !exist {
				itemName = item.ItemID
			} else {
				itemName = gameItem.Name
			}
			sb.WriteString(fmt.Sprintf("%dx %s", item.Quantity, itemName))
			sb.WriteRune('\n')
		}

		msg := tgbotapi.NewMessage(tgID, sb.String())
		go bot.Send(msg)
	}
}
