package models

import (
	"time"
)

type SpellScroll struct {
	ID           int64     `json:"id"`
	SpellID      int64     `json:"spell_id"`
	SpellName    string    `json:"spell_name,omitempty"` // For API responses
	CastingLevel int       `json:"casting_level"`
	Cost         float64   `json:"cost"`
	Weight       int       `json:"weight"`
	Description  string    `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateSpellScrollInput struct {
	SpellID      int64   `json:"spell_id"`
	CastingLevel int     `json:"casting_level"`
	Cost         float64 `json:"cost"`
	Weight       int     `json:"weight"`
	Description  string  `json:"description,omitempty"`
}

type UpdateSpellScrollInput struct {
	SpellID      int64   `json:"spell_id"`
	CastingLevel int     `json:"casting_level"`
	Cost         float64 `json:"cost"`
	Weight       int     `json:"weight"`
	Description  string  `json:"description,omitempty"`
}

func (i *CreateSpellScrollInput) Validate() error {
	if i.SpellID <= 0 {
		return NewValidationError("spell_id", "Spell ID must be positive")
	}

	if i.CastingLevel <= 0 || i.CastingLevel > 9 {
		return NewValidationError("casting_level", "Casting level must be between 1 and 9")
	}

	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}

	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}

	return nil
}

func (i *UpdateSpellScrollInput) Validate() error {
	if i.SpellID <= 0 {
		return NewValidationError("spell_id", "Spell ID must be positive")
	}

	if i.CastingLevel <= 0 || i.CastingLevel > 20 {
		return NewValidationError("casting_level", "Casting level must be between 1 and 9")
	}

	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}

	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}

	return nil
}
