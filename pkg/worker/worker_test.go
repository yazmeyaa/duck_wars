package worker_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/yazmeyaa/duck_wars/internal/models"
	"github.com/yazmeyaa/duck_wars/pkg/worker"
)

func TestWorker(t *testing.T) {
	ch := make(chan models.Task)
	taskResultChan := make(chan worker.TaskResult)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx, cancel := context.WithCancel(context.Background())
	worker := worker.NewWorker(ctx, cancel, ch, taskResultChan, logger, map[string]worker.TaskHandler{})

	go worker.Run()

	task := models.Task{
		ID:         "123",
		Command:    "finish_adventure",
		Payload:    "123",
		TaskStatus: models.STATUS_PENDING,
		CreatedAt:  time.Now(),
		SheduledAt: time.Now().Add(time.Second * 5),
	}

	ch <- task
	res := <-taskResultChan

	if res.Result != models.STATUS_COMPLETED {
		t.FailNow()
	}
}
