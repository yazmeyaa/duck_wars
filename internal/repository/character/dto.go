package character_repository

import "database/sql"

type CreateCharacterDTO struct {
	UserID int
	Name   string
}

type UpdateCharacterDTO struct {
	Name       sql.NullString
	Coins      sql.NullFloat64
	MaxEnergy  sql.NullInt16
	Experience sql.NullFloat64
}
