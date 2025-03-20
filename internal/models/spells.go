package models

import (
	"fmt"
	"time"
)

// Spell represents a spell in the game
type Spell struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	MagLevel     int       `json:"mag_level"`
	CryLevel     int       `json:"cry_level"`
	IllLevel     int       `json:"ill_level"`
	NecLevel     int       `json:"nec_level"`
	PyrLevel     int       `json:"pyr_level"`
	WchLevel     int       `json:"wch_level"`
	ClrLevel     int       `json:"clr_level"`
	DrdLevel     int       `json:"drd_level"`
	Range        string    `json:"range"`
	Duration     string    `json:"duration"`
	AreaOfEffect string    `json:"area_of_effect,omitempty"`
	Components   string    `json:"components,omitempty"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateSpellInput is used when creating a new spell
type CreateSpellInput struct {
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
	AreaOfEffect string `json:"area_of_effect,omitempty"`
	Components   string `json:"components,omitempty"`
	Description  string `json:"description"`
}

// Validate checks if the input is valid
func (i *CreateSpellInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name is required")
	}
	if i.Range == "" {
		return NewValidationError("range", "Range is required")
	}
	if i.Duration == "" {
		return NewValidationError("duration", "Duration is required")
	}
	if i.Description == "" {
		return NewValidationError("description", "Description is required")
	}

	// Check that at least one class level is set
	if i.MagLevel == 0 && i.CryLevel == 0 && i.IllLevel == 0 &&
		i.NecLevel == 0 && i.PyrLevel == 0 && i.WchLevel == 0 &&
		i.ClrLevel == 0 && i.DrdLevel == 0 {
		return NewValidationError("level", "At least one class level must be specified")
	}

	return nil
}

// UpdateSpellInput is used when updating an existing spell
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
	AreaOfEffect string `json:"area_of_effect,omitempty"`
	Components   string `json:"components,omitempty"`
	Description  string `json:"description"`
}

// Validate checks if the input is valid
func (i *UpdateSpellInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name is required")
	}
	if i.Range == "" {
		return NewValidationError("range", "Range is required")
	}
	if i.Duration == "" {
		return NewValidationError("duration", "Duration is required")
	}
	if i.Description == "" {
		return NewValidationError("description", "Description is required")
	}

	// Check that at least one class level is set
	if i.MagLevel == 0 && i.CryLevel == 0 && i.IllLevel == 0 &&
		i.NecLevel == 0 && i.PyrLevel == 0 && i.WchLevel == 0 &&
		i.ClrLevel == 0 && i.DrdLevel == 0 {
		return NewValidationError("level", "At least one class level must be specified")
	}

	return nil
}

// GetLevel returns the spell level for a specific class
func (s *Spell) GetLevel(class string) int {
	switch class {
	case "Magician":
		return s.MagLevel
	case "Cryomancer":
		return s.CryLevel
	case "Illusionist":
		return s.IllLevel
	case "Necromancer":
		return s.NecLevel
	case "Pyromancer":
		return s.PyrLevel
	case "Witch":
		return s.WchLevel
	case "Cleric":
		return s.ClrLevel
	case "Druid":
		return s.DrdLevel
	default:
		return 0
	}
}

// GetClassLevels returns a formatted string of all classes and levels
func (s *Spell) GetClassLevels() string {
	levels := ""

	if s.MagLevel > 0 {
		levels += fmt.Sprintf("Magician: %d, ", s.MagLevel)
	}
	if s.CryLevel > 0 {
		levels += fmt.Sprintf("Cryomancer: %d, ", s.CryLevel)
	}
	if s.IllLevel > 0 {
		levels += fmt.Sprintf("Illusionist: %d, ", s.IllLevel)
	}
	if s.NecLevel > 0 {
		levels += fmt.Sprintf("Necromancer: %d, ", s.NecLevel)
	}
	if s.PyrLevel > 0 {
		levels += fmt.Sprintf("Pyromancer: %d, ", s.PyrLevel)
	}
	if s.WchLevel > 0 {
		levels += fmt.Sprintf("Witch: %d, ", s.WchLevel)
	}
	if s.ClrLevel > 0 {
		levels += fmt.Sprintf("Cleric: %d, ", s.ClrLevel)
	}
	if s.DrdLevel > 0 {
		levels += fmt.Sprintf("Druid: %d, ", s.DrdLevel)
	}

	// Remove trailing comma and space
	if len(levels) > 2 {
		levels = levels[:len(levels)-2]
	}

	return levels
}
