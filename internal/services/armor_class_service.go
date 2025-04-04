package services

import (
	"context"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
)

// ACService handles armor class calculations
type ACService struct {
	inventoryRepo repositories.InventoryRepository
	characterRepo repositories.CharacterRepository
	armorRepo     repositories.ArmorRepository
	shieldRepo    repositories.ShieldRepository
}

// ACDetails represents the components and final AC value
type ACDetails struct {
	BaseAC         int    `json:"base_ac"`
	ArmorAC        int    `json:"armor_ac,omitempty"`
	ShieldBonus    int    `json:"shield_bonus,omitempty"`
	DexterityMod   int    `json:"dexterity_mod,omitempty"`
	NaturalAC      int    `json:"natural_ac,omitempty"`
	OtherBonuses   int    `json:"other_bonuses,omitempty"`
	FinalAC        int    `json:"final_ac"`
	ArmorEquipped  string `json:"armor_equipped,omitempty"`
	ShieldEquipped string `json:"shield_equipped,omitempty"`
}

// NewACService creates a new armor class service
func NewACService(
	inventoryRepo repositories.InventoryRepository,
	characterRepo repositories.CharacterRepository,
	armorRepo repositories.ArmorRepository,
	shieldRepo repositories.ShieldRepository,
) *ACService {
	return &ACService{
		inventoryRepo: inventoryRepo,
		characterRepo: characterRepo,
		armorRepo:     armorRepo,
		shieldRepo:    shieldRepo,
	}
}

// CalculateCharacterAC computes a character's armor class details
func (s *ACService) CalculateCharacterAC(ctx context.Context, characterID int64) (*ACDetails, error) {
	// Get the character to check if it exists and get dexterity modifier
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Get the character's inventory
	inventory, err := s.inventoryRepo.GetInventoryByCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Set up the response with default values
	details := &ACDetails{
		BaseAC:       9,
		DexterityMod: character.DefenceAdjustment,
	}

	// Get equipped armor and shield (if any)
	var equippedArmor models.InventoryItem
	var equippedShield models.InventoryItem
	var hasArmor bool
	var hasShield bool

	for _, item := range inventory.Items {
		if item.IsEquipped {
			if item.ItemType == "armor" {
				equippedArmor = item
				hasArmor = true
			} else if item.ItemType == "shield" {
				equippedShield = item
				hasShield = true
			}
		}
	}

	// Apply armor AC if equipped
	if hasArmor {
		armor, err := s.armorRepo.GetArmor(ctx, equippedArmor.ItemID)
		if err == nil {
			// In Hyperborea, equipping armor sets the base AC
			details.ArmorAC = armor.AC
			details.ArmorEquipped = armor.Name
		} else {
			logger.Error("Failed to fetch armor details: %v", err)
		}
	}

	// Apply shield bonus if equipped
	if hasShield {
		shield, err := s.shieldRepo.GetShield(ctx, equippedShield.ItemID)
		if err == nil {
			details.ShieldBonus = shield.DefenseModifier
			details.ShieldEquipped = shield.Name
		} else {
			logger.Error("Failed to fetch shield details: %v", err)
		}
	}

	// Check for natural armor bonuses (some classes might have this)
	if character.NaturalAC > 0 {
		details.NaturalAC = character.NaturalAC
	} else if character.Class == "Berserker" {
		// Special case for Berserker class which has natural armor
		if abilities, ok := character.Abilities.(map[string]interface{}); ok {
			if naturalAC, ok := abilities["natural_ac"].(int); ok {
				details.NaturalAC = naturalAC
			}
		}
	} else if character.Class == "Monk" {
		// Monks also have AC bonuses
		if abilities, ok := character.Abilities.(map[string]interface{}); ok {
			if acBonus, ok := abilities["ac_bonus"].(int); ok {
				details.OtherBonuses = acBonus
			}
		}
	}

	// Calculate the final AC
	finalAC := details.BaseAC

	// If armor is equipped, it replaces the base AC
	if details.ArmorAC > 0 {
		finalAC = details.ArmorAC
	}

	// Apply shield, dexterity and other bonuses (they reduce AC which is better)
	finalAC -= details.ShieldBonus
	finalAC -= details.DexterityMod
	finalAC -= details.NaturalAC
	finalAC -= details.OtherBonuses

	details.FinalAC = finalAC

	return details, nil
}
