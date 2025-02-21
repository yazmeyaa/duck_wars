package user_repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"github.com/yazmeyaa/ton_game/internal/game_data"
	"github.com/yazmeyaa/ton_game/internal/models"
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
	GetCharacterState(ctx context.Context, characterID int) (*models.CharacterState, error)
	SpendCharacterEnergy(ctx context.Context, characterID int, energy uint8) error
	GetCurrentEnergy(ctx context.Context, characterID int) (uint8, error)
}

type postgresRepository struct {
	db *sql.DB
	rc *redis.Client
	gd *game_data.GameData
}

func (r *postgresRepository) getCharacterKey(characterID int) string {
	return "char_state:" + strconv.Itoa(characterID)
}

func (r *postgresRepository) createState(ctx context.Context, hkey string) (*models.CharacterState, error) {
	// TODO: Store somewhere initial state
	res := models.CharacterState{
		Energy:        5,
		SpendEnergyAt: time.Now(),
		State:         models.STATE_IDLE,
	}

	_, err := r.rc.HSet(ctx, hkey, &res).Result()

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetCharacterState implements CharacterStateService.
func (r *postgresRepository) GetCharacterState(ctx context.Context, characterID int) (*models.CharacterState, error) {
	k := r.getCharacterKey(characterID)

	cmd := r.rc.HGetAll(ctx, k)
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	if len(cmd.Val()) == 0 {
		return r.createState(ctx, k)
	}

	var state models.CharacterState
	if err := cmd.Scan(&state); err != nil {
		return nil, err
	}

	return &state, nil
}

// GetCurrentEnergy implements CharacterStateService.
func (r *postgresRepository) GetCurrentEnergy(ctx context.Context, characterID int) (uint8, error) {
	k := r.getCharacterKey(characterID)
	v, err := r.rc.HGet(ctx, k, "energy").Result()

	if errors.Is(err, redis.Nil) {
		x, err := r.createState(ctx, k)
		if err != nil {
			return 0, err
		}
		return x.Energy, nil
	}

	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseUint(v, 10, 8)
	return uint8(res), err
}

// SpendCharacterEnergy implements CharacterStateService.
func (r *postgresRepository) SpendCharacterEnergy(ctx context.Context, characterID int, energy uint8) error {
	panic("unimplemented")
}

// GetByTelegramID implements UserRepository.
func (u *postgresRepository) GetByTelegramID(ctx context.Context, id int64) (*models.User, error) {
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
func (u *postgresRepository) Create(ctx context.Context, data CreateUserDTO) (int, error) {
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
func (u *postgresRepository) Delete(ctx context.Context, id int) error {
	_, err := sq.Delete(dbName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).
		ExecContext(ctx)

	return err
}

// GetByID implements UserRepository.
func (u *postgresRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
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

func New(db *sql.DB, rc *redis.Client, gd *game_data.GameData) UserRepository {
	return &postgresRepository{
		db: db,
		rc: rc,
		gd: gd,
	}
}
