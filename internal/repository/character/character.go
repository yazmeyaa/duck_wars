package character_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	db *sql.DB
}

func NewPostgres(db *sql.DB) CharacterRepository {
	return &postgresRepository{
		db: db,
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

	var char models.Character

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

	return &char, nil
}

func (cr *postgresRepository) GetByID(ctx context.Context, id int) (*models.Character, error) {
	row := sq.Select("id",
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

	var char models.Character

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

	return &char, nil
}

func (c *postgresRepository) Update(ctx context.Context, id int, data UpdateCharacterDTO) error {
	expr := sq.Update(dbName)

	if data.Coins.Valid {
		expr = expr.Set("coins", data.Coins.Float64)
	}

	if data.Name.Valid {
		expr = expr.Set("name", data.Name.String)
	}

	if data.MaxEnergy.Valid {
		expr = expr.Set("max_energy", data.MaxEnergy.Int16)
	}

	if data.Experience.Valid {
		fmt.Printf("EXP: %f\n", data.Experience.Float64)
		expr = expr.Set("experience", data.Experience.Float64)
	}

	if expr == sq.Update(dbName) {
		return errors.New("no fields to update")
	}

	_, err := expr.Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(c.db).
		ExecContext(ctx)

	return err
}
