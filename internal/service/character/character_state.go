package character_service

import "github.com/yazmeyaa/duck_wars/internal/models"

type CharacterStateService interface {
	GetCharacterState(characterID int) *models.CharacterState
	SpendCharacterEnergy(characterID int, energy uint8) error
	GetCurrentEnergy(characterID int) (uint8, error)
}
