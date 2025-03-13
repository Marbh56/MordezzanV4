package models

import (
	"time"
)

type Shield struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Cost            float64   `json:"cost"`
	Weight          int       `json:"weight"`
	DefenseModifier int       `json:"defense_modifier"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateShieldInput struct {
	Name            string  `json:"name"`
	Cost            float64 `json:"cost"`
	Weight          int     `json:"weight"`
	DefenseModifier int     `json:"defense_modifier"`
}

type UpdateShieldInput struct {
	Name            string  `json:"name"`
	Cost            float64 `json:"cost"`
	Weight          int     `json:"weight"`
	DefenseModifier int     `json:"defense_modifier"`
}

func (i *CreateShieldInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	if i.DefenseModifier < 0 {
		return NewValidationError("defense_modifier", "Defense modifier cannot be negative")
	}
	return nil
}

func (i *UpdateShieldInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	if i.DefenseModifier < 0 {
		return NewValidationError("defense_modifier", "Defense modifier cannot be negative")
	}
	return nil
}
