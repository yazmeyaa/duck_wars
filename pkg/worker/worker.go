package worker

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/yazmeyaa/duck_wars/internal/models"
)

type TaskResult struct {
	ID     string
	Result models.TaskStatus
}

type Worker struct {
	id             string
	logger         *slog.Logger
	ctx            context.Context
	cancelFunc     context.CancelFunc
	taskChan       <-chan models.Task
	taskResultChan chan<- TaskResult
}

func NewWorker(ctx context.Context, cancel context.CancelFunc, taskChan <-chan models.Task, taskResultChan chan TaskResult, logger *slog.Logger) *Worker {
	return &Worker{
		id:             uuid.NewString(),
		logger:         logger,
		ctx:            ctx,
		cancelFunc:     cancel,
		taskChan:       taskChan,
		taskResultChan: taskResultChan,
	}
}

func (w *Worker) GetID() string {
	return w.id
}

func (w *Worker) Run() {
	w.logger.Info("Start listing for tasks from worker", "workerID", w.id)

	for {
		select {
		case task, ok := <-w.taskChan:
			if !ok {
				w.logger.Info("Worker: task channel closed\n", "wokerID", w.id)
				break
			}
			w.processTask(task)
		case <-w.ctx.Done():
			w.logger.Info("Worker stopping\n", "wokerID", w.id, "error", w.ctx.Err().Error())
			return
		}
	}
}

func (w *Worker) processTask(task models.Task) {
	switch task.Command {
	case "finish_adventure":
		result := TaskResult{
			ID:     task.ID,
			Result: models.STATUS_COMPLETED,
		}
		w.taskResultChan <- result
	default:
		w.taskResultChan <- TaskResult{
			ID:     task.ID,
			Result: models.STATUS_FAILDED,
		}
	}
}

func (w *Worker) Stop() {
	w.cancelFunc()
}
