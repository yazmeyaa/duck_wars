package main

import (
	"errors"
	"github.com/yazmeyaa/ton_game/pkg/scheduler"
	"log/slog"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/ton_game/internal/bot/bot_updates_handler"
	"github.com/yazmeyaa/ton_game/internal/config"
	"github.com/yazmeyaa/ton_game/internal/database"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	"github.com/yazmeyaa/ton_game/internal/redis"
	characterRepository "github.com/yazmeyaa/ton_game/internal/repository/character"
	inventoryRepository "github.com/yazmeyaa/ton_game/internal/repository/inventory"
	userRepository "github.com/yazmeyaa/ton_game/internal/repository/user"
	characterService "github.com/yazmeyaa/ton_game/internal/service/character"
	inventoryService "github.com/yazmeyaa/ton_game/internal/service/inventory"
	itemsService "github.com/yazmeyaa/ton_game/internal/service/items"
	userService "github.com/yazmeyaa/ton_game/internal/service/user"
	"github.com/yazmeyaa/ton_game/internal/telegram_bot"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

var (
	resourcesPath = filepath.Join("resources")
)

func main() {
	cfg := config.New()
	db, err := database.New(cfg)

	if err != nil {
		logger.Error("Cannot open DB connection", "error", err.Error())
		os.Exit(1)
	}
	bot, err := telegram_bot.New(cfg)
	if err != nil {
		if errors.Is(err, telegram_bot.ErrTokenNotProvided) {
			logger.Error(err.Error())
		} else {
			logger.Error("Cannot get instance of telegram bot", "error", err.Error())
		}
		os.Exit(1)
		return
	}
	shed := scheduler.NewScheduler(db, logger)

	rdb, err := redis.NewRedisClient(cfg)
	if err != nil {
		logger.Error("Error while create redis client instance.", "error", err.Error())
		os.Exit(1)
		return
	}

	gd := game_data.NewGameData(logger)
	if err := gd.ReadFolderDefaultFiles(resourcesPath); err != nil {
		logger.Error("Error while reading game data files.", "error", err.Error())
		os.Exit(1)
		return
	}

	userRepo := userRepository.New(db, rdb, gd)
	characterRepo := characterRepository.NewPostgres(db)
	inventoryRepo := inventoryRepository.New(db)
	userSvc := userService.New(logger, userRepo, characterRepo)
	characterSvc := characterService.New(characterRepo, logger)
	inventorySvc := inventoryService.New(inventoryRepo, logger)
	itemsSvc := itemsService.New(inventorySvc, inventoryRepo, gd)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	updatesHandler := bot_updates_handler.New(gd, shed, userSvc, characterSvc, inventorySvc, itemsSvc, logger)

	for update := range updates {
		updatesHandler.HandleUpdate(bot, &update)
	}
}
