package services

import (
	"context"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
)

// ACService handles armor class calculations
type ACService struct {
	inventoryRepo      repositories.InventoryRepository
	characterRepo      repositories.CharacterRepository
	armorRepo          repositories.ArmorRepository
	shieldRepo         repositories.ShieldRepository
	encumbranceService *EncumbranceService
}

// ACDetails represents the components and final AC value
type ACDetails struct {
	BaseAC         int    `json:"base_ac"`
	ArmorAC        int    `json:"armor_ac,omitempty"`
	ShieldBonus    int    `json:"shield_bonus,omitempty"`
	DexterityMod   int    `json:"dexterity_mod,omitempty"`
	NaturalAC      int    `json:"natural_ac,omitempty"`
	AgileBonus     int    `json:"agile_bonus,omitempty"`
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
	encumbranceService *EncumbranceService,
) *ACService {
	return &ACService{
		inventoryRepo:      inventoryRepo,
		characterRepo:      characterRepo,
		armorRepo:          armorRepo,
		shieldRepo:         shieldRepo,
		encumbranceService: encumbranceService,
	}
}

func (s *ACService) CalculateCharacterAC(ctx context.Context, characterID int64) (*ACDetails, error) {
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	inventory, err := s.inventoryRepo.GetInventoryByCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	details := &ACDetails{
		BaseAC:       9,
		DexterityMod: character.DefenceAdjustment,
	}

	// Check for agile bonus (unarmored and unencumbered)
	if character.Class == "Thief" {
		wearingArmor := false
		for _, item := range inventory.Items {
			if item.ItemType == "armor" && item.IsEquipped {
				wearingArmor = true
				break
			}
		}

		// Check if heavily encumbered
		var isHeavyEncumbered bool
		if s.encumbranceService != nil {
			encumbranceDetails, err := s.encumbranceService.GetCharacterEncumbrance(ctx, characterID)
			if err == nil && encumbranceDetails != nil {
				isHeavyEncumbered = encumbranceDetails.Status.HeavyEncumbered
			}
		}

		// Apply the agile bonus if conditions are met
		if !wearingArmor && !isHeavyEncumbered {
			details.AgileBonus = 1
			logger.Info("Applied +1 agile AC bonus for unarmored Thief")
		}
	}

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

	if hasArmor {
		armor, err := s.armorRepo.GetArmor(ctx, equippedArmor.ItemID)
		if err == nil {
			details.ArmorAC = armor.AC
			details.ArmorEquipped = armor.Name
		} else {
			logger.Error("Failed to fetch armor details: %v", err)
		}
	}

	if hasShield {
		shield, err := s.shieldRepo.GetShield(ctx, equippedShield.ItemID)
		if err == nil {
			details.ShieldBonus = shield.DefenseModifier
			details.ShieldEquipped = shield.Name
		} else {
			logger.Error("Failed to fetch shield details: %v", err)
		}
	}

	if character.NaturalAC > 0 {
		details.NaturalAC = character.NaturalAC
	} else if character.Class == "Berserker" {
		if abilities, ok := character.Abilities.(map[string]interface{}); ok {
			if naturalAC, ok := abilities["natural_ac"].(int); ok {
				details.NaturalAC = naturalAC
			}
		}
	} else if character.Class == "Monk" {
		if abilities, ok := character.Abilities.(map[string]interface{}); ok {
			if acBonus, ok := abilities["ac_bonus"].(int); ok {
				details.OtherBonuses = acBonus
			}
		}
	}

	finalAC := details.BaseAC
	if details.ArmorAC > 0 {
		finalAC = details.ArmorAC
	}
	finalAC -= details.ShieldBonus
	finalAC -= details.DexterityMod
	finalAC -= details.AgileBonus
	finalAC -= details.NaturalAC
	finalAC -= details.OtherBonuses
	details.FinalAC = finalAC

	return details, nil
}
