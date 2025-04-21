package character_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/yazmeyaa/ton_game/internal/models"
)

const (
	dbName = "characters"
)

type CharacterRepository interface {
	Create(ctx context.Context, data CreateCharacterDTO) (id int, err error)
	GetByUserID(ctx context.Context, id int) (*models.Character, error)
	GetByID(ctx context.Context, id int) (*models.Character, error)
	Update(ctx context.Context, id int, data UpdateCharacterDTO) error
}

type postgresRepository struct {
	db     *sql.DB
	rdb    *redis.Client
	logger *slog.Logger
}

func NewPostgres(db *sql.DB, rbd *redis.Client, logger *slog.Logger) CharacterRepository {
	return &postgresRepository{
		db:     db,
		rdb:    rbd,
		logger: logger,
	}
}

func (cr *postgresRepository) Create(ctx context.Context, data CreateCharacterDTO) (id int, err error) {
	row := sq.Insert(dbName).Columns(
		"user_id",
		"name",
	).Values(
		data.UserID,
		data.Name,
	).
		Suffix("ON CONFLICT (user_id) DO UPDATE SET name = EXCLUDED.name, user_id = EXCLUDED.user_id").
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(cr.db).
		QueryRowContext(ctx)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (cr *postgresRepository) GetByUserID(ctx context.Context, id int) (*models.Character, error) {
	var char models.Character
	key := fmt.Sprintf("characters:character:user_id_%d", id)
	cmd := cr.rdb.HGetAll(ctx, key)
	res, err := cmd.Result()

	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		row := sq.Select(
			"id",
			"user_id",
			"name",
			"coins",
			"max_energy",
			"experience",
			"created_at",
			"updated_at",
		).
			From(dbName).
			Where(sq.Eq{"user_id": id}).
			PlaceholderFormat(sq.Dollar).
			RunWith(cr.db).
			QueryRowContext(ctx)

		if err := row.Scan(
			&char.ID,
			&char.UserID,
			&char.Name,
			&char.Coins,
			&char.MaxEnergy,
			&char.Experience,
			&char.CreatedAt,
			&char.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := cr.rdb.HSet(ctx, key, char).Err(); err != nil {
			cr.logger.Error("Error while saving character", "error", err.Error())
		}

	} else {
		err := cmd.Scan(&char)
		if err != nil {
			return nil, err
		}
	}

	return &char, nil
}

func (cr *postgresRepository) GetByID(ctx context.Context, id int) (*models.Character, error) {
	var char models.Character
	key := fmt.Sprintf("characters:character:%d", id)
	cmd := cr.rdb.HGetAll(ctx, key)
	res, err := cmd.Result()

	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		row := sq.Select(
			"id",
			"user_id",
			"name",
			"coins",
			"max_energy",
			"experience",
			"created_at",
			"updated_at",
		).
			From(dbName).
			Where(sq.Eq{"id": id}).
			PlaceholderFormat(sq.Dollar).
			RunWith(cr.db).
			QueryRowContext(ctx)

		if err := row.Scan(
			&char.ID,
			&char.UserID,
			&char.Name,
			&char.Coins,
			&char.MaxEnergy,
			&char.Experience,
			&char.CreatedAt,
			&char.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := cr.rdb.HSet(ctx, key, char).Err(); err != nil {
			cr.logger.Error("Error while saving character", "error", err.Error())
		}

	} else {
		err := cmd.Scan(&char)
		if err != nil {
			return nil, err
		}
	}

	return &char, nil
}

func (cr *postgresRepository) Update(ctx context.Context, id int, data UpdateCharacterDTO) error {
	expr := sq.Update(dbName)
	tx, err := cr.db.BeginTx(ctx, nil)
	if err != nil {
		cr.logger.Error("Error while starting transaction", "error", err.Error())
		return err
	}

	fieldsUpdated := false

	if data.Coins.Valid {
		expr = expr.Set("coins", data.Coins.Float64)
		fieldsUpdated = true
	}

	if data.Name.Valid {
		expr = expr.Set("name", data.Name.String)
		fieldsUpdated = true
	}

	if data.MaxEnergy.Valid {
		expr = expr.Set("max_energy", data.MaxEnergy.Int16)
		fieldsUpdated = true
	}

	if data.Experience.Valid {
		expr = expr.Set("experience", data.Experience.Float64)
		fieldsUpdated = true
	}

	if !fieldsUpdated {
		return errors.New("no fields to update")
	}

	_, err = expr.Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		ExecContext(ctx)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			cr.logger.Error("Error while rolling back transaction", "error", err.Error())
			return err
		}
		cr.logger.Error("Error while updating character", "error", err.Error())
		return err
	}

	key := fmt.Sprintf("characters:character:%d", id)
	if err := cr.rdb.Del(ctx, key).Err(); err != nil {
		cr.logger.Error("Error while deleting cache", "error", err.Error())

		if err := tx.Rollback(); err != nil {
			cr.logger.Error("Error while rolling back transaction", "error", err.Error())
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		cr.logger.Error("Error while committing transaction", "error", err.Error())
		return err
	}

	return nil
}
