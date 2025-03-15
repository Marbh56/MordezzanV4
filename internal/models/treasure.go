package models

import (
	"time"
)

type Treasure struct {
	ID             int64     `json:"id"`
	CharacterID    *int64    `json:"character_id,omitempty"`
	PlatinumCoins  int       `json:"platinum_coins"`
	GoldCoins      int       `json:"gold_coins"`
	ElectrumCoins  int       `json:"electrum_coins"`
	SilverCoins    int       `json:"silver_coins"`
	CopperCoins    int       `json:"copper_coins"`
	Gems           string    `json:"gems,omitempty"`
	ArtObjects     string    `json:"art_objects,omitempty"`
	OtherValuables string    `json:"other_valuables,omitempty"`
	TotalValueGold float64   `json:"total_value_gold"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateTreasureInput struct {
	CharacterID    *int64  `json:"character_id,omitempty"`
	PlatinumCoins  int     `json:"platinum_coins"`
	GoldCoins      int     `json:"gold_coins"`
	ElectrumCoins  int     `json:"electrum_coins"`
	SilverCoins    int     `json:"silver_coins"`
	CopperCoins    int     `json:"copper_coins"`
	Gems           string  `json:"gems,omitempty"`
	ArtObjects     string  `json:"art_objects,omitempty"`
	OtherValuables string  `json:"other_valuables,omitempty"`
	TotalValueGold float64 `json:"total_value_gold"`
}

type UpdateTreasureInput struct {
	PlatinumCoins  int     `json:"platinum_coins"`
	GoldCoins      int     `json:"gold_coins"`
	ElectrumCoins  int     `json:"electrum_coins"`
	SilverCoins    int     `json:"silver_coins"`
	CopperCoins    int     `json:"copper_coins"`
	Gems           string  `json:"gems,omitempty"`
	ArtObjects     string  `json:"art_objects,omitempty"`
	OtherValuables string  `json:"other_valuables,omitempty"`
	TotalValueGold float64 `json:"total_value_gold"`
}

func (i *CreateTreasureInput) Validate() error {
	if i.PlatinumCoins < 0 {
		return NewValidationError("platinum_coins", "Platinum coins cannot be negative")
	}
	if i.GoldCoins < 0 {
		return NewValidationError("gold_coins", "Gold coins cannot be negative")
	}
	if i.ElectrumCoins < 0 {
		return NewValidationError("electrum_coins", "Electrum coins cannot be negative")
	}
	if i.SilverCoins < 0 {
		return NewValidationError("silver_coins", "Silver coins cannot be negative")
	}
	if i.CopperCoins < 0 {
		return NewValidationError("copper_coins", "Copper coins cannot be negative")
	}
	if i.TotalValueGold < 0 {
		return NewValidationError("total_value_gold", "Total value cannot be negative")
	}
	return nil
}

func (i *UpdateTreasureInput) Validate() error {
	if i.PlatinumCoins < 0 {
		return NewValidationError("platinum_coins", "Platinum coins cannot be negative")
	}
	if i.GoldCoins < 0 {
		return NewValidationError("gold_coins", "Gold coins cannot be negative")
	}
	if i.ElectrumCoins < 0 {
		return NewValidationError("electrum_coins", "Electrum coins cannot be negative")
	}
	if i.SilverCoins < 0 {
		return NewValidationError("silver_coins", "Silver coins cannot be negative")
	}
	if i.CopperCoins < 0 {
		return NewValidationError("copper_coins", "Copper coins cannot be negative")
	}
	if i.TotalValueGold < 0 {
		return NewValidationError("total_value_gold", "Total value cannot be negative")
	}
	return nil
}
