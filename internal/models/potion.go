package models

import (
	"time"
)

type Potion struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Uses        int       `json:"uses"`
	Weight      int       `json:"weight"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreatePotionInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Uses        int    `json:"uses"`
	Weight      int    `json:"weight"`
}

type UpdatePotionInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Uses        int    `json:"uses"`
	Weight      int    `json:"weight"`
}

func (i *CreatePotionInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Description == "" {
		return NewValidationError("description", "Description cannot be empty")
	}
	if i.Uses < 1 {
		return NewValidationError("uses", "Uses must be at least 1")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	return nil
}

func (i *UpdatePotionInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Description == "" {
		return NewValidationError("description", "Description cannot be empty")
	}
	if i.Uses < 1 {
		return NewValidationError("uses", "Uses must be at least 1")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	return nil
}
