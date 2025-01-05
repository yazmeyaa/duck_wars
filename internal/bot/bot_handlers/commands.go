package bot_handlers

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
)

func (h *Handler) handleStartCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	userID, err := h.userService.Create(ctx, user_service.CreateUserDTO{
		TelegramID:    update.SentFrom().ID,
		CharacterName: update.SentFrom().FirstName,
	})

	if err != nil {
		h.log.Error("Cannot create user.", "error", err.Error())
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Cannot create user...")
		bot.Send(msg)
		return
	}

	char, err := h.characterService.GetByUserID(userID)

	if err != nil {
		h.log.Error("Cannot get users character.", "error", err.Error())
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Cannot create user...")
		bot.Send(msg)
		return
	}

	msgTxt := h.characterService.CharacterTextDecsription(char)

	msg := tgbotapi.NewMessage(update.FromChat().ID, msgTxt)
	bot.Send(msg)
}
