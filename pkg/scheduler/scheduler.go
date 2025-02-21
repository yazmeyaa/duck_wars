package scheduler

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/yazmeyaa/ton_game/pkg/models"
)

type ScheduleTask struct {
	Command     string
	Payload     string
	ScheduledAt time.Time
}

type Scheduler struct {
	db     *sql.DB
	logger *slog.Logger
}

func (s *Scheduler) ScheduleTask(task models.Task) (id string, err error) {
	row := s.db.QueryRow("INSERT INTO tasks (command, payload, sheduled_at, status) VALUES ($1, $2, $3, 'pending') RETURNING id",
		task.Command,
		task.Payload,
		task.ScheduledAt,
	)
	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func NewScheduler(db *sql.DB, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		db:     db,
		logger: logger,
	}
}
