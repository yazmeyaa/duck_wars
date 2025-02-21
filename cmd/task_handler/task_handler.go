package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/yazmeyaa/ton_game/internal/config"
	"github.com/yazmeyaa/ton_game/internal/database"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	characterRepository "github.com/yazmeyaa/ton_game/internal/repository/character"
	inventoryRepository "github.com/yazmeyaa/ton_game/internal/repository/inventory"
	characterService "github.com/yazmeyaa/ton_game/internal/service/character"
	inventoryService "github.com/yazmeyaa/ton_game/internal/service/inventory"
	"github.com/yazmeyaa/ton_game/internal/task_handlers"
	"github.com/yazmeyaa/ton_game/internal/telegram_bot"
	"github.com/yazmeyaa/ton_game/pkg/dispatcher"
	"github.com/yazmeyaa/ton_game/pkg/logger"
)

var (
	resourcesFilePath = filepath.Join("resources")
)

func main() {
	cfg := config.New()
	appLogger := logger.New("task_handler")
	db, err := database.New(cfg)
	if err != nil {
		appLogger.Error("[Task handler]: Error while connect to database", "error", err.Error())
		os.Exit(1)
		return
	}
	gd := game_data.NewGameData(appLogger)
	if err := gd.ReadFolderDefaultFiles(resourcesFilePath); err != nil {
		appLogger.Error("Error while reading game data files.", "error", err.Error())
		os.Exit(1)
		return
	}
	bot, err := telegram_bot.New(cfg)
	if err != nil {
		appLogger.Error("Error while create bot instance.", "error", err.Error())
		os.Exit(1)
		return
	}

	charRepo := characterRepository.NewPostgres(db)
	charService := characterService.New(charRepo, appLogger)
	inventoryRepo := inventoryRepository.New(db)
	inventorySvc := inventoryService.New(inventoryRepo, appLogger)

	h := task_handlers.NewTaskHandlersMap(appLogger, gd, bot, charService, inventorySvc)

	dCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d := dispatcher.NewDispatcher(dCtx, appLogger, db, h, 10)
	go d.Run()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	appLogger.Info("[Task handler]: Shutdown...")
}
