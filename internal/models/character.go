package models

import (
	"time"
)

type Character struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Name         string    `json:"name"`
	Class        string    `json:"class"`
	Level        int       `json:"level"` // Added level field
	Strength     int       `json:"strength"`
	Dexterity    int       `json:"dexterity"`
	Constitution int       `json:"constitution"`
	Wisdom       int       `json:"wisdom"`
	Intelligence int       `json:"intelligence"`
	Charisma     int       `json:"charisma"`
	HitPoints    int       `json:"hit_points"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateCharacterInput struct {
	UserID       int64  `json:"user_id"`
	Name         string `json:"name"`
	Class        string `json:"class"`
	Level        int    `json:"level"` // Added level field
	Strength     int    `json:"strength"`
	Dexterity    int    `json:"dexterity"`
	Constitution int    `json:"constitution"`
	Wisdom       int    `json:"wisdom"`
	Intelligence int    `json:"intelligence"`
	Charisma     int    `json:"charisma"`
	HitPoints    int    `json:"hit_points"`
}

type UpdateCharacterInput struct {
	Name         string `json:"name"`
	Class        string `json:"class"`
	Level        int    `json:"level"` // Added level field
	Strength     int    `json:"strength"`
	Dexterity    int    `json:"dexterity"`
	Constitution int    `json:"constitution"`
	Wisdom       int    `json:"wisdom"`
	Intelligence int    `json:"intelligence"`
	Charisma     int    `json:"charisma"`
	HitPoints    int    `json:"hit_points"`
}

func (i *CreateCharacterInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Class == "" {
		return NewValidationError("class", "Class cannot be empty")
	}
	if i.UserID <= 0 {
		return NewValidationError("user_id", "Invalid user ID")
	}
	if i.Level < 1 {
		return NewValidationError("level", "Level must be at least 1") // Added validation for level
	}
	if i.Strength < 3 || i.Strength > 18 {
		return NewValidationError("strength", "Strength must be between 3 and 18")
	}
	if i.Dexterity < 3 || i.Dexterity > 18 {
		return NewValidationError("dexterity", "Dexterity must be between 3 and 18")
	}
	if i.Constitution < 3 || i.Constitution > 18 {
		return NewValidationError("constitution", "Constitution must be between 3 and 18")
	}
	if i.Wisdom < 3 || i.Wisdom > 18 {
		return NewValidationError("wisdom", "Wisdom must be between 3 and 18")
	}
	if i.Intelligence < 3 || i.Intelligence > 18 {
		return NewValidationError("intelligence", "Intelligence must be between 3 and 18")
	}
	if i.Charisma < 3 || i.Charisma > 18 {
		return NewValidationError("charisma", "Charisma must be between 3 and 18")
	}
	if i.HitPoints < 1 {
		return NewValidationError("hit_points", "Hit points must be positive")
	}
	return nil
}

func (i *UpdateCharacterInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Class == "" {
		return NewValidationError("class", "Class cannot be empty")
	}
	if i.Level < 1 {
		return NewValidationError("level", "Level must be at least 1") // Added validation for level
	}
	if i.Strength < 3 || i.Strength > 18 {
		return NewValidationError("strength", "Strength must be between 3 and 18")
	}
	if i.Dexterity < 3 || i.Dexterity > 18 {
		return NewValidationError("dexterity", "Dexterity must be between 3 and 18")
	}
	if i.Constitution < 3 || i.Constitution > 18 {
		return NewValidationError("constitution", "Constitution must be between 3 and 18")
	}
	if i.Wisdom < 3 || i.Wisdom > 18 {
		return NewValidationError("wisdom", "Wisdom must be between 3 and 18")
	}
	if i.Intelligence < 3 || i.Intelligence > 18 {
		return NewValidationError("intelligence", "Intelligence must be between 3 and 18")
	}
	if i.Charisma < 3 || i.Charisma > 18 {
		return NewValidationError("charisma", "Charisma must be between 3 and 18")
	}
	if i.HitPoints < 1 {
		return NewValidationError("hit_points", "Hit points must be positive")
	}
	return nil
}

func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
