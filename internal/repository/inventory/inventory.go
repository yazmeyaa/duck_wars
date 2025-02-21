package inventory_repository

import (
	"context"
	"errors"

	"github.com/yazmeyaa/ton_game/internal/models"
)

type SingleItemDTO struct {
	ItemID   string
	Quantity uint16
}

var (
	ErrZeroQuantity    = errors.New("provided zero quantity")
	ErrItemNotFound    = errors.New("item not found")
	ErrNotEnoughtItems = errors.New("not enought items")
)

const (
	dbName = "inventory_items"
)

type InventoryRepository interface {
	GetItem(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error)
	GetItems(ctx context.Context, characterID int, itemsIds []string) ([]models.InventoryItem, error)
	GetInventory(ctx context.Context, characterID int) ([]models.InventoryItem, error)
	UpdateQuantity(ctx context.Context, characterID int, itemID string, quantityDelta int) error
	UpdateItemsBatch(ctx context.Context, characterID int, items []SingleItemDTO) error
	RemoveItem(ctx context.Context, itemID string, characterID int) error
	AddItem(ctx context.Context, itemID string, characterID int, quantity uint16) error
	AddItemsBatch(ctx context.Context, characterID int, items []SingleItemDTO) error
	RemoveItemsBatch(ctx context.Context, characterID int, itemsToDelete []SingleItemDTO) error
	CraftItems(ctx context.Context, characterID int, itemsToRemove []SingleItemDTO, itemsToAdd []SingleItemDTO) error
	TransferItems(ctx context.Context, fromID int, toID int, items []SingleItemDTO) error
}
