package models

import (
	"fmt"
	"time"
)

type Spellbook struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	TotalPages   int       `json:"total_pages"`
	UsedPages    int       `json:"used_pages"`
	SpellsStored []int64   `json:"spells_stored,omitempty"` // IDs of stored spells
	Value        int       `json:"value"`                   // Cost in gold pieces
	Weight       float64   `json:"weight"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateSpellbookInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TotalPages  int     `json:"total_pages"`
	Value       int     `json:"value"`
	Weight      float64 `json:"weight"`
}

type UpdateSpellbookInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TotalPages  int     `json:"total_pages"`
	UsedPages   int     `json:"used_pages"`
	Value       int     `json:"value"`
	Weight      float64 `json:"weight"`
}

// AddSpell adds a spell to the spellbook
func (s *Spellbook) AddSpell(spell *Spell, characterClass string) error {
	// Determine which level to use based on the character class
	var pagesRequired int

	switch characterClass {
	case "Magician":
		pagesRequired = spell.MagLevel
	case "Cryo-mancer":
		pagesRequired = spell.CryLevel
	case "Illusionist":
		pagesRequired = spell.IllLevel
	case "Necromancer":
		pagesRequired = spell.NecLevel
	case "Pyromancer":
		pagesRequired = spell.PyrLevel
	case "Witch":
		pagesRequired = spell.WchLevel
	case "Cleric":
		pagesRequired = spell.ClrLevel
	case "Druid":
		pagesRequired = spell.DrdLevel
	default:
		return fmt.Errorf("unsupported character class: %s", characterClass)
	}

	// Check if the spell is not available for this class (level 0)
	if pagesRequired == 0 {
		return fmt.Errorf("spell not available for %s class", characterClass)
	}

	// Check if there's enough space
	if s.UsedPages+pagesRequired > s.TotalPages {
		return fmt.Errorf("not enough pages in spellbook to add spell (requires %d pages, only %d available)",
			pagesRequired, s.TotalPages-s.UsedPages)
	}

	// Add the spell
	s.SpellsStored = append(s.SpellsStored, spell.ID)
	s.UsedPages += pagesRequired
	return nil
}

func (s *Spellbook) RemoveSpell(spell *Spell, characterClass string) error {
	// Find the spell in the stored spells
	index := -1
	for i, spellID := range s.SpellsStored {
		if spellID == spell.ID {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("spell not found in spellbook")
	}

	// Determine which level to use based on the character class
	var pagesUsed int

	switch characterClass {
	case "Magician":
		pagesUsed = spell.MagLevel
	case "Cryo-mancer":
		pagesUsed = spell.CryLevel
	case "Illusionist":
		pagesUsed = spell.IllLevel
	case "Necromancer":
		pagesUsed = spell.NecLevel
	case "Pyromancer":
		pagesUsed = spell.PyrLevel
	case "Witch":
		pagesUsed = spell.WchLevel
	case "Cleric":
		pagesUsed = spell.ClrLevel
	case "Druid":
		pagesUsed = spell.DrdLevel
	default:
		return fmt.Errorf("unsupported character class: %s", characterClass)
	}

	// Remove the spell
	s.SpellsStored = append(s.SpellsStored[:index], s.SpellsStored[index+1:]...)
	s.UsedPages -= pagesUsed
	return nil
}

// AvailablePages returns the number of free pages
func (s *Spellbook) AvailablePages() int {
	return s.TotalPages - s.UsedPages
}

func (i *CreateSpellbookInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.TotalPages <= 0 {
		return NewValidationError("total_pages", "Total pages must be positive")
	}
	if i.Value < 0 {
		return NewValidationError("value", "Value cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	return nil
}

func (i *UpdateSpellbookInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.TotalPages <= 0 {
		return NewValidationError("total_pages", "Total pages must be positive")
	}
	if i.UsedPages < 0 || i.UsedPages > i.TotalPages {
		return NewValidationError("used_pages", "Used pages must be between 0 and total pages")
	}
	if i.Value < 0 {
		return NewValidationError("value", "Value cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	return nil
}
