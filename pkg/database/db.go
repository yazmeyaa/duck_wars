package database

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/yazmeyaa/duck_wars/config"
	"github.com/yazmeyaa/duck_wars/pkg/logger"
)

func NewDB(config *config.Config, log *slog.Logger) (*sql.DB, error) {
	const fn string = "NewDB"
	dsn, err := config.Database.DSN()
	if err != nil {
		logger.LogError(log, fn, "Failed to get DSN.", err)
	}
	db, err := sql.Open("postgres", dsn)

	return db, err
}
