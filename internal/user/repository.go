package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/yazmeyaa/duck_wars/pkg/logger"
)

var (
	ErrorUserNotFound error = errors.New("user not found")
)

type UserRepository interface {
	Create(data *CreateUserDTO) (*User, error)
	Delete(user *User) error
	DeleteById(id uint32) error
	Update(user *User, data *UpdateUserDTO) error
	GetByUsername(username string) (*User, error)
	GetById(id uint32) (*User, error)
}

type userRepository struct {
	log *slog.Logger
	db  *sql.DB
}

// GetById implements UserRepository.
func (u userRepository) GetById(id uint32) (*User, error) {
	const fn string = "UserRepository.GetById"
	user := User{}

	query := `SELECT "id", "username", "telegramId", "password", "gmail", "created_at"
	FROM users
	WHERE id=$1`

	row := u.db.QueryRow(query, id)
	if row == nil {
		return nil, ErrorUserNotFound
	}
	err := u.bindUser(row, &user, fn)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByUsername implements UserRepository.
func (u userRepository) GetByUsername(username string) (*User, error) {
	const fn string = "UserRepository.GetById"
	user := User{}

	query := `SELECT "id", "username", "telegramId", "password", "gmail", "created_at"
	FROM users
	WHERE username=$1`

	row := u.db.QueryRow(query, username)
	if row == nil {
		return nil, ErrorUserNotFound
	}

	err := u.bindUser(row, &user, fn)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u userRepository) Create(data *CreateUserDTO) (*User, error) {
	const fn string = "UserRepository.Create"

	user := User{}

	if data.Username == "" || data.Password == "" {
		return nil, fmt.Errorf("missed required values in CreateUserDTO. Username: %s, Password: %s", data.Username, data.Password)
	}

	query := `INSERT INTO users ("username", "password")
	VALUES ($1, $2)
	RETURNING "id", "username", "password", "created_at";`

	row := u.db.QueryRow(
		query,
		data.Username,
		data.Password,
	)

	if row == nil {
		return nil, errors.New("row is nil")
	}

	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		logger.LogError(u.log, fn, "Error while binding values to user.", err)
		return nil, err
	}

	return &user, nil
}

func (u userRepository) Delete(user *User) error {
	const fn string = "UserRepository.Delete"

	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	query := `DELETE FROM users WHERE id=$1`
	result, err := u.db.Exec(query, user.ID)
	if err != nil {
		logger.LogError(u.log, fn, "Error while deleting user.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.LogError(u.log, fn, "Error while checking affected rows.", err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id: %d", user.ID)
	}

	return nil
}

func (u userRepository) DeleteById(id uint32) error {
	return u.Delete(&User{ID: id})
}

func (u userRepository) Update(user *User, data *UpdateUserDTO) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	if user.ID == 0 {
		return fmt.Errorf("user id cannot be 0")
	}

	query := `UPDATE users SET`
	params := []interface{}{}
	setClauses := []string{}

	if data.Username != nil {
		setClauses = append(setClauses, `username = $1`)
		params = append(params, *data.Username)
	}
	if data.Password != nil {
		setClauses = append(setClauses, `password = $2`)
		params = append(params, *data.Password)
	}
	if data.TelegramID != nil {
		setClauses = append(setClauses, `telegramId = $3`)
		params = append(params, *data.TelegramID)
	}
	if data.Gmail != nil {
		setClauses = append(setClauses, `gmail = $4`)
		params = append(params, *data.Gmail)
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query += " " + strings.Join(setClauses, ", ") + " WHERE id = $" + strconv.Itoa(len(params)+1)
	params = append(params, user.ID)

	_, err := u.db.Exec(query, params...)
	if err != nil {
		logger.LogError(u.log, "UserRepository.Update", "Error while updating user.", err)
		return err
	}

	return nil
}

func (u userRepository) bindUser(row *sql.Row, user *User, fn string) error {
	var telegramID sql.NullInt64
	var gmail sql.NullString

	if err := row.Scan(&user.ID, &user.Username, &user.TelegramID, &user.Password, &gmail, &user.CreatedAt); err != nil {
		logger.LogError(u.log, fn, "Error while binding values to user.", err)
		return err
	}

	if telegramID.Valid {
		user.TelegramID = new(int64)
		*user.TelegramID = telegramID.Int64
	}

	if gmail.Valid {
		user.Gmail = new(string)
		*user.Gmail = gmail.String
	}

	return nil
}

func NewRepository(log *slog.Logger, db *sql.DB) UserRepository {
	return userRepository{
		log: log,
		db:  db,
	}
}
