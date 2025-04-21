package models

import "time"

type Character struct {
	ID         int       `redis:"id"`
	UserID     int       `redis:"user_id"`
	Name       string    `redis:"name"`
	Coins      float64   `redis:"coins"`
	MaxEnergy  uint8     `redis:"max_energy"`
	Experience float64   `redis:"experience"`
	CreatedAt  time.Time `redis:"created_at"`
	UpdatedAt  time.Time `redis:"updated_at"`
}
