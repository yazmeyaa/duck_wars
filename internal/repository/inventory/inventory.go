package inventory_repository

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/yazmeyaa/duck_wars/internal/models"
)

type SingleItemDTO struct {
	ItemID   string
	Quantity uint16
}

var (
	ErrZeroQuantity = errors.New("provided zero quantity")
	ErrItemNotFound = errors.New("item not found")
)

const (
	dbName = "inventory_items"
)

type InventoryRepository interface {
	GetItem(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error)
	GetItems(ctx context.Context, characterID, itemsIds []string) ([]models.InventoryItem, error)
	GetInventory(ctx context.Context, characterID int) ([]models.InventoryItem, error)
	UpdateQuantity(ctx context.Context, characterID int, itemID string, quantityDelta int) error
	UpdateItemsBatch(ctx context.Context, characterID int, items []SingleItemDTO) error
	RemoveItem(ctx context.Context, itemID string, characterID int) error
	AddItem(ctx context.Context, itemID string, characterID int, quantity uint16) error
	AddItemsBatch(ctx context.Context, characterID int, items []SingleItemDTO) error
	RemoveItemsBatch(ctx context.Context, characterID int, itemsToDelete []SingleItemDTO) error
	CraftItems(ctx context.Context, characterID int, itemsToRemove []SingleItemDTO, itemsToAdd []SingleItemDTO)
}

type inventoryRepositoryPostgres struct {
	db *sql.DB
}

// GetItems implements InventoryRepository.
func (i *inventoryRepositoryPostgres) GetItems(ctx context.Context, characterID []string, itemIds []string) ([]models.InventoryItem, error) {
	rows, err := sq.Select("character_id", "item_id", "quantity").
		From("inventory_items").
		Where(sq.Eq{"item_id": itemIds}).
		PlaceholderFormat(sq.Dollar).
		RunWith(i.db).
		QueryContext(ctx)

	if err != nil {
		return nil, err
	}

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

// CraftItems implements InventoryRepository.
func (i *inventoryRepositoryPostgres) CraftItems(ctx context.Context, characterID int, itemsToRemove []SingleItemDTO, itemsToAdd []SingleItemDTO) {

	panic("unimplemented")
}

// UpdateItemsBatch implements InventoryRepository.
func (i *inventoryRepositoryPostgres) UpdateItemsBatch(ctx context.Context, characterID int, items []SingleItemDTO) error {
	if len(items) == 0 {
		return nil
	}
	caseStmt := sq.Case("item_id")
	itemIDs := make([]string, 0, len(items))
	for _, item := range items {
		caseStmt = caseStmt.When(item.ItemID, item.Quantity)
		itemIDs = append(itemIDs, item.ItemID)
	}

	_, err := sq.Update("inventory_items").
		Set("quantity", caseStmt).
		Where(sq.Eq{"character_id": characterID}).
		Where(sq.Eq{"item_id": itemIDs}).
		PlaceholderFormat(sq.Dollar).
		RunWith(i.db).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	_, err = sq.Delete("inventory_items").
		Where(sq.Eq{"character_id": characterID}).
		Where(sq.Eq{"item_id": itemIDs}).
		Where("quantity = 0").
		PlaceholderFormat(sq.Dollar).
		RunWith(i.db).
		ExecContext(ctx)
	return err
}

// AddItemsBatch implements InventoryRepository.
func (i *inventoryRepositoryPostgres) AddItemsBatch(ctx context.Context, characterID int, items []SingleItemDTO) error {
	if len(items) == 0 {
		return nil
	}

	stmt := sq.Insert(dbName).
		Columns("character_id", "item_id", "quantity")

	for _, item := range items {
		stmt = stmt.Values(characterID, item.ItemID, item.Quantity)
	}

	stmt = stmt.PlaceholderFormat(sq.Dollar)
	_, err := stmt.RunWith(i.db).ExecContext(ctx)

	return err
}

// RemoveItemsBatch implements InventoryRepository.
func (i *inventoryRepositoryPostgres) RemoveItemsBatch(ctx context.Context, characterID int, itemsToDelete []SingleItemDTO) error {
	ids := make([]string, len(itemsToDelete))
	for _, item := range itemsToDelete {
		ids = append(ids, item.ItemID)
	}

	_, err := sq.Delete(dbName).
		Where(sq.Eq{"item_id": ids, "character_id": characterID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(i.db).
		ExecContext(ctx)

	return err
}

func New(db *sql.DB) InventoryRepository {
	return &inventoryRepositoryPostgres{
		db: db,
	}
}

func (i *inventoryRepositoryPostgres) AddItem(ctx context.Context, itemID string, characterID int, quantity uint16) error {
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

func (i *inventoryRepositoryPostgres) RemoveItem(ctx context.Context, itemID string, characterID int) error {
	_, err := i.db.ExecContext(
		ctx,
		"DELETE FROM inventory_items WHERE item_id=$1 AND character_id=$2",
		itemID,
		characterID,
	)

	return err
}

func (r *inventoryRepositoryPostgres) UpdateQuantity(ctx context.Context, characterID int, itemID string, quantityDelta int) error {
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

func (i *inventoryRepositoryPostgres) GetInventory(ctx context.Context, characterID int) ([]models.InventoryItem, error) {
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

func (i *inventoryRepositoryPostgres) GetItem(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error) {
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
