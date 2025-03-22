package models

import (
	"time"
)

type Container struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	MaxWeight    int       `json:"max_weight"`
	AllowedItems string    `json:"allowed_items"`
	Cost         float64   `json:"cost"`
	Weight       int       `json:"weight"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateContainerInput struct {
	Name         string  `json:"name"`
	MaxWeight    int     `json:"max_weight"`
	AllowedItems string  `json:"allowed_items"`
	Cost         float64 `json:"cost"`
}

type UpdateContainerInput struct {
	Name         string  `json:"name"`
	MaxWeight    int     `json:"max_weight"`
	AllowedItems string  `json:"allowed_items"`
	Cost         float64 `json:"cost"`
}

func (i *CreateContainerInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.MaxWeight <= 0 {
		return NewValidationError("max_weight", "Maximum weight must be positive")
	}
	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	return nil
}

func (i *UpdateContainerInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.MaxWeight <= 0 {
		return NewValidationError("max_weight", "Maximum weight must be positive")
	}
	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	return nil
}
