package inventory_repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yazmeyaa/duck_wars/internal/models"
)

type InventoryRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{
		db: db,
	}
}

var (
	ErrZeroQuantity = errors.New("provided zero quantity")
	ErrItemNotFound = errors.New("item not found")
)

func (i *InventoryRepository) AddItem(ctx context.Context, itemID string, characterID int, quantity uint16) error {
	if quantity == 0 {
		return ErrZeroQuantity
	}

	_, err := i.db.ExecContext(
		ctx,
		"INSERT INTO inventory_items (item_id, character_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (item_id, character_id) DO UPDATE SET quantity = inventory_items.quantity + EXCLUDED.quantity",
		itemID,
		characterID,
		quantity,
	)

	return err
}

func (i *InventoryRepository) RemoveItem(ctx context.Context, itemID string, characterID int) error {
	_, err := i.db.ExecContext(
		ctx,
		"DELETE FROM inventory_items WHERE item_id=$1 AND character_id=$2",
		itemID,
		characterID,
	)

	return err
}

func (r *InventoryRepository) UpdateQuantity(ctx context.Context, characterID int, itemID string, quantityDelta int) error {
	if quantityDelta == 0 {
		return ErrZeroQuantity
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.ExecContext(
		ctx,
		`
        INSERT INTO inventory_items (character_id, item_id, quantity)
        VALUES ($1, $2, GREATEST($3, 0))
        ON CONFLICT (character_id, item_id)
        DO UPDATE SET quantity = GREATEST(inventory_items.quantity + $3, 0)`,
		characterID, itemID, quantityDelta)
	if err != nil {
		return err
	}

	// Delete item if quantity is 0
	_, err = tx.ExecContext(
		ctx,
		`
		DELETE FROM inventory_items
        WHERE character_id = $1 AND item_id = $2 AND quantity = 0`,
		characterID, itemID)
	return err
}

func (i *InventoryRepository) GetInventory(ctx context.Context, characterID int) ([]models.InventoryItem, error) {
	rows, err := i.db.QueryContext(
		ctx,
		`
		SELECT character_id, item_id, quantity
		FROM inventory_items
		WHERE character_id=$1
		`,
		characterID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InventoryItem

	for rows.Next() {
		var item models.InventoryItem

		if err := rows.Scan(&item.CharacterID, &item.ItemID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (i *InventoryRepository) GetItem(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error) {
	row := i.db.QueryRowContext(
		ctx,
		`SELECT character_id, item_id, quantity
		FROM inventory_items
		WHERE character_id = $1 AND item_id = $2`,
		characterID,
		itemID,
	)

	var item models.InventoryItem
	if err := row.Scan(&item.CharacterID, &item.ItemID, &item.Quantity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}

	return &item, nil
}
