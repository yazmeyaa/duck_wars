package dispatcher_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/yazmeyaa/duck_wars/internal/config"
	"github.com/yazmeyaa/duck_wars/internal/database"
	"github.com/yazmeyaa/duck_wars/pkg/dispatcher"
	"github.com/yazmeyaa/duck_wars/pkg/models"
	"github.com/yazmeyaa/duck_wars/pkg/sheduler"
	"github.com/yazmeyaa/duck_wars/pkg/worker"
)

func handleTask(ctx context.Context, t models.Task) worker.TaskResult {
	return worker.TaskResult{
		ID:     t.ID,
		Result: models.STATUS_COMPLETED,
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
	defer db.Close()

	sheduler := sheduler.NewSheduler(db, logger)

	_, err = sheduler.SheduleTask(models.Task{
		Command:    "finish_adventure",
		Payload:    "forest 1233321",
		SheduledAt: time.Now().Add(time.Second),
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
