package worker

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/yazmeyaa/duck_wars/pkg/models"
)

type TaskResult struct {
	ID     string
	Result models.TaskStatus
}

type TaskHandler func(context.Context, models.Task) TaskResult

type Worker struct {
	id             string
	logger         *slog.Logger
	ctx            context.Context
	cancelFunc     context.CancelFunc
	taskChan       <-chan models.Task
	taskResultChan chan<- TaskResult
	handlers       map[string]TaskHandler
}

func NewWorker(ctx context.Context, cancel context.CancelFunc, taskChan <-chan models.Task, taskResultChan chan TaskResult, logger *slog.Logger, handlers map[string]TaskHandler) *Worker {
	return &Worker{
		id:             uuid.NewString(),
		logger:         logger,
		ctx:            ctx,
		cancelFunc:     cancel,
		taskChan:       taskChan,
		taskResultChan: taskResultChan,
		handlers:       handlers,
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
	handler, exists := w.handlers[task.Command]
	if !exists {
		w.logger.Warn("No handler for command", "command", task.Command)
		w.taskResultChan <- TaskResult{
			ID:     task.ID,
			Result: models.STATUS_FAILDED,
		}
		return
	}

	result := handler(w.ctx, task)
	w.taskResultChan <- result
}

func (w *Worker) Stop() {
	w.cancelFunc()
}
