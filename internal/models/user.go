package models

import "time"

type User struct {
	ID          int
	TelegramID  int64
	CharacterID int
	Username    string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}
