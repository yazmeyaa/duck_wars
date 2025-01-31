package character_service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yazmeyaa/duck_wars/internal/game_data"
	"github.com/yazmeyaa/duck_wars/internal/models"
)

type CharacterStateService interface {
	GetCharacterState(ctx context.Context, characterID int) (*models.CharacterState, error)
	SpendCharacterEnergy(ctx context.Context, characterID int, energy uint8) error
	GetCurrentEnergy(ctx context.Context, characterID int) (uint8, error)
}

type redisCharacterStateService struct {
	rc *redis.Client
	gd *game_data.GameData
}

func (r *redisCharacterStateService) getCharacterKey(characterID int) string {
	return "char_state:" + strconv.Itoa(characterID)
}

func (r *redisCharacterStateService) createState(ctx context.Context, hkey string) (*models.CharacterState, error) {
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
func (r *redisCharacterStateService) GetCharacterState(ctx context.Context, characterID int) (*models.CharacterState, error) {
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
func (r *redisCharacterStateService) GetCurrentEnergy(ctx context.Context, characterID int) (uint8, error) {
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
func (r *redisCharacterStateService) SpendCharacterEnergy(ctx context.Context, characterID int, energy uint8) error {
	panic("unimplemented")
}

func NewCharacterStateService(rc *redis.Client, gd *game_data.GameData) CharacterStateService {
	return &redisCharacterStateService{
		rc: rc,
		gd: gd,
	}
}
