package models

import (
	"time"
)

type Armor struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	ArmorType       string    `json:"armor_type"`
	AC              int       `json:"ac"`
	Cost            float64   `json:"cost"`
	DamageReduction int       `json:"damage_reduction"`
	Weight          int       `json:"weight"`
	WeightClass     string    `json:"weight_class"`
	MovementRate    int       `json:"movement_rate"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateArmorInput struct {
	Name            string  `json:"name"`
	ArmorType       string  `json:"armor_type"`
	AC              int     `json:"ac"`
	Cost            float64 `json:"cost"`
	DamageReduction int     `json:"damage_reduction"`
	Weight          int     `json:"weight"`
	WeightClass     string  `json:"weight_class"`
	MovementRate    int     `json:"movement_rate"`
}

type UpdateArmorInput struct {
	Name            string  `json:"name"`
	ArmorType       string  `json:"armor_type"`
	AC              int     `json:"ac"`
	Cost            float64 `json:"cost"`
	DamageReduction int     `json:"damage_reduction"`
	Weight          int     `json:"weight"`
	WeightClass     string  `json:"weight_class"`
	MovementRate    int     `json:"movement_rate"`
}

func (i *CreateArmorInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.ArmorType == "" {
		return NewValidationError("armor_type", "Armor type cannot be empty")
	}
	if i.AC < 1 || i.AC > 9 {
		return NewValidationError("ac", "AC must be between 1 and 9")
	}
	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	if i.DamageReduction < 0 {
		return NewValidationError("damage_reduction", "Damage reduction cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	if i.WeightClass == "" {
		return NewValidationError("weight_class", "Weight class cannot be empty")
	}
	if !isValidWeightClass(i.WeightClass) {
		return NewValidationError("weight_class", "Weight class must be 'Light', 'Medium', or 'Heavy'")
	}
	if i.MovementRate <= 0 {
		return NewValidationError("movement_rate", "Movement rate must be positive")
	}
	return nil
}

func (i *UpdateArmorInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.ArmorType == "" {
		return NewValidationError("armor_type", "Armor type cannot be empty")
	}
	if i.AC < 1 || i.AC > 9 {
		return NewValidationError("ac", "AC must be between 1 and 9")
	}
	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	if i.DamageReduction < 0 {
		return NewValidationError("damage_reduction", "Damage reduction cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	if i.WeightClass == "" {
		return NewValidationError("weight_class", "Weight class cannot be empty")
	}
	if !isValidWeightClass(i.WeightClass) {
		return NewValidationError("weight_class", "Weight class must be 'Light', 'Medium', or 'Heavy'")
	}
	if i.MovementRate <= 0 {
		return NewValidationError("movement_rate", "Movement rate must be positive")
	}
	return nil
}

func isValidWeightClass(class string) bool {
	return class == "Light" || class == "Medium" || class == "Heavy"
}
