package main

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yazmeyaa/duck_wars/internal/bot/bot_updates_handler"
	"github.com/yazmeyaa/duck_wars/internal/config"
	"github.com/yazmeyaa/duck_wars/internal/database"
	"github.com/yazmeyaa/duck_wars/internal/game_data"
	character_repository "github.com/yazmeyaa/duck_wars/internal/repository/character"
	inventory_repository "github.com/yazmeyaa/duck_wars/internal/repository/inventory"
	user_repository "github.com/yazmeyaa/duck_wars/internal/repository/user"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	inventory_service "github.com/yazmeyaa/duck_wars/internal/service/inventory"
	items_service "github.com/yazmeyaa/duck_wars/internal/service/items"
	user_service "github.com/yazmeyaa/duck_wars/internal/service/user"
	"github.com/yazmeyaa/duck_wars/internal/telegram_bot"
	"github.com/yazmeyaa/duck_wars/pkg/sheduler"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

var (
	testDataPath              = filepath.Join("resources")
	craftFilePath             = filepath.Join(testDataPath, "crafts.json")
	itemsFilePath             = filepath.Join(testDataPath, "items.json")
	weaponAttributesFilePath  = filepath.Join(testDataPath, "weapon_attributes.json")
	adventuresFilePath        = filepath.Join(testDataPath, "adventures.json")
	adventuresRewardsFilePath = filepath.Join(testDataPath, "adventures_rewards.json")
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
		logger.Error("Error while getting telegram bot instance", "error", err.Error())
		os.Exit(1)
	}
	sheduler := sheduler.NewSheduler(db, logger)

	gd := game_data.NewGameData(logger)
	if err := gd.Read(game_data.GameDataConstructor{
		CraftFile:            craftFilePath,
		ItemsFile:            itemsFilePath,
		WeaponAttributesFile: weaponAttributesFilePath,
		AdventuresFile:       adventuresFilePath,
		AdventureRewardsFile: adventuresRewardsFilePath,
	}); err != nil {
		logger.Error("Error while reading game data files.", "error", err.Error())
		os.Exit(1)
		return
	}

	userRepository := user_repository.New(db)
	characterRepository := character_repository.New(db)
	inventoryRepository := inventory_repository.New(db)
	userService := user_service.New(logger, userRepository, characterRepository)
	characterService := character_service.New(characterRepository, logger)
	inventoryService := inventory_service.New(inventoryRepository, logger)
	itemsService := items_service.New(inventoryService, inventoryRepository, gd)

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

	updatesHandler := bot_updates_handler.New(gd, sheduler, userService, characterService, inventoryService, itemsService, logger)

	for update := range updates {
		updatesHandler.HandleUpdate(bot, &update)
	}
}
