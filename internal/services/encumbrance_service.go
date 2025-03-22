package services

import (
	"context"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"sort"
)

// EncumbranceService handles inventory weight and encumbrance calculations
type EncumbranceService struct {
	inventoryRepo   repositories.InventoryRepository
	characterRepo   repositories.CharacterRepository
	weaponRepo      repositories.WeaponRepository
	armorRepo       repositories.ArmorRepository
	shieldRepo      repositories.ShieldRepository
	potionRepo      repositories.PotionRepository
	magicItemRepo   repositories.MagicItemRepository
	ringRepo        repositories.RingRepository
	ammoRepo        repositories.AmmoRepository
	spellScrollRepo repositories.SpellScrollRepository
	containerRepo   repositories.ContainerRepository
	equipmentRepo   repositories.EquipmentRepository
	treasureRepo    repositories.TreasureRepository
}

// NewEncumbranceService creates a new encumbrance service
func NewEncumbranceService(
	inventoryRepo repositories.InventoryRepository,
	characterRepo repositories.CharacterRepository,
	weaponRepo repositories.WeaponRepository,
	armorRepo repositories.ArmorRepository,
	shieldRepo repositories.ShieldRepository,
	potionRepo repositories.PotionRepository,
	magicItemRepo repositories.MagicItemRepository,
	ringRepo repositories.RingRepository,
	ammoRepo repositories.AmmoRepository,
	spellScrollRepo repositories.SpellScrollRepository,
	containerRepo repositories.ContainerRepository,
	equipmentRepo repositories.EquipmentRepository,
	treasureRepo repositories.TreasureRepository,
) *EncumbranceService {
	return &EncumbranceService{
		inventoryRepo:   inventoryRepo,
		characterRepo:   characterRepo,
		weaponRepo:      weaponRepo,
		armorRepo:       armorRepo,
		shieldRepo:      shieldRepo,
		potionRepo:      potionRepo,
		magicItemRepo:   magicItemRepo,
		ringRepo:        ringRepo,
		ammoRepo:        ammoRepo,
		spellScrollRepo: spellScrollRepo,
		containerRepo:   containerRepo,
		equipmentRepo:   equipmentRepo,
		treasureRepo:    treasureRepo,
	}
}

// GetCharacterEncumbrance calculates the complete encumbrance details for a character
func (s *EncumbranceService) GetCharacterEncumbrance(ctx context.Context, characterID int64) (*models.InventoryWeightDetails, error) {
	// Get character info (for strength/constitution)
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Calculate encumbrance thresholds based on attributes
	thresholds := models.CalculateEncumbranceThresholds(character.Strength, character.Constitution)

	// Get inventory
	inventory, err := s.inventoryRepo.GetInventoryByCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Initialize weight details
	details := &models.InventoryWeightDetails{
		TotalWeight:  inventory.CurrentWeight,
		WeightByType: make(map[string]float64),
		Thresholds:   thresholds,
	}

	// Calculate encumbrance status
	details.Status = models.CalculateEncumbranceStatus(inventory.CurrentWeight, thresholds)

	// Get weight by type and collect item details for heaviest items
	weightedItems, err := s.calculateWeightByType(ctx, inventory)
	if err != nil {
		return nil, err
	}

	// Add any carried treasure weight
	if inventory.Treasure != nil {
		coinWeight := calculateTreasureWeight(inventory.Treasure)
		details.TotalWeight += coinWeight
		details.WeightByType["treasure"] = coinWeight
	}

	// Sort weighted items by total weight (descending)
	sort.Slice(weightedItems, func(i, j int) bool {
		return weightedItems[i].TotalWeight > weightedItems[j].TotalWeight
	})

	// Get the 5 heaviest items (or all if fewer than 5)
	heaviestCount := 5
	if len(weightedItems) < heaviestCount {
		heaviestCount = len(weightedItems)
	}
	details.HeaviestItems = weightedItems[:heaviestCount]

	// Recalculate status with most up-to-date weight
	details.Status = models.CalculateEncumbranceStatus(details.TotalWeight, thresholds)

	return details, nil
}

