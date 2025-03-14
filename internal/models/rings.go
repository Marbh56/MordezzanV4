package models

import (
	"time"
)

type Ring struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost"`
	Weight      int       `json:"weight"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateRingInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
	Weight      int     `json:"weight"`
}

type UpdateRingInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
	Weight      int     `json:"weight"`
}

func (i *CreateRingInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Description == "" {
		return NewValidationError("description", "Description cannot be empty")
	}
	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	return nil
}

func (i *UpdateRingInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Description == "" {
		return NewValidationError("description", "Description cannot be empty")
	}
	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	return nil
}
