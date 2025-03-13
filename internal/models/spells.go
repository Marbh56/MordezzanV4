package models

import (
	"time"
)

// Spell represents a magic spell that can be cast
type Spell struct {
	ID           int64     `json:"id"`
	CharacterID  int64     `json:"character_id"`
	Name         string    `json:"name"`
	MagLevel     int       `json:"mag_level"` // Magician level
	CryLevel     int       `json:"cry_level"` // Cryomancer level
	IllLevel     int       `json:"ill_level"` // Illusionist level
	NecLevel     int       `json:"nec_level"` // Necromancer level
	PyrLevel     int       `json:"pyr_level"` // Pyromancer level
	WchLevel     int       `json:"wch_level"` // Witch level
	ClrLevel     int       `json:"clr_level"` // Cleric level
	DrdLevel     int       `json:"drd_level"` // Druid level
	Range        string    `json:"range"`
	Duration     string    `json:"duration"`
	AreaOfEffect string    `json:"area_of_effect"`
	Components   string    `json:"components"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateSpellInput is used for creating a new spell
type CreateSpellInput struct {
	CharacterID  int64  `json:"character_id"`
	Name         string `json:"name"`
	MagLevel     int    `json:"mag_level"`
	CryLevel     int    `json:"cry_level"`
	IllLevel     int    `json:"ill_level"`
	NecLevel     int    `json:"nec_level"`
	PyrLevel     int    `json:"pyr_level"`
	WchLevel     int    `json:"wch_level"`
	ClrLevel     int    `json:"clr_level"`
	DrdLevel     int    `json:"drd_level"`
	Range        string `json:"range"`
	Duration     string `json:"duration"`
	AreaOfEffect string `json:"area_of_effect"`
	Components   string `json:"components"`
	Description  string `json:"description"`
}

// UpdateSpellInput is used for updating an existing spell
type UpdateSpellInput struct {
	Name         string `json:"name"`
	MagLevel     int    `json:"mag_level"`
	CryLevel     int    `json:"cry_level"`
	IllLevel     int    `json:"ill_level"`
	NecLevel     int    `json:"nec_level"`
	PyrLevel     int    `json:"pyr_level"`
	WchLevel     int    `json:"wch_level"`
	ClrLevel     int    `json:"clr_level"`
	DrdLevel     int    `json:"drd_level"`
	Range        string `json:"range"`
	Duration     string `json:"duration"`
	AreaOfEffect string `json:"area_of_effect"`
	Components   string `json:"components"`
	Description  string `json:"description"`
}

// Validate checks if the input for creating a spell is valid
func (i *CreateSpellInput) Validate() error {
	if i.CharacterID <= 0 {
		return NewValidationError("character_id", "Invalid character ID")
	}
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}

	// Check if at least one school has a valid level
	if i.MagLevel <= 0 && i.CryLevel <= 0 && i.IllLevel <= 0 &&
		i.NecLevel <= 0 && i.PyrLevel <= 0 && i.WchLevel <= 0 &&
		i.ClrLevel <= 0 && i.DrdLevel <= 0 {
		return NewValidationError("level", "At least one school must have a valid spell level")
	}

	// Validate individual levels if set
	if i.MagLevel < 0 || i.MagLevel > 9 {
		return NewValidationError("mag_level", "Magician level must be between 0 and 9")
	}
	if i.CryLevel < 0 || i.CryLevel > 9 {
		return NewValidationError("cry_level", "Cryomancer level must be between 0 and 9")
	}
	if i.IllLevel < 0 || i.IllLevel > 9 {
		return NewValidationError("ill_level", "Illusionist level must be between 0 and 9")
	}
	if i.NecLevel < 0 || i.NecLevel > 9 {
		return NewValidationError("nec_level", "Necromancer level must be between 0 and 9")
	}
	if i.PyrLevel < 0 || i.PyrLevel > 9 {
		return NewValidationError("pyr_level", "Pyromancer level must be between 0 and 9")
	}
	if i.WchLevel < 0 || i.WchLevel > 9 {
		return NewValidationError("wch_level", "Witch level must be between 0 and 9")
	}
	if i.ClrLevel < 0 || i.ClrLevel > 9 {
		return NewValidationError("clr_level", "Cleric level must be between 0 and 9")
	}
	if i.DrdLevel < 0 || i.DrdLevel > 9 {
		return NewValidationError("drd_level", "Druid level must be between 0 and 9")
	}

	if i.Range == "" {
		return NewValidationError("range", "Range cannot be empty")
	}
	if i.Duration == "" {
		return NewValidationError("duration", "Duration cannot be empty")
	}

	return nil
}

// Validate checks if the input for updating a spell is valid
func (i *UpdateSpellInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}

	// Check if at least one school has a valid level
	if i.MagLevel <= 0 && i.CryLevel <= 0 && i.IllLevel <= 0 &&
		i.NecLevel <= 0 && i.PyrLevel <= 0 && i.WchLevel <= 0 &&
		i.ClrLevel <= 0 && i.DrdLevel <= 0 {
		return NewValidationError("level", "At least one school must have a valid spell level")
	}

	// Validate individual levels if set
	if i.MagLevel < 0 || i.MagLevel > 9 {
		return NewValidationError("mag_level", "Magician level must be between 0 and 9")
	}
	if i.CryLevel < 0 || i.CryLevel > 9 {
		return NewValidationError("cry_level", "Cryomancer level must be between 0 and 9")
	}
	if i.IllLevel < 0 || i.IllLevel > 9 {
		return NewValidationError("ill_level", "Illusionist level must be between 0 and 9")
	}
	if i.NecLevel < 0 || i.NecLevel > 9 {
		return NewValidationError("nec_level", "Necromancer level must be between 0 and 9")
	}
	if i.PyrLevel < 0 || i.PyrLevel > 9 {
		return NewValidationError("pyr_level", "Pyromancer level must be between 0 and 9")
	}
	if i.WchLevel < 0 || i.WchLevel > 9 {
		return NewValidationError("wch_level", "Witch level must be between 0 and 9")
	}
	if i.ClrLevel < 0 || i.ClrLevel > 9 {
		return NewValidationError("clr_level", "Cleric level must be between 0 and 9")
	}
	if i.DrdLevel < 0 || i.DrdLevel > 9 {
		return NewValidationError("drd_level", "Druid level must be between 0 and 9")
	}

	if i.Range == "" {
		return NewValidationError("range", "Range cannot be empty")
	}
	if i.Duration == "" {
		return NewValidationError("duration", "Duration cannot be empty")
	}

	return nil
}
