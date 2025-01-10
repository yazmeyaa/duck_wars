package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/yazmeyaa/duck_wars/internal/game_data"
	"github.com/yazmeyaa/duck_wars/internal/task_handlers"
	"github.com/yazmeyaa/duck_wars/pkg/config"
	"github.com/yazmeyaa/duck_wars/pkg/database"
	"github.com/yazmeyaa/duck_wars/pkg/dispatcher"
	"github.com/yazmeyaa/duck_wars/pkg/logger"
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
	h := task_handlers.NewTaskHandlersMap(logger, gd)

	dCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d := dispatcher.NewDispatcher(dCtx, logger, db, h, 10)
	go d.Run()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	logger.Info("[Task handler]: Shutdown...")
}
