package main

import (
	"errors"
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/duck_wars/internal/bot/bot_handlers"
	user_repository "github.com/yazmeyaa/duck_wars/internal/repository/user"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
	"github.com/yazmeyaa/duck_wars/pkg/config"
	"github.com/yazmeyaa/duck_wars/pkg/db"
	"github.com/yazmeyaa/duck_wars/pkg/telegram_bot"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

func main() {
	cfg := config.New()
	bot, err := telegram_bot.New(cfg)
	if err != nil {
		logger.Error("Error while getting telegram bot instance", "error", err.Error())
		os.Exit(1)
	}

	db, err := db.New(cfg)

	if err != nil {
		logger.Error("Cannot open DB connection", "error", err.Error())
		os.Exit(1)
	}

	userRepository := user_repository.New(db)
	userService := user_service.New(logger, userRepository)

	if err != nil {
		if errors.Is(err, telegram_bot.ErrTokenNotProvided) {
			logger.Error(err.Error())
		} else {
			logger.Error("Cannot get instance of telegram bot", "error", err.Error())
		}
		os.Exit(1)
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	updatesHandler := bot_handlers.New(logger, userService)

	for update := range updates {
		updatesHandler.HandleUpdate(bot, &update)
	}
}
