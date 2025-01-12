package command_handlers

import (
	"context"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
	"github.com/yazmeyaa/duck_wars/pkg/bot_commands_handler"
)

func HandleStartCommand(userService user_service.UserService, characterService *character_service.CharacterService, logger *slog.Logger) bot_commands_handler.HandleBotCommand {
	return func(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		userID, err := userService.Create(ctx, user_service.CreateUserDTO{
			TelegramID:    update.SentFrom().ID,
			CharacterName: update.SentFrom().FirstName,
		})

		if err != nil {
			logger.Error("Cannot create user.", "error", err.Error())
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Cannot create user...")
			bot.Send(msg)
			return
		}

		char, err := characterService.GetByUserID(userID)

		if err != nil {
			logger.Error("Cannot get users character.", "error", err.Error())
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Cannot create user...")
			bot.Send(msg)
			return
		}

		msgTxt := characterService.CharacterTextDecsription(char)

		msg := tgbotapi.NewMessage(update.FromChat().ID, msgTxt)
		go bot.Send(msg)
	}
}
