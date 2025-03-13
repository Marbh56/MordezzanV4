package models

import (
	"time"
)

type MagicItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	ItemType    string    `json:"item_type"` // 'Rod', 'Wand', or 'Staff'
	Description string    `json:"description"`
	Charges     *int      `json:"charges,omitempty"`
	Cost        float64   `json:"cost"`
	Weight      int       `json:"weight"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateMagicItemInput struct {
	Name        string  `json:"name"`
	ItemType    string  `json:"item_type"` // 'Rod', 'Wand', or 'Staff'
	Description string  `json:"description"`
	Charges     *int    `json:"charges,omitempty"`
	Cost        float64 `json:"cost"`
	Weight      int     `json:"weight"`
}

type UpdateMagicItemInput struct {
	Name        string  `json:"name"`
	ItemType    string  `json:"item_type"` // 'Rod', 'Wand', or 'Staff'
	Description string  `json:"description"`
	Charges     *int    `json:"charges,omitempty"`
	Cost        float64 `json:"cost"`
	Weight      int     `json:"weight"`
}

func (i *CreateMagicItemInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}

	if i.ItemType == "" {
		return NewValidationError("item_type", "Item type cannot be empty")
	}

	// Validate item_type is one of the allowed values
	if !IsValidItemType(i.ItemType) {
		return NewValidationError("item_type", "Item type must be 'Rod', 'Wand', or 'Staff'")
	}

	if i.Description == "" {
		return NewValidationError("description", "Description cannot be empty")
	}

	if i.Charges != nil && *i.Charges < 0 {
		return NewValidationError("charges", "Charges cannot be negative")
	}

	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}

	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}

	return nil
}

func (i *UpdateMagicItemInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}

	if i.ItemType == "" {
		return NewValidationError("item_type", "Item type cannot be empty")
	}

	// Validate item_type is one of the allowed values
	if !IsValidItemType(i.ItemType) {
		return NewValidationError("item_type", "Item type must be 'Rod', 'Wand', or 'Staff'")
	}

	if i.Description == "" {
		return NewValidationError("description", "Description cannot be empty")
	}

	if i.Charges != nil && *i.Charges < 0 {
		return NewValidationError("charges", "Charges cannot be negative")
	}

	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}

	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}

	return nil
}

func IsValidItemType(itemType string) bool {
	return itemType == "Rod" || itemType == "Wand" || itemType == "Staff"
}
