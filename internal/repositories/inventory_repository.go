package repositories

import (
	"context"
	"database/sql"
	"errors"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type InventoryRepository interface {
	GetInventory(ctx context.Context, id int64) (*models.Inventory, error)
	GetInventoryByCharacter(ctx context.Context, characterID int64) (*models.Inventory, error)
	ListInventories(ctx context.Context) ([]*models.Inventory, error)
	CreateInventory(ctx context.Context, input *models.CreateInventoryInput) (int64, error)
	UpdateInventory(ctx context.Context, id int64, input *models.UpdateInventoryInput) error
	DeleteInventory(ctx context.Context, id int64) error

	GetInventoryItems(ctx context.Context, inventoryID int64) ([]models.InventoryItem, error)
	GetInventoryItem(ctx context.Context, id int64) (*models.InventoryItem, error)
	GetInventoryItemsByType(ctx context.Context, inventoryID int64, itemType string) ([]models.InventoryItem, error)
	GetInventoryItemByTypeAndItemID(ctx context.Context, inventoryID int64, itemType string, itemID int64) (*models.InventoryItem, error)
	AddInventoryItem(ctx context.Context, inventoryID int64, input *models.AddItemInput) (int64, error)
	UpdateInventoryItem(ctx context.Context, id int64, input *models.UpdateItemInput) error
	RemoveInventoryItem(ctx context.Context, id int64) error
	RemoveAllInventoryItems(ctx context.Context, inventoryID int64) error

	GetEquippedItems(ctx context.Context, inventoryID int64) ([]models.InventoryItem, error)
	GetItemsBySlot(ctx context.Context, inventoryID int64, slot string) ([]models.InventoryItem, error)

	UpdateInventoryWeight(ctx context.Context, id int64, weight float64) error
	RecalculateInventoryWeight(ctx context.Context, id int64) error
}

type SQLCInventoryRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

type ItemSummary struct {
	ID        int64  `json:"id"`
	ItemType  string `json:"item_type"`
	ItemID    int64  `json:"item_id"`
	Name      string `json:"name"`
	TwoHanded bool   `json:"two_handed,omitempty"`
}

type EquipmentStatus struct {
	EquippedSlots  map[string]ItemSummary `json:"equipped_slots"`
	AvailableSlots []string               `json:"available_slots"`
}

func NewSQLCInventoryRepository(db *sql.DB) *SQLCInventoryRepository {
	return &SQLCInventoryRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCInventoryRepository) GetInventory(ctx context.Context, id int64) (*models.Inventory, error) {
	inventory, err := r.q.GetInventory(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("inventory", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	items, err := r.GetInventoryItems(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.Inventory{
		ID:            inventory.ID,
		CharacterID:   inventory.CharacterID,
		MaxWeight:     inventory.MaxWeight,
		CurrentWeight: inventory.CurrentWeight,
		Items:         items,
		CreatedAt:     inventory.CreatedAt,
		UpdatedAt:     inventory.UpdatedAt,
	}, nil
}

func (r *SQLCInventoryRepository) GetInventoryByCharacter(ctx context.Context, characterID int64) (*models.Inventory, error) {
	inventory, err := r.q.GetInventoryByCharacter(ctx, characterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("inventory for character", characterID)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	items, err := r.GetInventoryItems(ctx, inventory.ID)
	if err != nil {
		return nil, err
	}

	return &models.Inventory{
		ID:            inventory.ID,
		CharacterID:   inventory.CharacterID,
		MaxWeight:     inventory.MaxWeight,
		CurrentWeight: inventory.CurrentWeight,
		Items:         items,
		CreatedAt:     inventory.CreatedAt,
		UpdatedAt:     inventory.UpdatedAt,
	}, nil
}

func (r *SQLCInventoryRepository) ListInventories(ctx context.Context) ([]*models.Inventory, error) {
	inventories, err := r.q.ListInventories(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]*models.Inventory, len(inventories))
	for i, inv := range inventories {
		items, err := r.GetInventoryItems(ctx, inv.ID)
		if err != nil {
			return nil, err
		}

		result[i] = &models.Inventory{
			ID:            inv.ID,
			CharacterID:   inv.CharacterID,
			MaxWeight:     inv.MaxWeight,
			CurrentWeight: inv.CurrentWeight,
			Items:         items,
			CreatedAt:     inv.CreatedAt,
			UpdatedAt:     inv.UpdatedAt,
		}
	}

	return result, nil
}

func (r *SQLCInventoryRepository) CreateInventory(ctx context.Context, input *models.CreateInventoryInput) (int64, error) {
	// Check if inventory already exists for this character
	_, err := r.q.GetInventoryByCharacter(ctx, input.CharacterID)
	if err == nil {
		return 0, apperrors.NewValidationError("character_id", "Inventory already exists for this character")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, apperrors.NewDatabaseError(err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)
	result, err := qtx.CreateInventory(ctx, sqlcdb.CreateInventoryParams{
		CharacterID: input.CharacterID,
		MaxWeight:   input.MaxWeight,
	})
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	if err := tx.Commit(); err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	return id, nil
}

func (r *SQLCInventoryRepository) UpdateInventory(ctx context.Context, id int64, input *models.UpdateInventoryInput) error {
	// First verify the inventory exists
	_, err := r.GetInventory(ctx, id)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	defer tx.Rollback()

	// Prepare update parameters
	var maxWeight, currentWeight sql.NullFloat64
	if input.MaxWeight != nil {
		maxWeight.Float64 = *input.MaxWeight
		maxWeight.Valid = true
	}
	if input.CurrentWeight != nil {
		currentWeight.Float64 = *input.CurrentWeight
		currentWeight.Valid = true
	}

	// Update inventory with timestamp
	_, err = tx.ExecContext(ctx, `
        UPDATE inventories 
        SET max_weight = COALESCE(?, max_weight),
            current_weight = COALESCE(?, current_weight),
            updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `,
		maxWeight,
		currentWeight,
		id)

	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

func (r *SQLCInventoryRepository) DeleteInventory(ctx context.Context, id int64) error {
	_, err := r.GetInventory(ctx, id)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// Delete all items in the inventory first
	err = qtx.RemoveAllInventoryItems(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	// Then delete the inventory itself
	err = qtx.DeleteInventory(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	if err := tx.Commit(); err != nil {
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

func (r *SQLCInventoryRepository) GetInventoryItems(ctx context.Context, inventoryID int64) ([]models.InventoryItem, error) {
	items, err := r.q.GetInventoryItems(ctx, inventoryID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]models.InventoryItem, len(items))
	for i, item := range items {
		result[i] = models.InventoryItem{
			ID:          item.ID,
			InventoryID: item.InventoryID,
			ItemType:    item.ItemType,
			ItemID:      item.ItemID,
			Quantity:    int(item.Quantity),
			IsEquipped:  item.IsEquipped,
			Notes:       item.Notes.String,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	return result, nil
}

func (r *SQLCInventoryRepository) GetInventoryItem(ctx context.Context, id int64) (*models.InventoryItem, error) {
	item, err := r.q.GetInventoryItem(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("inventory item", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	return &models.InventoryItem{
		ID:          item.ID,
		InventoryID: item.InventoryID,
		ItemType:    item.ItemType,
		ItemID:      item.ItemID,
		Quantity:    int(item.Quantity),
		IsEquipped:  item.IsEquipped,
		Notes:       item.Notes.String,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}

func (r *SQLCInventoryRepository) GetInventoryItemsByType(ctx context.Context, inventoryID int64, itemType string) ([]models.InventoryItem, error) {
	items, err := r.q.GetInventoryItemsByType(ctx, sqlcdb.GetInventoryItemsByTypeParams{
		InventoryID: inventoryID,
		ItemType:    itemType,
	})
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]models.InventoryItem, len(items))
	for i, item := range items {
		result[i] = models.InventoryItem{
			ID:          item.ID,
			InventoryID: item.InventoryID,
			ItemType:    item.ItemType,
			ItemID:      item.ItemID,
			Quantity:    int(item.Quantity),
			IsEquipped:  item.IsEquipped,
			Notes:       item.Notes.String,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	return result, nil
}

func (r *SQLCInventoryRepository) GetInventoryItemByTypeAndItemID(ctx context.Context, inventoryID int64, itemType string, itemID int64) (*models.InventoryItem, error) {
	item, err := r.q.GetInventoryItemByTypeAndItemID(ctx, sqlcdb.GetInventoryItemByTypeAndItemIDParams{
		InventoryID: inventoryID,
		ItemType:    itemType,
		ItemID:      itemID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("inventory item", itemID)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	return &models.InventoryItem{
		ID:          item.ID,
		InventoryID: item.InventoryID,
		ItemType:    item.ItemType,
		ItemID:      item.ItemID,
		Quantity:    int(item.Quantity),
		IsEquipped:  item.IsEquipped,
		Notes:       item.Notes.String,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}

func (r *SQLCInventoryRepository) AddInventoryItem(ctx context.Context, inventoryID int64, input *models.AddItemInput) (int64, error) {
	// Create a nullable string for notes
	notesParam := sql.NullString{
		String: input.Notes,
		Valid:  input.Notes != "",
	}

	// Create a nullable string for slot
	slotParam := sql.NullString{
		String: input.Slot,
		Valid:  input.Slot != "",
	}

	params := sqlcdb.AddInventoryItemParams{
		InventoryID: inventoryID,
		ItemType:    input.ItemType,
		ItemID:      input.ItemID,
		Quantity:    int64(input.Quantity),
		IsEquipped:  input.IsEquipped,
		Slot:        slotParam, // Add the slot parameter here
		Notes:       notesParam,
	}

	result, err := r.q.AddInventoryItem(ctx, params)
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	return id, nil
}

func (r *SQLCInventoryRepository) UpdateInventoryItem(ctx context.Context, id int64, input *models.UpdateItemInput) error {
	// Initialize parameters
	var quantity sql.NullInt64
	var isEquipped sql.NullBool
	var notes sql.NullString
	var slot sql.NullString

	if input.Quantity != nil {
		quantity.Int64 = int64(*input.Quantity)
		quantity.Valid = true
	}
	if input.IsEquipped != nil {
		isEquipped.Bool = *input.IsEquipped
		isEquipped.Valid = true
	}
	if input.Notes != nil {
		notes.String = *input.Notes
		notes.Valid = true
	}
	if input.Slot != nil {
		slot.String = *input.Slot
		slot.Valid = true
	}

	logger.Debug("Updating inventory item %d with: quantity=%v (valid=%v), isEquipped=%v (valid=%v), slot=%q (valid=%v), notes=%q (valid=%v)",
		id,
		quantity.Int64, quantity.Valid,
		isEquipped.Bool, isEquipped.Valid,
		slot.String, slot.Valid,
		notes.String, notes.Valid)

	// Use direct SQL instead of SQLC's function to avoid parameter issues
	query := `
		UPDATE inventory_items
		SET 
			quantity = CASE WHEN ? THEN ? ELSE quantity END,
			is_equipped = CASE WHEN ? THEN ? ELSE is_equipped END,
			slot = CASE WHEN ? THEN ? ELSE slot END,
			notes = CASE WHEN ? THEN ? ELSE notes END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	// Execute the query with all parameters
	_, err := r.db.ExecContext(ctx, query,
		quantity.Valid, quantity.Int64,
		isEquipped.Valid, isEquipped.Bool,
		slot.Valid, slot.String,
		notes.Valid, notes.String,
		id)

	if err != nil {
		logger.Error("Failed to update inventory item: %v", err)
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

func (r *SQLCInventoryRepository) RemoveInventoryItem(ctx context.Context, id int64) error {
	_, err := r.GetInventoryItem(ctx, id)
	if err != nil {
		return err
	}

	err = r.q.RemoveInventoryItem(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

func (r *SQLCInventoryRepository) RemoveAllInventoryItems(ctx context.Context, inventoryID int64) error {
	err := r.q.RemoveAllInventoryItems(ctx, inventoryID)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

func (r *SQLCInventoryRepository) GetEquippedItems(ctx context.Context, inventoryID int64) ([]models.InventoryItem, error) {
	items, err := r.q.GetEquippedItems(ctx, inventoryID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]models.InventoryItem, len(items))
	for i, item := range items {
		result[i] = models.InventoryItem{
			ID:          item.ID,
			InventoryID: item.InventoryID,
			ItemType:    item.ItemType,
			ItemID:      item.ItemID,
			Quantity:    int(item.Quantity),
			IsEquipped:  item.IsEquipped,
			Slot:        item.Slot.String,
			Notes:       item.Notes.String,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	return result, nil
}

func (r *SQLCInventoryRepository) GetItemsBySlot(ctx context.Context, inventoryID int64, slot string) ([]models.InventoryItem, error) {
	params := sqlcdb.GetItemsBySlotParams{
		InventoryID: inventoryID,
		Slot: sql.NullString{
			String: slot,
			Valid:  slot != "",
		},
	}

	items, err := r.q.GetItemsBySlot(ctx, params)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]models.InventoryItem, len(items))
	for i, item := range items {
		result[i] = models.InventoryItem{
			ID:          item.ID,
			InventoryID: item.InventoryID,
			ItemType:    item.ItemType,
			ItemID:      item.ItemID,
			Quantity:    int(item.Quantity),
			IsEquipped:  item.IsEquipped,
			Slot:        item.Slot.String,
			Notes:       item.Notes.String,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}

	return result, nil
}

func (r *SQLCInventoryRepository) UpdateInventoryWeight(ctx context.Context, id int64, weight float64) error {
	err := r.q.UpdateInventoryWeight(ctx, sqlcdb.UpdateInventoryWeightParams{
		CurrentWeight: weight,
		ID:            id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCInventoryRepository) RecalculateInventoryWeight(ctx context.Context, id int64) error {
	err := r.q.RecalculateInventoryWeight(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}
