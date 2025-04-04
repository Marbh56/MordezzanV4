package models

import (
	"time"
)

// InventoryItem represents a generic item in an inventory
type InventoryItem struct {
	ID          int64     `json:"id"`
	InventoryID int64     `json:"inventory_id"`
	ItemType    string    `json:"item_type"`
	ItemID      int64     `json:"item_id"`
	Quantity    int       `json:"quantity"`
	IsEquipped  bool      `json:"is_equipped"`
	Slot        string    `json:"slot,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Inventory represents a character's inventory
type Inventory struct {
	ID            int64           `json:"id"`
	CharacterID   int64           `json:"character_id"`
	MaxWeight     float64         `json:"max_weight"`
	CurrentWeight float64         `json:"current_weight"`
	Items         []InventoryItem `json:"items,omitempty"`
	Treasure      *Treasure       `json:"treasure,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// CreateInventoryInput represents input data for creating a new inventory
type CreateInventoryInput struct {
	CharacterID int64   `json:"character_id"`
	MaxWeight   float64 `json:"max_weight"`
}

// UpdateInventoryInput represents input data for updating an inventory
type UpdateInventoryInput struct {
	MaxWeight     *float64 `json:"max_weight,omitempty"`
	CurrentWeight *float64 `json:"current_weight,omitempty"`
}

// AddItemInput represents input data for adding an item to inventory
type AddItemInput struct {
	ItemType   string `json:"item_type"`
	ItemID     int64  `json:"item_id"`
	Quantity   int    `json:"quantity"`
	IsEquipped bool   `json:"is_equipped"`
	Slot       string `json:"slot,omitempty"`
	Notes      string `json:"notes,omitempty"`
}

// UpdateItemInput represents input data for updating an inventory item
type UpdateItemInput struct {
	Quantity   *int    `json:"quantity,omitempty"`
	IsEquipped *bool   `json:"is_equipped,omitempty"`
	Slot       *string `json:"slot,omitempty"`
	Notes      *string `json:"notes,omitempty"`
}

func (i *CreateInventoryInput) Validate() error {
	if i.CharacterID <= 0 {
		return NewValidationError("character_id", "Character ID must be positive")
	}
	if i.MaxWeight < 0 {
		return NewValidationError("max_weight", "Max weight cannot be negative")
	}
	return nil
}

func (i *UpdateInventoryInput) Validate() error {
	if i.MaxWeight != nil && *i.MaxWeight < 0 {
		return NewValidationError("max_weight", "Max weight cannot be negative")
	}
	if i.CurrentWeight != nil && *i.CurrentWeight < 0 {
		return NewValidationError("current_weight", "Current weight cannot be negative")
	}
	return nil
}

func (i *AddItemInput) Validate() error {
	if i.ItemType == "" {
		return NewValidationError("item_type", "Item type cannot be empty")
	}
	if i.ItemID <= 0 {
		return NewValidationError("item_id", "Item ID must be positive")
	}
	if i.Quantity <= 0 {
		return NewValidationError("quantity", "Quantity must be positive")
	}
	return nil
}

func (i *UpdateItemInput) Validate() error {
	if i.Quantity != nil && *i.Quantity < 0 {
		return NewValidationError("quantity", "Quantity cannot be negative")
	}
	return nil
}
