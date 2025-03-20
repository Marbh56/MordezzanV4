package models

import (
	"fmt"
	"time"
)

// PreparedSpell represents a spell prepared by a character
type PreparedSpell struct {
	ID           int64     `json:"id"`
	CharacterID  int64     `json:"character_id"`
	SpellID      int64     `json:"spell_id"`
	SlotLevel    int       `json:"slot_level"`
	PreparedAt   time.Time `json:"prepared_at"`
	SpellDetails *Spell    `json:"spell_details,omitempty"`
}

// PrepareSpellInput is the input for preparing a spell
type PrepareSpellInput struct {
	CharacterID int64 `json:"character_id"`
	SpellID     int64 `json:"spell_id"`
	SlotLevel   int   `json:"slot_level"`
}

// Validate validates the PrepareSpellInput
func (i *PrepareSpellInput) Validate() error {
	if i.CharacterID <= 0 {
		return NewValidationError("character_id", "Character ID must be positive")
	}
	if i.SpellID <= 0 {
		return NewValidationError("spell_id", "Spell ID must be positive")
	}
	if i.SlotLevel <= 0 || i.SlotLevel > 9 {
		return NewValidationError("slot_level", "Slot level must be between 1 and 9")
	}
	return nil
}

// UnprepareSpellInput is the input for unpreparing a spell
type UnprepareSpellInput struct {
	CharacterID int64 `json:"character_id"`
	SpellID     int64 `json:"spell_id"`
}

// Validate validates the UnprepareSpellInput
func (i *UnprepareSpellInput) Validate() error {
	if i.CharacterID <= 0 {
		return NewValidationError("character_id", "Character ID must be positive")
	}
	if i.SpellID <= 0 {
		return NewValidationError("spell_id", "Spell ID must be positive")
	}
	return nil
}

// SpellSlots represents the number of spell slots by level
type SpellSlots struct {
	Level1 int `json:"level_1"`
	Level2 int `json:"level_2"`
	Level3 int `json:"level_3"`
	Level4 int `json:"level_4"`
	Level5 int `json:"level_5"`
	Level6 int `json:"level_6"`
	Level7 int `json:"level_7"`
	Level8 int `json:"level_8"`
	Level9 int `json:"level_9"`
}

// GetSlotsByLevel returns the number of slots for a specific level
func (s *SpellSlots) GetSlotsByLevel(level int) int {
	switch level {
	case 1:
		return s.Level1
	case 2:
		return s.Level2
	case 3:
		return s.Level3
	case 4:
		return s.Level4
	case 5:
		return s.Level5
	case 6:
		return s.Level6
	case 7:
		return s.Level7
	case 8:
		return s.Level8
	case 9:
		return s.Level9
	default:
		return 0
	}
}

// SetSlotsByLevel sets the number of slots for a specific level
func (s *SpellSlots) SetSlotsByLevel(level, slots int) error {
	if level < 1 || level > 9 {
		return fmt.Errorf("invalid spell slot level: %d", level)
	}

	switch level {
	case 1:
		s.Level1 = slots
	case 2:
		s.Level2 = slots
	case 3:
		s.Level3 = slots
	case 4:
		s.Level4 = slots
	case 5:
		s.Level5 = slots
	case 6:
		s.Level6 = slots
	case 7:
		s.Level7 = slots
	case 8:
		s.Level8 = slots
	case 9:
		s.Level9 = slots
	}
	return nil
}

// CharacterSpellPreparation contains information about a character's prepared spells
type CharacterSpellPreparation struct {
	Character          *Character       `json:"character"`
	AvailableSlots     SpellSlots       `json:"available_slots"`
	UsedSlots          SpellSlots       `json:"used_slots"`
	PreparedSpells     []*PreparedSpell `json:"prepared_spells"`
	SpellsInSpellbooks []*Spell         `json:"spells_in_spellbooks"`
}
