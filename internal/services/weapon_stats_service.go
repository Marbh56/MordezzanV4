package services

import (
	"context"
	"fmt"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"strings"
)

type WeaponStatsService struct {
	inventoryRepo     repositories.InventoryRepository
	characterRepo     repositories.CharacterRepository
	weaponRepo        repositories.WeaponRepository
	weaponMasteryRepo repositories.WeaponMasteryRepository
}

type WeaponStats struct {
	Weapon             *models.Weapon         `json:"weapon"`
	InventoryItem      *models.InventoryItem  `json:"inventory_item"`
	BaseToHit          int                    `json:"base_to_hit"`
	ToHitBonus         int                    `json:"to_hit_bonus"`
	FinalToHit         int                    `json:"final_to_hit"`
	BaseDamage         string                 `json:"base_damage"`
	DamageBonus        int                    `json:"damage_bonus"`
	FinalDamage        string                 `json:"final_damage"`
	BaseAttackRate     string                 `json:"base_attack_rate"`
	ImprovedAttackRate bool                   `json:"improved_attack_rate"`
	FinalAttackRate    string                 `json:"final_attack_rate"`
	IsMastered         bool                   `json:"is_mastered"`
	MasteryLevel       string                 `json:"mastery_level,omitempty"`
	MasteryBonuses     map[string]interface{} `json:"mastery_bonuses,omitempty"`
}

func NewWeaponStatsService(
	inventoryRepo repositories.InventoryRepository,
	characterRepo repositories.CharacterRepository,
	weaponRepo repositories.WeaponRepository,
	weaponMasteryRepo repositories.WeaponMasteryRepository,
) *WeaponStatsService {
	return &WeaponStatsService{
		inventoryRepo:     inventoryRepo,
		characterRepo:     characterRepo,
		weaponRepo:        weaponRepo,
		weaponMasteryRepo: weaponMasteryRepo,
	}
}

