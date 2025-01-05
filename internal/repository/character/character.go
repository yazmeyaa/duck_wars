package character_repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/yazmeyaa/duck_wars/internal/models"
)

const (
	dbName = "characters"
)

type CharacterRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *CharacterRepository {
	return &CharacterRepository{
		db: db,
	}
}

func (cr *CharacterRepository) Create(ctx context.Context, data CreateCharacterDTO) (id int, err error) {
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

func (cr *CharacterRepository) GetByUserID(ctx context.Context, id int) (*models.Character, error) {
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

func (cr *CharacterRepository) GetByID(ctx context.Context, id int) (*models.Character, error) {
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
