package worker_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/yazmeyaa/ton_game/pkg/models"
	"github.com/yazmeyaa/ton_game/pkg/worker"
)

func TestWorker(t *testing.T) {
	ch := make(chan models.Task)
	taskResultChan := make(chan worker.TaskResult)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx, cancel := context.WithCancel(context.Background())
	newWorker := worker.NewWorker(ctx, cancel, ch, taskResultChan, logger, map[string]worker.TaskHandler{})

	go newWorker.Run()

	task := models.Task{
		ID:          "123",
		Command:     "finish_adventure",
		Payload:     "123",
		TaskStatus:  models.StatusPending,
		CreatedAt:   time.Now(),
		ScheduledAt: time.Now().Add(time.Second * 5),
	}

	ch <- task
	res := <-taskResultChan

	if res.Result != models.StatusCompleted {
		t.FailNow()
	}
}
