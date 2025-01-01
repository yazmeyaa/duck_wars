package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/yazmeyaa/duck_wars/pkg/config"
)

func New(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Database.DSN())

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
