package dispatcher

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"time"

	"github.com/yazmeyaa/ton_game/pkg/models"
	"github.com/yazmeyaa/ton_game/pkg/worker"
)

type Dispatcher struct {
	logger          *slog.Logger
	db              *sql.DB
	shutdownCtx     context.Context
	taskChan        chan models.Task
	taskResultChan  chan worker.TaskResult
	workers         map[string]*worker.Worker
	requestInterval int
}

func (d *Dispatcher) GetTaskChannel() chan models.Task {
	return d.taskChan
}

func (d *Dispatcher) GetTaskResultChannel() chan worker.TaskResult {
	return d.taskResultChan
}

func NewDispatcher(ctx context.Context, logger *slog.Logger, db *sql.DB, taskHandlers map[string]worker.TaskHandler, workersCount int) *Dispatcher {
	workers := map[string]*worker.Worker{}
	taskChannel := make(chan models.Task, 100)
	taskResultChannel := make(chan worker.TaskResult, 100)

	for range workersCount {
		ctx, cancel := context.WithCancel(ctx)
		w := worker.NewWorker(ctx, cancel, taskChannel, taskResultChannel, logger, taskHandlers)
		workers[w.GetID()] = w
	}

	return &Dispatcher{
		logger:          logger,
		db:              db,
		shutdownCtx:     ctx,
		taskChan:        taskChannel,
		taskResultChan:  taskResultChannel,
		workers:         workers,
		requestInterval: 1000,
	}
}

func (d *Dispatcher) loadTasks() ([]models.Task, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(
		"SELECT id, command, payload, status, sheduled_at, created_at FROM tasks WHERE sheduled_at <= NOW() AND status = 'pending' FOR UPDATE SKIP LOCKED",
	)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	tasks := make([]models.Task, 0)

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID,
			&task.Command,
			&task.Payload,
			&task.TaskStatus,
			&task.ScheduledAt,
			&task.CreatedAt,
		); err != nil {
			err := tx.Rollback()
			if err != nil {
				return nil, err
			}
			return nil, err
		}

		tasks = append(tasks, task)
	}

	for _, task := range tasks {
		if _, err := tx.Exec("UPDATE tasks SET status = 'processing' WHERE id = $1", task.ID); err != nil {
			err := tx.Rollback()
			if err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (d *Dispatcher) runWorkers() {
	for _, w := range d.workers {
		go w.Run()
	}
}

func (d *Dispatcher) Run() {
	d.runWorkers()
	ticker := time.NewTicker(time.Millisecond * time.Duration(d.requestInterval))

	for {
		select {
		case <-d.shutdownCtx.Done():
			d.logger.Info("Stop dispatcher...")
			ticker.Stop()
			d.shutdown()
			return
		case result := <-d.taskResultChan:
			if _, err := d.db.Exec("UPDATE tasks SET status = $1 WHERE id = $2", result.Result, result.ID); err != nil {
				d.logger.Error("Error while updating task status.", "error", err.Error())
			}
		case <-ticker.C:
			tasks, err := d.loadTasks()
			wg := sync.WaitGroup{}
			wg.Add(len(tasks))

			if err != nil {
				d.logger.Error("Error while loading tasks by scheduler", "error", err.Error())
				continue
			}

			for _, task := range tasks {
				task := task
				go func() {
					_, err := d.db.Exec("UPDATE tasks SET status='processing' WHERE id=$1", task.ID)
					if err != nil {
						d.logger.Error("Error while update task status...", "error", err.Error())
					}
					d.taskChan <- task
					wg.Done()
				}()
			}

			wg.Wait()
			if len(tasks) > 0 {
				d.logger.Info("Tasks were sent to workers", "tasks count", len(tasks))
			}
		}
	}
}

func (d *Dispatcher) shutdown() {
	close(d.taskChan)
	close(d.taskResultChan)
	for _, w := range d.workers {
		w.Stop()
	}
}
