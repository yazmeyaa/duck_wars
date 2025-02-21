package dispatcher_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/yazmeyaa/ton_game/internal/config"
	"github.com/yazmeyaa/ton_game/internal/database"
	"github.com/yazmeyaa/ton_game/pkg/dispatcher"
	"github.com/yazmeyaa/ton_game/pkg/models"
	"github.com/yazmeyaa/ton_game/pkg/scheduler"
	"github.com/yazmeyaa/ton_game/pkg/worker"
)

func handleTask(_ context.Context, t models.Task) worker.TaskResult {
	return worker.TaskResult{
		ID:     t.ID,
		Result: models.StatusCompleted,
	}
}

func TestDispatcher(t *testing.T) {
	stop, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := config.New()
	db, err := database.New(cfg)
	if err != nil {
		logger.Error(err.Error())
		t.FailNow()
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(db)

	newScheduler := scheduler.NewScheduler(db, logger)

	_, err = newScheduler.ScheduleTask(models.Task{
		Command:     "finish_adventure",
		Payload:     "forest 1233321",
		ScheduledAt: time.Now().Add(time.Second),
	})

	h := map[string]worker.TaskHandler{
		"finish_adventure": handleTask,
	}

	if err != nil {
		logger.Error(err.Error())
		t.FailNow()
		return
	}
	d := dispatcher.NewDispatcher(stop, logger, db, h, 5)

	go d.Run()

	<-time.After(time.Second * 5)

}
