package user_repository

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/yazmeyaa/duck_wars/internal/models"
)

const dbName = "users"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrTelegramIDAlreadyExist = errors.New("user with same telegram id already exists")
)

type UserRepository interface {
	Create(ctx context.Context, data CreateUserDTO) (int, error)
	Delete(ctx context.Context, id int) error
	// Update(ctx context.Context, userID int, data UpdateUserDTO) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByTelegramID(ctx context.Context, id int64) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

// GetByTelegramID implements UserRepository.
func (u *userRepository) GetByTelegramID(ctx context.Context, id int64) (*models.User, error) {
	row := sq.Select(
		"id",
		"telegram_id",
		"updated_at",
		"created_at",
	).
		From(dbName).
		Where(sq.Eq{"telegram_id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).
		QueryRowContext(ctx)

	var user models.User

	if err := row.Scan(
		&user.ID,
		&user.TelegramID,
		&user.UpdatedAt,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Create implements UserRepository.
func (u *userRepository) Create(ctx context.Context, data CreateUserDTO) (int, error) {
	row := sq.Insert(dbName).
		Columns("telegram_id").
		Values(data.TelegramID).
		Suffix("ON CONFLICT (telegram_id) DO UPDATE SET telegram_id = EXCLUDED.telegram_id").
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).
		QueryRowContext(ctx)

	var id int

	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "unique_user_telegram_id":
				return 0, ErrTelegramIDAlreadyExist
			}
		}

		return 0, err
	}

	return id, nil
}

// Delete implements UserRepository.
func (u *userRepository) Delete(ctx context.Context, id int) error {
	_, err := sq.Delete(dbName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).
		ExecContext(ctx)

	return err
}

// GetByID implements UserRepository.
func (u *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	row := sq.Select(
		"id",
		"telegram_id",
		"updated_at",
		"created_at",
	).
		From(dbName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).
		QueryRowContext(ctx)

	var user models.User

	if err := row.Scan(
		&user.ID,
		&user.TelegramID,
		&user.UpdatedAt,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func New(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
