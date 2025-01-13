package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/yazmeyaa/duck_wars/internal/config"
	"github.com/yazmeyaa/duck_wars/internal/database"
	"github.com/yazmeyaa/duck_wars/internal/game_data"
	character_repository "github.com/yazmeyaa/duck_wars/internal/repository/character"
	inventory_repository "github.com/yazmeyaa/duck_wars/internal/repository/inventory"
	character_service "github.com/yazmeyaa/duck_wars/internal/service/character"
	"github.com/yazmeyaa/duck_wars/internal/service/inventory_service"
	"github.com/yazmeyaa/duck_wars/internal/task_handlers"
	"github.com/yazmeyaa/duck_wars/internal/telegram_bot"
	"github.com/yazmeyaa/duck_wars/pkg/dispatcher"
	"github.com/yazmeyaa/duck_wars/pkg/logger"
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
	logger := logger.New()
	db, err := database.New(cfg)
	if err != nil {
		logger.Error("[Task handler]: Error while connect to database", "error", err.Error())
		os.Exit(1)
		return
	}
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
	bot, err := telegram_bot.New(cfg)
	if err != nil {
		logger.Error("Error while create bot instance.", "error", err.Error())
		os.Exit(1)
		return
	}

	charRepo := character_repository.New(db)
	charService := character_service.New(charRepo, logger)
	inventoryRepo := inventory_repository.New(db)
	inventoryService := inventory_service.New(inventoryRepo, logger)

	h := task_handlers.NewTaskHandlersMap(logger, gd, bot, charService, inventoryService)

	dCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d := dispatcher.NewDispatcher(dCtx, logger, db, h, 10)
	go d.Run()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	logger.Info("[Task handler]: Shutdown...")
}
