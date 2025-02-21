package character_service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/yazmeyaa/ton_game/internal/models"
	character_repository "github.com/yazmeyaa/ton_game/internal/repository/character"
)

type characterService struct {
	characterRepository character_repository.CharacterRepository
	logger              *slog.Logger
}

// GetCharacterState implements CharacterService.
func (c *characterService) GetCharacterState(ctx context.Context, characterID int) (*models.CharacterState, error) {
	panic("unimplemented")
}

// GetCurrentEnergy implements CharacterService.
func (c *characterService) GetCurrentEnergy(ctx context.Context, characterID int) (uint8, error) {
	panic("unimplemented")
}

// SpendCharacterEnergy implements CharacterService.
func (c *characterService) SpendCharacterEnergy(ctx context.Context, characterID int, energy uint8) error {
	panic("unimplemented")
}

type CharacterService interface {
	GetByID(id int) (*models.Character, error)
	GetByUserID(id int) (*models.Character, error)
	AddExperience(id int, xp float64) error
	CharacterTextDecsription(char *models.Character) string
	GetCharacterState(ctx context.Context, characterID int) (*models.CharacterState, error)
	SpendCharacterEnergy(ctx context.Context, characterID int, energy uint8) error
	GetCurrentEnergy(ctx context.Context, characterID int) (uint8, error)
}

func New(characterRepo character_repository.CharacterRepository, logger *slog.Logger) CharacterService {
	return &characterService{
		characterRepository: characterRepo,
		logger:              logger,
	}
}

func (c *characterService) GetByID(id int) (*models.Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return c.characterRepository.GetByID(ctx, id)
}

func (c *characterService) GetByUserID(id int) (*models.Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return c.characterRepository.GetByUserID(ctx, id)
}

func (c *characterService) AddExperience(id int, xp float64) error {
	char, err := c.GetByID(id)
	if err != nil {
		c.logger.Error("Error while get user by id", "error", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err = c.characterRepository.Update(ctx, id, character_repository.UpdateCharacterDTO{
		Experience: sql.NullFloat64{
			Valid:   true,
			Float64: char.Experience + xp,
		},
	})

	if err != nil {
		c.logger.Error("Error while update user XP", "character_id", id, "error", err.Error())
	}

	return err
}

func (c *characterService) CharacterTextDecsription(char *models.Character) string {
	var sb strings.Builder

	sb.WriteString("Name: ")
	sb.WriteString(char.Name)
	sb.WriteRune('\n')
	sb.WriteString("Experience: ")
	sb.WriteString(fmt.Sprintf("%f", char.Experience))

	return sb.String()
}
