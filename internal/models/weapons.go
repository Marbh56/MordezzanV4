package models

import (
	"time"
)

type Weapon struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Category        string    `json:"category"`
	WeaponClass     int       `json:"weapon_class"`
	Cost            float64   `json:"cost"`
	Weight          int       `json:"weight"`
	RangeShort      *int      `json:"range_short,omitempty"`
	RangeMedium     *int      `json:"range_medium,omitempty"`
	RangeLong       *int      `json:"range_long,omitempty"`
	RateOfFire      string    `json:"rate_of_fire,omitempty"`
	Damage          string    `json:"damage"`
	DamageTwoHanded string    `json:"damage_two_handed,omitempty"`
	Properties      string    `json:"properties,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateWeaponInput struct {
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	WeaponClass     int     `json:"weapon_class"`
	Cost            float64 `json:"cost"`
	Weight          int     `json:"weight"`
	RangeShort      *int    `json:"range_short,omitempty"`
	RangeMedium     *int    `json:"range_medium,omitempty"`
	RangeLong       *int    `json:"range_long,omitempty"`
	RateOfFire      string  `json:"rate_of_fire,omitempty"`
	Damage          string  `json:"damage"`
	DamageTwoHanded string  `json:"damage_two_handed,omitempty"`
	Properties      string  `json:"properties,omitempty"`
}

type UpdateWeaponInput struct {
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	WeaponClass     int     `json:"weapon_class"`
	Cost            float64 `json:"cost"`
	Weight          int     `json:"weight"`
	RangeShort      *int    `json:"range_short,omitempty"`
	RangeMedium     *int    `json:"range_medium,omitempty"`
	RangeLong       *int    `json:"range_long,omitempty"`
	RateOfFire      string  `json:"rate_of_fire,omitempty"`
	Damage          string  `json:"damage"`
	DamageTwoHanded string  `json:"damage_two_handed,omitempty"`
	Properties      string  `json:"properties,omitempty"`
}

func (i *CreateWeaponInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Category == "" {
		return NewValidationError("category", "Category cannot be empty")
	}

	// For ranged weapons, we'll default to weapon class 1 if not specified
	if i.Category == "Ranged" || i.Category == "Hurled" {
		// No need to validate weapon class for ranged/hurled weapons
		if i.WeaponClass <= 0 {
			i.WeaponClass = 1 // Default weapon class for ranged/hurled weapons
		}
	} else {
		// For melee weapons, validate weapon class
		if i.WeaponClass <= 0 {
			return NewValidationError("weapon_class", "Weapon class must be positive for melee weapons")
		}
	}

	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	if i.Damage == "" {
		return NewValidationError("damage", "Damage cannot be empty")
	}
	return nil
}

func (i *UpdateWeaponInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Category == "" {
		return NewValidationError("category", "Category cannot be empty")
	}
	if i.WeaponClass < 0 {
		return NewValidationError("weight_class", "Weight class cannot be negative")
	}
	if i.Cost < 0 {
		return NewValidationError("cost", "Cost cannot be negative")
	}
	if i.Weight <= 0 {
		return NewValidationError("weight", "Weight must be positive")
	}
	if i.Damage == "" {
		return NewValidationError("damage", "Damage cannot be empty")
	}
	return nil
}