// calculateWeightByType calculates the weight breakdown by item type
func (s *EncumbranceService) calculateWeightByType(ctx context.Context, inventory *models.Inventory) ([]models.WeightedInventoryItem, error) {
	weightByType := make(map[string]float64)
	var weightedItems []models.WeightedInventoryItem

	// Process each inventory item
	for _, item := range inventory.Items {
		var itemWeight float64
		var itemName string

		// Get weight based on item type
		switch item.ItemType {
		case "weapon":
			if weapon, err := s.weaponRepo.GetWeapon(ctx, item.ItemID); err == nil {
				itemWeight = float64(weapon.Weight)
				itemName = weapon.Name
				weightByType["weapons"] += itemWeight * float64(item.Quantity)
			}
		case "armor":
			if armor, err := s.armorRepo.GetArmor(ctx, item.ItemID); err == nil {
				itemWeight = float64(armor.Weight)
				itemName = armor.Name
				weightByType["armor"] += itemWeight * float64(item.Quantity)
			}
		case "shield":
			if shield, err := s.shieldRepo.GetShield(ctx, item.ItemID); err == nil {
				itemWeight = float64(shield.Weight)
				itemName = shield.Name
				weightByType["shields"] += itemWeight * float64(item.Quantity)
			}
		case "potion":
			if potion, err := s.potionRepo.GetPotion(ctx, item.ItemID); err == nil {
				itemWeight = 0.5
				itemName = potion.Name
				weightByType["potions"] += itemWeight * float64(item.Quantity)
			}
		case "magic_item":
			if magicItem, err := s.magicItemRepo.GetMagicItem(ctx, item.ItemID); err == nil {
				itemWeight = float64(magicItem.Weight)
				itemName = magicItem.Name
				weightByType["magic_items"] += itemWeight * float64(item.Quantity)
			}
		case "ring":
			itemWeight = 0.1 // Standard ring weight (negligible)
			if ring, err := s.ringRepo.GetRing(ctx, item.ItemID); err == nil {
				itemName = ring.Name
			}
			weightByType["rings"] += itemWeight * float64(item.Quantity)
		case "ammo":
			if ammo, err := s.ammoRepo.GetAmmo(ctx, item.ItemID); err == nil {
				// Ammo weight is usually per bundle
				itemWeight = float64(ammo.Weight)
				itemName = ammo.Name
				weightByType["ammunition"] += itemWeight * float64(item.Quantity)
			}
		case "spell_scroll":
			itemWeight = 0.1 // Standard scroll weight
			if scroll, err := s.spellScrollRepo.GetSpellScroll(ctx, item.ItemID); err == nil {
				itemName = "Scroll of " + scroll.SpellName
			}
			weightByType["scrolls"] += itemWeight * float64(item.Quantity)
		case "container":
			if container, err := s.containerRepo.GetContainer(ctx, item.ItemID); err == nil {
				itemWeight = float64(container.Weight)
				itemName = container.Name
				weightByType["containers"] += itemWeight * float64(item.Quantity)
			}
		case "equipment":
			if equipment, err := s.equipmentRepo.GetEquipment(ctx, item.ItemID); err == nil {
				itemWeight = float64(equipment.Weight)
				itemName = equipment.Name
				weightByType["equipment"] += itemWeight * float64(item.Quantity)
			}
		}

		totalItemWeight := itemWeight * float64(item.Quantity)
		if itemWeight > 0 {
			weightedItems = append(weightedItems, models.WeightedInventoryItem{
				ID:          item.ID,
				Name:        itemName,
				ItemType:    item.ItemType,
				Weight:      itemWeight,
				TotalWeight: totalItemWeight,
				Quantity:    item.Quantity,
			})
		}
	}

	return weightedItems, nil
}

// calculateTreasureWeight calculates the weight of coins and gems
func calculateTreasureWeight(treasure *models.Treasure) float64 {
	// Standard weight assumptions:
	// - 50 coins = 1 lb
	// - Gems are negligible weight

	totalCoins := treasure.GoldCoins + treasure.SilverCoins + treasure.CopperCoins +
		treasure.ElectrumCoins + treasure.PlatinumCoins

	coinWeight := float64(totalCoins) / 50.0

	// Add weight of any gems/jewelry if tracking those by weight
	// gemWeight := float64(treasure.Gems) * gemWeightEach

	return coinWeight
}

// UpdateInventoryWeights recalculates and updates all weights for a character's inventory
func (s *EncumbranceService) UpdateInventoryWeights(ctx context.Context, characterID int64) error {
	// Get inventory
	inventory, err := s.inventoryRepo.GetInventoryByCharacter(ctx, characterID)
	if err != nil {
		return err
	}

	// Calculate new total weight
	weightedItems, err := s.calculateWeightByType(ctx, inventory)
	if err != nil {
		return err
	}

	var totalWeight float64
	for _, item := range weightedItems {
		totalWeight += item.TotalWeight
	}

	// Add treasure weight
	if inventory.Treasure != nil {
		totalWeight += calculateTreasureWeight(inventory.Treasure)
	}

	// Update inventory weight
	return s.inventoryRepo.UpdateInventoryWeight(ctx, inventory.ID, totalWeight)
}
