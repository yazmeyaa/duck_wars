package character_service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/yazmeyaa/duck_wars/internal/models"
	character_repository "github.com/yazmeyaa/duck_wars/internal/repository/character"
)

type CharacterService struct {
	characterRepository character_repository.CharacterRepository
	logger              *slog.Logger
}

func New(characterRepo character_repository.CharacterRepository, logger *slog.Logger) *CharacterService {
	return &CharacterService{
		characterRepository: characterRepo,
		logger:              logger,
	}
}

func (c *CharacterService) GetByID(id int) (*models.Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return c.characterRepository.GetByID(ctx, id)
}

func (c *CharacterService) GetByUserID(id int) (*models.Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return c.characterRepository.GetByUserID(ctx, id)
}

func (c *CharacterService) AddExperience(id int, xp float64) error {
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

func (c *CharacterService) CharacterTextDecsription(char *models.Character) string {
	var sb strings.Builder

	sb.WriteString("Name: ")
	sb.WriteString(char.Name)
	sb.WriteRune('\n')
	sb.WriteString("Experience: ")
	sb.WriteString(fmt.Sprintf("%f", char.Experience))

	return sb.String()
}
