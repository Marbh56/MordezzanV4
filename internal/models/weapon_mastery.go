package models

import (
	"time"
)

// WeaponMastery represents a character's mastery of a specific weapon type
type WeaponMastery struct {
	ID             int64     `json:"id"`
	CharacterID    int64     `json:"character_id"`
	WeaponBaseName string    `json:"weapon_base_name"`
	MasteryLevel   string    `json:"mastery_level"` // "mastered" or "grand_mastery"
	WeaponName     string    `json:"weapon_name,omitempty"`
	WeaponCategory string    `json:"weapon_category,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AddWeaponMasteryInput is used for creating a new weapon mastery
type AddWeaponMasteryInput struct {
	CharacterID    int64  `json:"character_id"`
	WeaponBaseName string `json:"weapon_base_name"`
	MasteryLevel   string `json:"mastery_level"`
}

// UpdateWeaponMasteryInput is used for updating an existing weapon mastery
type UpdateWeaponMasteryInput struct {
	MasteryLevel string `json:"mastery_level"`
}

// Validate ensures that the input for adding a weapon mastery is valid
func (i *AddWeaponMasteryInput) Validate() error {
	if i.CharacterID <= 0 {
		return NewValidationError("character_id", "Character ID must be positive")
	}
	if i.WeaponBaseName == "" {
		return NewValidationError("weapon_base_name", "Weapon base name cannot be empty")
	}
	if i.MasteryLevel != "mastered" && i.MasteryLevel != "grand_mastery" {
		return NewValidationError("mastery_level", "Mastery level must be either 'mastered' or 'grand_mastery'")
	}
	return nil
}

// Validate ensures that the input for updating a weapon mastery is valid
func (i *UpdateWeaponMasteryInput) Validate() error {
	if i.MasteryLevel != "mastered" && i.MasteryLevel != "grand_mastery" {
		return NewValidationError("mastery_level", "Mastery level must be either 'mastered' or 'grand_mastery'")
	}
	return nil
}

// WeaponMasteryRule defines the game rules for weapon mastery by class and level
type WeaponMasteryRule struct {
	ClassName  string `json:"class_name"`
	LevelRange [2]int `json:"level_range"` // [min_level, max_level]
	MaxSlots   int    `json:"max_slots"`
}

// GetWeaponMasteryRules returns the rules for how many weapon masteries a character can have
func GetWeaponMasteryRules() []WeaponMasteryRule {
	return []WeaponMasteryRule{
		{ClassName: "Fighter", LevelRange: [2]int{1, 3}, MaxSlots: 2},
		{ClassName: "Fighter", LevelRange: [2]int{4, 7}, MaxSlots: 3},
		{ClassName: "Fighter", LevelRange: [2]int{8, 11}, MaxSlots: 4},
		{ClassName: "Fighter", LevelRange: [2]int{12, 100}, MaxSlots: 5}, // 100 as upper bound for simplicity

		{ClassName: "Ranger", LevelRange: [2]int{1, 3}, MaxSlots: 2},
		{ClassName: "Ranger", LevelRange: [2]int{4, 7}, MaxSlots: 3},
		{ClassName: "Ranger", LevelRange: [2]int{8, 11}, MaxSlots: 4},
		{ClassName: "Ranger", LevelRange: [2]int{12, 100}, MaxSlots: 5},

		{ClassName: "Paladin", LevelRange: [2]int{1, 3}, MaxSlots: 2},
		{ClassName: "Paladin", LevelRange: [2]int{4, 7}, MaxSlots: 3},
		{ClassName: "Paladin", LevelRange: [2]int{8, 11}, MaxSlots: 4},
		{ClassName: "Paladin", LevelRange: [2]int{12, 100}, MaxSlots: 5},

		{ClassName: "Barbarian", LevelRange: [2]int{1, 3}, MaxSlots: 2},
		{ClassName: "Barbarian", LevelRange: [2]int{4, 7}, MaxSlots: 3},
		{ClassName: "Barbarian", LevelRange: [2]int{8, 11}, MaxSlots: 4},
		{ClassName: "Barbarian", LevelRange: [2]int{12, 100}, MaxSlots: 5},

		{ClassName: "Cataphract", LevelRange: [2]int{1, 3}, MaxSlots: 2},
		{ClassName: "Cataphract", LevelRange: [2]int{4, 7}, MaxSlots: 3},
		{ClassName: "Cataphract", LevelRange: [2]int{8, 11}, MaxSlots: 4},
		{ClassName: "Cataphract", LevelRange: [2]int{12, 100}, MaxSlots: 5},

		{ClassName: "Huntsman", LevelRange: [2]int{1, 3}, MaxSlots: 2},
		{ClassName: "Huntsman", LevelRange: [2]int{4, 7}, MaxSlots: 3},
		{ClassName: "Huntsman", LevelRange: [2]int{8, 11}, MaxSlots: 4},
		{ClassName: "Huntsman", LevelRange: [2]int{12, 100}, MaxSlots: 5},

		{ClassName: "Berserker", LevelRange: [2]int{1, 3}, MaxSlots: 2},
		{ClassName: "Berserker", LevelRange: [2]int{4, 7}, MaxSlots: 3},
		{ClassName: "Berserker", LevelRange: [2]int{8, 11}, MaxSlots: 4},
		{ClassName: "Berserker", LevelRange: [2]int{12, 100}, MaxSlots: 5},
	}
}

// GetWeaponMasteryBonuses returns the bonuses provided by weapon mastery
func GetWeaponMasteryBonuses(masteryLevel string) map[string]interface{} {
	if masteryLevel == "grand_mastery" {
		return map[string]interface{}{
			"to_hit_bonus":      2,
			"damage_bonus":      2,
			"improved_rate":     true,
			"critical_improved": true,
		}
	}

	// Regular mastery
	return map[string]interface{}{
		"to_hit_bonus":      1,
		"damage_bonus":      1,
		"improved_rate":     false,
		"critical_improved": false,
	}
}

// CanHaveWeaponMastery checks if a character class can have weapon mastery
func CanHaveWeaponMastery(className string) bool {
	weaponMasteryClasses := []string{
		"Fighter", "Ranger", "Paladin", "Barbarian",
		"Cataphract", "Huntsman", "Berserker",
	}

	for _, c := range weaponMasteryClasses {
		if c == className {
			return true
		}
	}

	return false
}

// GetAvailableMasterySlots returns the number of mastery slots available for a class and level
func GetAvailableMasterySlots(className string, level int) int {
	if !CanHaveWeaponMastery(className) {
		return 0
	}

	rules := GetWeaponMasteryRules()
	for _, rule := range rules {
		if rule.ClassName == className &&
			level >= rule.LevelRange[0] && level <= rule.LevelRange[1] {
			return rule.MaxSlots
		}
	}

	// Default fallback
	return 0
}

// CanHaveGrandMastery checks if a character can have grand mastery based on their level
func CanHaveGrandMastery(level int) bool {
	return level >= 4
}

// GetMasteryProgress returns a struct detailing mastery progression for a class
type MasteryProgress struct {
	BaseMasteries         int  `json:"base_masteries"`
	Level4Unlock          bool `json:"level_4_unlock"`
	Level8Unlock          bool `json:"level_8_unlock"`
	Level12Unlock         bool `json:"level_12_unlock"`
	GrandMasteryAvailable bool `json:"grand_mastery_available"`
}

// GetMasteryProgress returns the mastery progression information for a character
func GetMasteryProgressInfo(className string, level int) MasteryProgress {
	if !CanHaveWeaponMastery(className) {
		return MasteryProgress{}
	}

	return MasteryProgress{
		BaseMasteries:         2,
		Level4Unlock:          level >= 4,
		Level8Unlock:          level >= 8,
		Level12Unlock:         level >= 12,
		GrandMasteryAvailable: level >= 4,
	}
}

// WeaponMasteryEffect represents the calculated combat effects of weapon mastery
type WeaponMasteryEffect struct {
	BaseWeaponName string `json:"base_weapon_name"`
	Level          string `json:"level"`
	ToHitBonus     int    `json:"to_hit_bonus"`
	DamageBonus    int    `json:"damage_bonus"`
	ImprovedRate   bool   `json:"improved_rate"`

	// Combat values with mastery applied
	ModifiedAttackRate string `json:"modified_attack_rate,omitempty"`

	// For display in UI
	Description string `json:"description"`
}

func CalculateAttackRateWithMastery(baseRate string, masteryLevel string) string {
	// These rates are for standard melee weapons
	switch baseRate {
	case "1/2":
		return "1/1" // 1/2 → 1/1
	case "1/1":
		return "3/2" // 1/1 → 3/2
	case "3/2":
		return "2/1" // 3/2 → 2/1
	case "2/1":
		return "5/2" // 2/1 → 5/2
	case "5/2":
		return "3/1" // 5/2 → 3/1
	}

	// If we didn't match any expected pattern, return the base rate
	return baseRate
}