func (s *WeaponStatsService) CalculateCharacterWeaponStats(ctx context.Context, characterID int64) ([]*WeaponStats, error) {
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	inventory, err := s.inventoryRepo.GetInventoryByCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Get all character's weapon masteries
	masteries, err := s.weaponMasteryRepo.GetWeaponMasteriesByCharacter(ctx, characterID)
	if err != nil {
		logger.Error("Failed to fetch weapon masteries: %v", err)
		// Continue without masteries if error
		masteries = []*models.WeaponMastery{}
	}

	// Create a map for quick mastery lookup
	masteryMap := make(map[string]*models.WeaponMastery)
	for _, mastery := range masteries {
		masteryMap[mastery.WeaponBaseName] = mastery
	}

	var weaponStats []*WeaponStats

	// Process all weapons in inventory
	for _, item := range inventory.Items {
		if item.ItemType != "weapon" {
			continue
		}

		weapon, err := s.weaponRepo.GetWeapon(ctx, item.ItemID)
		if err != nil {
			logger.Error("Failed to fetch weapon details for ID %d: %v", item.ItemID, err)
			continue
		}

		// Calculate base weapon name for mastery lookup
		baseWeaponName := extractBaseWeaponName(weapon.Name)

		// Determine if this is a missile weapon (ranged or hurled)
		isMissileWeapon := isRangedWeapon(weapon)

		// Initialize weapon stats with correct attack rate
		stats := &WeaponStats{
			Weapon:        weapon,
			InventoryItem: &item,
			BaseToHit:     0,
			BaseDamage:    weapon.Damage,
		}

		// Handle attack rate differently for melee vs missile weapons
		if isMissileWeapon {
			// For missile weapons, use RateOfFire from database
			stats.BaseAttackRate = weapon.RateOfFire
			if stats.BaseAttackRate == "" {
				// Default missile weapons to "1/1" if not specified
				stats.BaseAttackRate = "1/1"
			}
		} else {
			// All melee weapons have a standard attack rate of "1/1"
			stats.BaseAttackRate = "1/1"
		}

		// Check if weapon is mastered
		mastery, hasMastery := masteryMap[baseWeaponName]
		if hasMastery {
			stats.IsMastered = true
			stats.MasteryLevel = mastery.MasteryLevel
			stats.MasteryBonuses = models.GetWeaponMasteryBonuses(mastery.MasteryLevel)
		}

		// Calculate to-hit bonus
		if isMissileWeapon {
			// Ranged weapon uses DEX
			stats.ToHitBonus = character.RangedModifier
		} else {
			// Melee weapon uses STR
			stats.ToHitBonus = character.MeleeModifier
		}

		// Add weapon bonus if it exists (for magic weapons)
		weaponBonus := extractWeaponBonus(weapon.Name)
		stats.ToHitBonus += weaponBonus

		// Add mastery bonuses
		if stats.IsMastered {
			if bonus, ok := stats.MasteryBonuses["to_hit_bonus"].(int); ok {
				stats.ToHitBonus += bonus
			}
		}

		// Calculate final to-hit
		stats.FinalToHit = stats.BaseToHit + stats.ToHitBonus

		// Calculate damage bonus
		if !isMissileWeapon {
			// Only melee weapons get STR damage bonus
			stats.DamageBonus = character.DamageAdjustment
		}

		// Add weapon bonus to damage as well
		stats.DamageBonus += weaponBonus

		// Add mastery damage bonus
		if stats.IsMastered {
			if bonus, ok := stats.MasteryBonuses["damage_bonus"].(int); ok {
				stats.DamageBonus += bonus
			}
		}

		// Format final damage string
		stats.FinalDamage = formatDamageWithBonus(stats.BaseDamage, stats.DamageBonus)

		// Calculate attack rate based on mastery
		stats.FinalAttackRate = stats.BaseAttackRate
		if stats.IsMastered {
			if improvedRate, ok := stats.MasteryBonuses["improved_rate"].(bool); ok && improvedRate {
				stats.ImprovedAttackRate = true
				stats.FinalAttackRate = models.CalculateAttackRateWithMastery(stats.BaseAttackRate, stats.MasteryLevel)
			}
		}

		// Fighter-type classes get improved melee attack rate at level 7+
		if character.Class == "Fighter" || character.Class == "Ranger" ||
			character.Class == "Paladin" || character.Class == "Barbarian" ||
			character.Class == "Berserker" || character.Class == "Cataphract" ||
			character.Class == "Huntsman" {
			if character.Level >= 7 && !isMissileWeapon && stats.BaseAttackRate == "1/1" {
				if stats.FinalAttackRate == "1/1" { // Only upgrade if not already improved by mastery
					stats.FinalAttackRate = "3/2"
					stats.ImprovedAttackRate = true
				}
			}
		}

		weaponStats = append(weaponStats, stats)
	}

	return weaponStats, nil
}

// Helper functions
func isRangedWeapon(weapon *models.Weapon) bool {
	return weapon.Category == "Ranged" || weapon.Category == "Hurled"
}

func extractBaseWeaponName(name string) string {
	if idx := strings.Index(name, " +"); idx != -1 {
		name = name[:idx]
	}

	suffixes := []string{
		" of Slaying",
		" of Fire",
		" of Frost",
		" of Lightning",
		" of Venom",
		" of Speed",
		" of Accuracy",
		" of Power",
	}

	for _, suffix := range suffixes {
		if idx := strings.Index(name, suffix); idx != -1 {
			name = name[:idx]
			break
		}
	}

	return strings.TrimSpace(name)
}

func extractWeaponBonus(name string) int {
	if idx := strings.Index(name, " +"); idx != -1 {
		bonusPart := name[idx+2:]
		if spaceIdx := strings.Index(bonusPart, " "); spaceIdx != -1 {
			bonusPart = bonusPart[:spaceIdx]
		}

		var bonus int
		_, err := fmt.Sscanf(bonusPart, "%d", &bonus)
		if err == nil {
			return bonus
		}
	}
	return 0
}

func formatDamageWithBonus(baseDamage string, bonus int) string {
	if bonus == 0 {
		return baseDamage
	}

	if bonus > 0 {
		return fmt.Sprintf("%s+%d", baseDamage, bonus)
	}

	return fmt.Sprintf("%s%d", baseDamage, bonus)
}
