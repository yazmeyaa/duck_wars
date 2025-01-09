package sheduler

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/yazmeyaa/duck_wars/internal/models"
)

type SheduleTask struct {
	Command    string
	Payload    string
	SheduledAt time.Time
}

type Sheduler struct {
	db     *sql.DB
	logger *slog.Logger
}

func (s *Sheduler) SheduleTask(task models.Task) (id string, err error) {
	row := s.db.QueryRow("INSERT INTO tasks (command, payload, sheduled_at, status) VALUES ($1, $2, $3, 'pending') RETURNING id",
		task.Command,
		task.Payload,
		task.SheduledAt,
	)
	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func NewSheduler(db *sql.DB, logger *slog.Logger) *Sheduler {
	return &Sheduler{
		db:     db,
		logger: logger,
	}
}
