package models

import (
	"time"
)

type Ammo struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Cost      float64   `json:"cost"`
	Weight    int       `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAmmoInput struct {
	Name   string  `json:"name"`
	Cost   float64 `json:"cost"`
	Weight int     `json:"weight"`
}

type UpdateAmmoInput struct {
	Name   string  `json:"name"`
	Cost   float64 `json:"cost"`
	Weight int     `json:"weight"`
}

func (i *CreateAmmoInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	return nil
}

func (i *UpdateAmmoInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	return nil
}
