package inventory_repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/yazmeyaa/ton_game/internal/models"
)

type postgresRepository struct {
	db *sql.DB
}

func New(db *sql.DB) InventoryRepository {
	return &postgresRepository{
		db: db,
	}
}

type MissedItem struct {
	ItemID   string
	Quantity uint16
}

type Queryer interface {
	QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error)
	Query(q string, args ...any) (*sql.Rows, error)
}

type Executer interface {
	Exec(q string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error)
}

// TransferItems implements InventoryRepository.
func (i *postgresRepository) TransferItems(ctx context.Context, fromID int, toID int, items []SingleItemDTO) error {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	missedItems, err := i.characterHasEnoughItems(ctx, tx, fromID, items)
	if err != nil {
		return err
	}
	if len(missedItems) > 0 {
		return ErrNotEnoughItems
	}

	if err := i.removeItems(ctx, tx, fromID, items); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	if err := i.addItems(ctx, tx, toID, items); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetItems implements InventoryRepository.
func (i *postgresRepository) GetItems(ctx context.Context, characterID int, itemIds []string) ([]models.InventoryItem, error) {
	return i.getItems(ctx, i.db, characterID, itemIds)
}

// CraftItems implements InventoryRepository.
func (i *postgresRepository) CraftItems(ctx context.Context, characterID int, itemsToRemove []SingleItemDTO, itemsToAdd []SingleItemDTO) error {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	removeItemsIds := i.getItemsIds(itemsToRemove)
	availableItems, err := i.getItems(ctx, tx, characterID, removeItemsIds)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	fmt.Printf("%+v\n", availableItems)

	if len(removeItemsIds) != len(availableItems) {
		return ErrNotEnoughItems
	}

	for _, itemToRemove := range itemsToRemove {
		isAvailable := false
		for _, availableItem := range availableItems {
			if itemToRemove.ItemID == availableItem.ItemID && availableItem.Quantity >= itemToRemove.Quantity {
				isAvailable = true
				break
			}
		}
		if !isAvailable {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return ErrNotEnoughItems
		}
	}

	if err := i.removeItems(ctx, tx, characterID, itemsToRemove); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	if err := i.addItems(ctx, tx, characterID, itemsToAdd); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// UpdateItemsBatch implements InventoryRepository.
func (i *postgresRepository) UpdateItemsBatch(ctx context.Context, characterID int, items []SingleItemDTO) error {
	if len(items) == 0 {
		return nil
	}
	caseStmt := sq.Case()
	itemIDs := make([]string, 0, len(items))
	for _, item := range items {
		caseStmt = caseStmt.When(sq.Expr("item_id = ?", item.ItemID), item.Quantity)
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

	return err
}

// AddItemsBatch implements InventoryRepository.
func (i *postgresRepository) AddItemsBatch(ctx context.Context, characterID int, items []SingleItemDTO) error {
	return i.addItems(ctx, i.db, characterID, items)
}

// RemoveItemsBatch implements InventoryRepository.
func (i *postgresRepository) RemoveItemsBatch(ctx context.Context, characterID int, itemsToDelete []SingleItemDTO) error {
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

func (i *postgresRepository) AddItem(ctx context.Context, itemID string, characterID int, quantity uint16) error {
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

func (i *postgresRepository) RemoveItem(ctx context.Context, itemID string, characterID int) error {
	_, err := i.db.ExecContext(
		ctx,
		"DELETE FROM inventory_items WHERE item_id=$1 AND character_id=$2",
		itemID,
		characterID,
	)

	return err
}

func (i *postgresRepository) UpdateQuantity(ctx context.Context, characterID int, itemID string, quantityDelta int) error {
	if quantityDelta == 0 {
		return ErrZeroQuantity
	}

	tx, err := i.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			panic(p)
		} else if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		} else {
			err := tx.Commit()
			if err != nil {
				return
			}
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

	return err
}

func (i *postgresRepository) GetInventory(ctx context.Context, characterID int) ([]models.InventoryItem, error) {
	rows, err := i.db.QueryContext(
		ctx,
		`
		SELECT character_id, item_id, quantity
		FROM inventory_items
		WHERE character_id=$1 AND quantity > 0
		`,
		characterID)

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

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

func (i *postgresRepository) GetItem(ctx context.Context, characterID int, itemID string) (*models.InventoryItem, error) {
	row := i.db.QueryRowContext(
		ctx,
		`SELECT character_id, item_id, quantity
		FROM inventory_items
		WHERE character_id = $1 AND item_id = $2 AND quantity > 0`,
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

func (i *postgresRepository) getItems(ctx context.Context, querier Queryer, characterId int, itemIds []string) ([]models.InventoryItem, error) {
	stmt := sq.Select("character_id", "item_id", "quantity").
		From("inventory_items").
		Where(
			sq.And{
				sq.Eq{
					"item_id":      itemIds,
					"character_id": characterId,
				},
				sq.Expr("quantity > 0"),
			},
		).
		PlaceholderFormat(sq.Dollar)

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := querier.QueryContext(ctx, q, args...)
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

func (i *postgresRepository) removeItems(ctx context.Context, exec Executer, characterID int, items []SingleItemDTO) error {
	if len(items) == 0 {
		return nil
	}

	ids := make([]string, len(items))
	caseStmt := sq.Case("item_id")

	for idx, item := range items {
		ids[idx] = item.ItemID
		caseStmt = caseStmt.When(
			sq.Expr("?", item.ItemID),
			sq.Expr("GREATEST(quantity - ?, 0)", item.Quantity),
		)
	}
	whereStmt := sq.Eq{
		"character_id": characterID,
		"item_id":      ids,
	}

	query, args, err := sq.Update(dbName).
		Set("quantity", caseStmt).
		Where(whereStmt).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = exec.ExecContext(ctx, query, args...)

	return err
}

func (i *postgresRepository) getItemsIds(items []SingleItemDTO) []string {
	r := make([]string, len(items))
	for idx, item := range items {
		r[idx] = item.ItemID
	}
	return r
}

func (i *postgresRepository) addItems(ctx context.Context, exec Executer, characterID int, items []SingleItemDTO) error {
	if len(items) == 0 {
		return nil
	}

	stmt := sq.Insert(dbName).
		Columns("character_id", "item_id", "quantity")

	for _, item := range items {
		stmt = stmt.Values(characterID, item.ItemID, item.Quantity)
	}

	stmt = stmt.Suffix("ON CONFLICT (item_id, character_id) DO UPDATE SET quantity = inventory_items.quantity + EXCLUDED.quantity")

	stmt = stmt.PlaceholderFormat(sq.Dollar)
	q, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	_, err = exec.ExecContext(ctx, q, args...)

	return err
}

func (i *postgresRepository) characterHasEnoughItems(ctx context.Context, q Queryer, characterID int, requiredItems []SingleItemDTO) ([]MissedItem, error) {
	ids := i.getItemsIds(requiredItems)
	availableItems, err := i.getItems(ctx, q, characterID, ids)
	if err != nil {
		return nil, err
	}

	var missedItems []MissedItem

	for _, requiredItem := range requiredItems {
		isAvailable := false
		for _, item := range availableItems {
			if item.Quantity == requiredItem.Quantity && item.Quantity >= requiredItem.Quantity {
				isAvailable = true
				break
			}
		}

		if !isAvailable {
			missedItems = append(missedItems, MissedItem{
				ItemID: requiredItem.ItemID,
			})
		}
	}

	return missedItems, nil
}
