package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"mordezzanV4/internal/services"
)

type InventoryController struct {
	inventoryRepo      repositories.InventoryRepository
	characterRepo      repositories.CharacterRepository
	weaponRepo         repositories.WeaponRepository
	armorRepo          repositories.ArmorRepository
	shieldRepo         repositories.ShieldRepository
	potionRepo         repositories.PotionRepository
	magicItemRepo      repositories.MagicItemRepository
	ringRepo           repositories.RingRepository
	ammoRepo           repositories.AmmoRepository
	spellScrollRepo    repositories.SpellScrollRepository
	containerRepo      repositories.ContainerRepository
	equipmentRepo      repositories.EquipmentRepository
	treasureRepo       repositories.TreasureRepository
	encumbranceService *services.EncumbranceService
	tmpl               *template.Template
}

// EnrichedInventoryItem contains detailed item information
type EnrichedInventoryItem struct {
	ID          int64       `json:"id"`
	InventoryID int64       `json:"inventory_id"`
	ItemType    string      `json:"item_type"`
	ItemID      int64       `json:"item_id"`
	ItemDetails interface{} `json:"item_details"`
	Quantity    int         `json:"quantity"`
	IsEquipped  bool        `json:"is_equipped"`
	Slot        string      `json:"slot,omitempty"`
	Notes       string      `json:"notes,omitempty"`
}

type EquipmentStatus struct {
	EquippedSlots  map[string]ItemSummary `json:"equipped_slots"`
	AvailableSlots []string               `json:"available_slots"`
}

type ItemSummary struct {
	ID        int64  `json:"id"`
	ItemType  string `json:"item_type"`
	ItemID    int64  `json:"item_id"`
	Name      string `json:"name"`
	TwoHanded bool   `json:"two_handed,omitempty"`
}

func NewInventoryController(
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
	encumbranceService *services.EncumbranceService,
	tmpl *template.Template,
) *InventoryController {
	return &InventoryController{
		inventoryRepo:      inventoryRepo,
		characterRepo:      characterRepo,
		weaponRepo:         weaponRepo,
		armorRepo:          armorRepo,
		shieldRepo:         shieldRepo,
		potionRepo:         potionRepo,
		magicItemRepo:      magicItemRepo,
		ringRepo:           ringRepo,
		ammoRepo:           ammoRepo,
		spellScrollRepo:    spellScrollRepo,
		containerRepo:      containerRepo,
		equipmentRepo:      equipmentRepo,
		treasureRepo:       treasureRepo,
		encumbranceService: encumbranceService,
		tmpl:               tmpl,
	}
}

// Inventory handlers
func (c *InventoryController) GetInventory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid inventory ID format"))
		return
	}

	inventory, err := c.inventoryRepo.GetInventory(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(inventory); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	if err := c.tmpl.ExecuteTemplate(w, "inventory.html", inventory); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *InventoryController) GetInventoryByCharacter(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	// Verify character exists first
	_, err = c.characterRepo.GetCharacter(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Try to get inventory
	inventory, err := c.inventoryRepo.GetInventoryByCharacter(r.Context(), characterID)
	if err != nil {
		// Log the exact error to help diagnose the issue
		logger.Debug("GetInventoryByCharacter error: %v (type %T)", err, err)

		// Check for SQL no rows error specifically
		if errors.Is(err, sql.ErrNoRows) || apperrors.IsNotFound(err) {
			// Create a new inventory for this character
			logger.Info("Creating new inventory for character %d", characterID)
			input := &models.CreateInventoryInput{
				CharacterID: characterID,
				MaxWeight:   100.0, // Default capacity
			}
			inventoryID, err := c.inventoryRepo.CreateInventory(r.Context(), input)
			if err != nil {
				logger.Error("Failed to create inventory: %v", err)
				apperrors.HandleError(w, err)
				return
			}

			// Get the newly created inventory
			inventory, err = c.inventoryRepo.GetInventory(r.Context(), inventoryID)
			if err != nil {
				logger.Error("Failed to retrieve new inventory: %v", err)
				apperrors.HandleError(w, err)
				return
			}
		} else {
			// Handle other errors
			logger.Error("Unexpected error getting inventory: %v", err)
			apperrors.HandleError(w, err)
			return
		}
	}

	// Get inventory items and enrich them
	enrichedItems, err := c.enrichInventoryItems(r.Context(), inventory.Items)
	if err != nil {
		logger.Error("Failed to enrich inventory items: %v", err)
	}

	// Get treasure if available
	if c.treasureRepo != nil {
		treasure, err := c.treasureRepo.GetTreasureByCharacter(r.Context(), characterID)
		if err == nil && treasure != nil {
			inventory.Treasure = treasure
		}
	}

	// Get encumbrance details if the service is available
	var encumbranceDetails *models.InventoryWeightDetails
	if c.encumbranceService != nil {
		encumbranceDetails, err = c.encumbranceService.GetCharacterEncumbrance(r.Context(), characterID)
		if err != nil {
			logger.Error("Failed to get encumbrance details: %v", err)
		}
	}

	// Update inventory max weight if encumbrance service provided a capacity
	if encumbranceDetails != nil && encumbranceDetails.Thresholds.MaximumCapacity > 0 {
		maxCapacity := encumbranceDetails.Thresholds.MaximumCapacity
		if inventory.MaxWeight != maxCapacity {
			updateInput := &models.UpdateInventoryInput{
				MaxWeight: &maxCapacity,
			}
			if err := c.inventoryRepo.UpdateInventory(r.Context(), inventory.ID, updateInput); err != nil {
				logger.Error("Failed to update inventory max weight: %v", err)
			}
			inventory.MaxWeight = maxCapacity
		}
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"inventory":   inventory,
		"items":       enrichedItems,
		"encumbrance": encumbranceDetails,
	}); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *InventoryController) ListInventories(w http.ResponseWriter, r *http.Request) {
	inventories, err := c.inventoryRepo.ListInventories(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(inventories); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *InventoryController) CreateInventory(w http.ResponseWriter, r *http.Request) {
	var input models.CreateInventoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

	if err := input.Validate(); err != nil {
		var validationErr *models.ValidationError
		if errors.As(err, &validationErr) {
			validationErrors := map[string]string{
				validationErr.Field: validationErr.Message,
			}
			apperrors.HandleValidationErrors(w, validationErrors)
			return
		}
		apperrors.HandleError(w, err)
		return
	}

	// Check if character exists
	if c.characterRepo != nil {
		_, err := c.characterRepo.GetCharacter(r.Context(), input.CharacterID)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}
	}

	id, err := c.inventoryRepo.CreateInventory(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	inventory, err := c.inventoryRepo.GetInventory(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(inventory); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *InventoryController) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid inventory ID format"))
		return
	}

	var input models.UpdateInventoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

	if err := input.Validate(); err != nil {
		var validationErr *models.ValidationError
		if errors.As(err, &validationErr) {
			validationErrors := map[string]string{
				validationErr.Field: validationErr.Message,
			}
			apperrors.HandleValidationErrors(w, validationErrors)
			return
		}
		apperrors.HandleError(w, err)
		return
	}

	if err := c.inventoryRepo.UpdateInventory(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	updatedInventory, err := c.inventoryRepo.GetInventory(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedInventory); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *InventoryController) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid inventory ID format"))
		return
	}

	if err := c.inventoryRepo.DeleteInventory(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Inventory item handlers
func (c *InventoryController) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "itemId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid inventory item ID format"))
		return
	}

	item, err := c.inventoryRepo.GetInventoryItem(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Enrich item with details
	enrichedItem, err := c.enrichInventoryItem(r.Context(), *item)
	if err != nil {
		logger.Error("Failed to enrich inventory item: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(enrichedItem); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *InventoryController) AddInventoryItem(w http.ResponseWriter, r *http.Request) {
	inventoryID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid inventory ID format"))
		return
	}

	// Check if inventory exists
	_, err = c.inventoryRepo.GetInventory(r.Context(), inventoryID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	var input models.AddItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

	if err := input.Validate(); err != nil {
		var validationErr *models.ValidationError
		if errors.As(err, &validationErr) {
			validationErrors := map[string]string{
				validationErr.Field: validationErr.Message,
			}
			apperrors.HandleValidationErrors(w, validationErrors)
			return
		}
		apperrors.HandleError(w, err)
		return
	}

	// Verify that the item exists based on item type
	if err := c.validateItemExists(r.Context(), input.ItemType, input.ItemID); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// If we're adding it as equipped, validate slot assignment
	if input.IsEquipped {
		if err := c.validateEquipItem(r.Context(), inventoryID, input.ItemID, input.ItemType, input.Slot); err != nil {
			apperrors.HandleError(w, err)
			return
		}
	}

	// Add the item to inventory
	id, err := c.inventoryRepo.AddInventoryItem(r.Context(), inventoryID, &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Recalculate inventory weight
	if err := c.inventoryRepo.RecalculateInventoryWeight(r.Context(), inventoryID); err != nil {
		logger.Error("Failed to recalculate inventory weight: %v", err)
	}

	// Get updated inventory
	updatedInventory, err := c.inventoryRepo.GetInventory(r.Context(), inventoryID)
	if err != nil {
		logger.Error("Failed to get updated inventory: %v", err)
	}

	// Get the item
	item, err := c.inventoryRepo.GetInventoryItem(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Enrich item with details
	enrichedItem, err := c.enrichInventoryItem(r.Context(), *item)
	if err != nil {
		logger.Error("Failed to enrich inventory item: %v", err)
	}

	// Include weight information in the response
	response := struct {
		Item           EnrichedInventoryItem `json:"item"`
		TotalWeight    float64               `json:"total_weight"`
		WeightCapacity float64               `json:"weight_capacity"`
		IsOverweight   bool                  `json:"is_overweight"`
	}{
		Item:           enrichedItem,
		TotalWeight:    updatedInventory.CurrentWeight,
		WeightCapacity: updatedInventory.MaxWeight,
		IsOverweight:   updatedInventory.CurrentWeight > updatedInventory.MaxWeight,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *InventoryController) UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid inventory item ID format"))
		return
	}

	var input models.UpdateItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

	if err := input.Validate(); err != nil {
		var validationErr *models.ValidationError
		if errors.As(err, &validationErr) {
			validationErrors := map[string]string{
				validationErr.Field: validationErr.Message,
			}
			apperrors.HandleValidationErrors(w, validationErrors)
			return
		}
		apperrors.HandleError(w, err)
		return
	}

	// Get the existing item
	existingItem, err := c.inventoryRepo.GetInventoryItem(r.Context(), itemID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Handle equipment slot validation if we're equipping an item
	if input.IsEquipped != nil && *input.IsEquipped && !existingItem.IsEquipped {
		proposedSlot := ""
		if input.Slot != nil {
			proposedSlot = *input.Slot
		}

		// Validate that this item can be equipped in this slot
		if err := c.validateEquipItem(r.Context(), existingItem.InventoryID, existingItem.ItemID, existingItem.ItemType, proposedSlot); err != nil {
			apperrors.HandleError(w, err)
			return
		}

		// If no slot was provided but we found a valid one, use it
		if input.Slot == nil && proposedSlot != "" {
			input.Slot = &proposedSlot
		}
	} else if input.IsEquipped != nil && !*input.IsEquipped && existingItem.IsEquipped {
		// If unequipping, clear the slot
		emptySlot := ""
		input.Slot = &emptySlot
	}

	// Update the item
	if err := c.inventoryRepo.UpdateInventoryItem(r.Context(), itemID, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Recalculate inventory weight
	if err := c.inventoryRepo.RecalculateInventoryWeight(r.Context(), existingItem.InventoryID); err != nil {
		logger.Error("Failed to recalculate inventory weight: %v", err)
	}

	// Get updated inventory
	inventory, err := c.inventoryRepo.GetInventory(r.Context(), existingItem.InventoryID)
	if err != nil {
		logger.Error("Failed to get updated inventory: %v", err)
	}

	// Get the updated item
	updatedItem, err := c.inventoryRepo.GetInventoryItem(r.Context(), itemID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Enrich item with details
	enrichedItem, err := c.enrichInventoryItem(r.Context(), *updatedItem)
	if err != nil {
		logger.Error("Failed to enrich inventory item: %v", err)
	}

	// Include weight information in the response
	response := struct {
		Item           EnrichedInventoryItem `json:"item"`
		TotalWeight    float64               `json:"total_weight"`
		WeightCapacity float64               `json:"weight_capacity"`
		IsOverweight   bool                  `json:"is_overweight"`
	}{
		Item:           enrichedItem,
		TotalWeight:    inventory.CurrentWeight,
		WeightCapacity: inventory.MaxWeight,
		IsOverweight:   inventory.CurrentWeight > inventory.MaxWeight,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *InventoryController) RemoveInventoryItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid inventory item ID format"))
		return
	}

	// Get the existing item to find its inventory ID
	existingItem, err := c.inventoryRepo.GetInventoryItem(r.Context(), itemID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Store the inventory ID before deleting the item
	inventoryID := existingItem.InventoryID

	// Remove the item
	if err := c.inventoryRepo.RemoveInventoryItem(r.Context(), itemID); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Recalculate inventory weight
	if err := c.inventoryRepo.RecalculateInventoryWeight(r.Context(), inventoryID); err != nil {
		logger.Error("Failed to recalculate inventory weight: %v", err)
	}

	// Get updated inventory
	inventory, err := c.inventoryRepo.GetInventory(r.Context(), inventoryID)
	if err != nil {
		logger.Error("Failed to get updated inventory: %v", err)
	}

	// Return success response with updated weight info
	response := struct {
		Success        bool    `json:"success"`
		ItemID         int64   `json:"item_id"`
		TotalWeight    float64 `json:"total_weight"`
		WeightCapacity float64 `json:"weight_capacity"`
		IsOverweight   bool    `json:"is_overweight"`
	}{
		Success:        true,
		ItemID:         itemID,
		TotalWeight:    inventory.CurrentWeight,
		WeightCapacity: inventory.MaxWeight,
		IsOverweight:   inventory.CurrentWeight > inventory.MaxWeight,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// Helper methods
func (c *InventoryController) validateItemExists(ctx context.Context, itemType string, itemID int64) error {
	switch itemType {
	case "weapon":
		if c.weaponRepo != nil {
			if _, err := c.weaponRepo.GetWeapon(ctx, itemID); err != nil {
				return err
			}
		}
	case "armor":
		if c.armorRepo != nil {
			if _, err := c.armorRepo.GetArmor(ctx, itemID); err != nil {
				return err
			}
		}
	case "shield":
		if c.shieldRepo != nil {
			if _, err := c.shieldRepo.GetShield(ctx, itemID); err != nil {
				return err
			}
		}
	case "potion":
		if c.potionRepo != nil {
			if _, err := c.potionRepo.GetPotion(ctx, itemID); err != nil {
				return err
			}
		}
	case "magic_item":
		if c.magicItemRepo != nil {
			if _, err := c.magicItemRepo.GetMagicItem(ctx, itemID); err != nil {
				return err
			}
		}
	case "ring":
		if c.ringRepo != nil {
			if _, err := c.ringRepo.GetRing(ctx, itemID); err != nil {
				return err
			}
		}
	case "ammo":
		if c.ammoRepo != nil {
			if _, err := c.ammoRepo.GetAmmo(ctx, itemID); err != nil {
				return err
			}
		}
	case "spell_scroll":
		if c.spellScrollRepo != nil {
			if _, err := c.spellScrollRepo.GetSpellScroll(ctx, itemID); err != nil {
				return err
			}
		}
	case "container":
		if c.containerRepo != nil {
			if _, err := c.containerRepo.GetContainer(ctx, itemID); err != nil {
				return err
			}
		}
	case "equipment":
		if c.equipmentRepo != nil {
			if _, err := c.equipmentRepo.GetEquipment(ctx, itemID); err != nil {
				return err
			}
		}
	default:
		return apperrors.NewBadRequest("Invalid item type: " + itemType)
	}
	return nil
}

func (c *InventoryController) enrichInventoryItems(ctx context.Context, items []models.InventoryItem) ([]EnrichedInventoryItem, error) {
	enrichedItems := make([]EnrichedInventoryItem, 0, len(items))

	for _, item := range items {
		enrichedItem, err := c.enrichInventoryItem(ctx, item)
		if err != nil {
			logger.Error("Failed to enrich item %d of type %s: %v", item.ItemID, item.ItemType, err)
			// Add the item without details
			enrichedItems = append(enrichedItems, EnrichedInventoryItem{
				ID:          item.ID,
				InventoryID: item.InventoryID,
				ItemType:    item.ItemType,
				ItemID:      item.ItemID,
				ItemDetails: nil,
				Quantity:    item.Quantity,
				IsEquipped:  item.IsEquipped,
				Slot:        item.Slot, // Include the slot field
				Notes:       item.Notes,
			})
			continue
		}
		enrichedItems = append(enrichedItems, enrichedItem)
	}

	return enrichedItems, nil
}

func (c *InventoryController) enrichInventoryItem(ctx context.Context, item models.InventoryItem) (EnrichedInventoryItem, error) {
	var details interface{}
	var err error
	switch item.ItemType {
	case "weapon":
		if c.weaponRepo != nil {
			details, err = c.weaponRepo.GetWeapon(ctx, item.ItemID)
		}
	case "armor":
		if c.armorRepo != nil {
			details, err = c.armorRepo.GetArmor(ctx, item.ItemID)
		}
	case "shield":
		if c.shieldRepo != nil {
			details, err = c.shieldRepo.GetShield(ctx, item.ItemID)
		}
	case "potion":
		if c.potionRepo != nil {
			details, err = c.potionRepo.GetPotion(ctx, item.ItemID)
		}
	case "magic_item":
		if c.magicItemRepo != nil {
			details, err = c.magicItemRepo.GetMagicItem(ctx, item.ItemID)
		}
	case "ring":
		if c.ringRepo != nil {
			details, err = c.ringRepo.GetRing(ctx, item.ItemID)
		}
	case "ammo":
		if c.ammoRepo != nil {
			details, err = c.ammoRepo.GetAmmo(ctx, item.ItemID)
		}
	case "spell_scroll":
		if c.spellScrollRepo != nil {
			details, err = c.spellScrollRepo.GetSpellScroll(ctx, item.ItemID)
		}
	case "container":
		if c.containerRepo != nil {
			details, err = c.containerRepo.GetContainer(ctx, item.ItemID)
		}
	case "equipment":
		if c.equipmentRepo != nil {
			details, err = c.equipmentRepo.GetEquipment(ctx, item.ItemID)
		}
	default:
		logger.Warning("Unknown item type: %s for item ID %d", item.ItemType, item.ItemID)
		// Instead of returning an error, we'll return a basic info object
		details = map[string]interface{}{
			"name":        "Unknown Item",
			"description": "Item details could not be loaded",
			"weight":      0,
		}
	}

	if err != nil {
		logger.Error("Failed to fetch details for item type %s, ID %d: %v", item.ItemType, item.ItemID, err)
		// Instead of propagating the error, provide a fallback
		details = map[string]interface{}{
			"name":        fmt.Sprintf("%s (ID: %d)", item.ItemType, item.ItemID),
			"description": "Failed to load item details",
			"weight":      0,
		}
	}

	return EnrichedInventoryItem{
		ID:          item.ID,
		InventoryID: item.InventoryID,
		ItemType:    item.ItemType,
		ItemID:      item.ItemID,
		ItemDetails: details,
		Quantity:    item.Quantity,
		IsEquipped:  item.IsEquipped,
		Slot:        item.Slot,
		Notes:       item.Notes,
	}, nil
}

func (c *InventoryController) GetEncumbranceStatus(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	// Check if character exists
	if c.characterRepo != nil {
		_, err := c.characterRepo.GetCharacter(r.Context(), characterID)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}
	}

	// Get encumbrance details from service
	details, err := c.encumbranceService.GetCharacterEncumbrance(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(details); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// RecalculateEncumbrance recalculates a character's inventory weights and encumbrance status
func (c *InventoryController) RecalculateEncumbrance(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	// Check if character exists
	if c.characterRepo != nil {
		_, err := c.characterRepo.GetCharacter(r.Context(), characterID)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}
	}

	// Update inventory weights
	if err := c.encumbranceService.UpdateInventoryWeights(r.Context(), characterID); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Get updated encumbrance details
	details, err := c.encumbranceService.GetCharacterEncumbrance(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(details); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// UpdateInventoryCapacity updates a character's inventory capacity manually
func (c *InventoryController) UpdateInventoryCapacity(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	var input struct {
		MaxWeight float64 `json:"max_weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

	// Validate input
	if input.MaxWeight < 0 {
		apperrors.HandleError(w, apperrors.NewBadRequest("Maximum weight cannot be negative"))
		return
	}

	// Get inventory for this character
	inventory, err := c.inventoryRepo.GetInventoryByCharacter(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Update the max weight
	maxWeight := input.MaxWeight
	updateInput := &models.UpdateInventoryInput{
		MaxWeight: &maxWeight,
	}

	if err := c.inventoryRepo.UpdateInventory(r.Context(), inventory.ID, updateInput); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Get the updated encumbrance details
	details, err := c.encumbranceService.GetCharacterEncumbrance(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(details); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *InventoryController) validateEquipItem(ctx context.Context, inventoryID int64, itemID int64, itemType string, proposedSlot string) error {
	// Get all currently equipped items
	equippedItems, err := c.inventoryRepo.GetEquippedItems(ctx, inventoryID)
	if err != nil {
		return err
	}

	// Determine if this is a two-handed weapon
	isTwoHanded := false
	if itemType == "weapon" {
		weapon, err := c.weaponRepo.GetWeapon(ctx, itemID)
		if err != nil {
			return err
		}
		isTwoHanded = models.IsTwoHanded(weapon.Properties)
		if isTwoHanded {
			for _, item := range equippedItems {
				if item.Slot == string(models.SlotMainHand) || item.Slot == string(models.SlotOffHand) {
					return apperrors.NewBadRequest("Cannot equip a two-handed weapon while hands are occupied")
				}
			}
		}
	}

	if proposedSlot != "" {
		slot := models.EquipmentSlot(proposedSlot)
		// Verify the slot is valid for this item type
		validSlots := models.GetItemTypeSlots(itemType)
		validSlot := false
		for _, s := range validSlots {
			if s == slot {
				validSlot = true
				break
			}
		}
		if !validSlot {
			return apperrors.NewBadRequest("Invalid slot for this item type")
		}

		// Check if slot is already occupied
		for _, item := range equippedItems {
			if item.Slot == proposedSlot {
				return apperrors.NewBadRequest("Slot is already occupied")
			}
			if itemType == "shield" && item.ItemType == "weapon" &&
				models.IsTwoHanded(item.Notes) && item.Slot == string(models.SlotMainHand) {
				return apperrors.NewBadRequest("Cannot equip a shield while wielding a two-handed weapon")
			}
			if isTwoHanded && item.ItemType == "shield" &&
				item.Slot == string(models.SlotOffHand) {
				return apperrors.NewBadRequest("Cannot equip a two-handed weapon while using a shield")
			}
		}
	} else {
		// Auto-assign a slot if none provided
		var suggestedSlot string

		switch itemType {
		case "ring":
			hasLeftRing := false
			hasRightRing := false
			for _, item := range equippedItems {
				if item.ItemType == "ring" {
					if item.Slot == string(models.SlotRingLeft) {
						hasLeftRing = true
					} else if item.Slot == string(models.SlotRingRight) {
						hasRightRing = true
					}
				}
			}
			if !hasLeftRing {
				suggestedSlot = string(models.SlotRingLeft)
			} else if !hasRightRing {
				suggestedSlot = string(models.SlotRingRight)
			} else {
				return apperrors.NewBadRequest("Cannot equip more than two rings")
			}
		case "weapon":
			hasMainHand := false
			hasOffHand := false
			for _, item := range equippedItems {
				if item.Slot == string(models.SlotMainHand) {
					hasMainHand = true
				} else if item.Slot == string(models.SlotOffHand) {
					hasOffHand = true
				}
			}
			if isTwoHanded {
				if hasMainHand || hasOffHand {
					return apperrors.NewBadRequest("Two-handed weapons require both hands to be free")
				}
				suggestedSlot = string(models.SlotMainHand)
			} else {
				if !hasMainHand {
					suggestedSlot = string(models.SlotMainHand)
				} else if !hasOffHand {
					suggestedSlot = string(models.SlotOffHand)
				} else {
					return apperrors.NewBadRequest("Both hands are already occupied")
				}
			}
		case "shield":
			for _, item := range equippedItems {
				if item.Slot == string(models.SlotOffHand) {
					return apperrors.NewBadRequest("Off-hand is already occupied")
				}
				if item.ItemType == "weapon" && models.IsTwoHanded(item.Notes) {
					return apperrors.NewBadRequest("Cannot equip a shield with a two-handed weapon")
				}
			}
			suggestedSlot = string(models.SlotOffHand)
		case "armor":
			for _, item := range equippedItems {
				if item.ItemType == "armor" {
					return apperrors.NewBadRequest("Already wearing armor")
				}
			}
			suggestedSlot = string(models.SlotBody)
		default:
			validSlots := models.GetItemTypeSlots(itemType)
			if len(validSlots) == 0 {
				return apperrors.NewBadRequest("This item type cannot be equipped")
			}
			for _, slot := range validSlots {
				occupied := false
				for _, item := range equippedItems {
					if item.Slot == string(slot) {
						occupied = true
						break
					}
				}
				if !occupied {
					suggestedSlot = string(slot)
					break
				}
			}
			if suggestedSlot == "" {
				return apperrors.NewBadRequest("No available slots for this item")
			}
		}
		return c.validateEquipItem(ctx, inventoryID, itemID, itemType, suggestedSlot)
	}

	return nil
}

func (c *InventoryController) GetEquipmentStatus(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID: %v", err)
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	// Log the request for debugging
	logger.Debug("Getting equipment status for character ID: %d", characterID)

	// First, check if the character exists
	_, err = c.characterRepo.GetCharacter(r.Context(), characterID)
	if err != nil {
		logger.Error("Failed to get character: %v", err)
		apperrors.HandleError(w, err)
		return
	}

	// Next, get the inventory for this character
	inventory, err := c.inventoryRepo.GetInventoryByCharacter(r.Context(), characterID)
	if err != nil {
		// If inventory doesn't exist yet, create one
		if errors.Is(err, sql.ErrNoRows) || apperrors.IsNotFound(err) {
			logger.Info("Creating new inventory for character %d", characterID)
			input := &models.CreateInventoryInput{
				CharacterID: characterID,
				MaxWeight:   100.0,
			}
			inventoryID, err := c.inventoryRepo.CreateInventory(r.Context(), input)
			if err != nil {
				logger.Error("Failed to create inventory: %v", err)
				apperrors.HandleError(w, err)
				return
			}
			inventory, err = c.inventoryRepo.GetInventory(r.Context(), inventoryID)
			if err != nil {
				logger.Error("Failed to retrieve new inventory: %v", err)
				apperrors.HandleError(w, err)
				return
			}
		} else {
			logger.Error("Failed to get inventory: %v", err)
			apperrors.HandleError(w, err)
			return
		}
	}

	// Get the equipped items and slots
	status, err := c.getEquipmentStatus(r.Context(), inventory.ID)
	if err != nil {
		logger.Error("Failed to get equipment status: %v", err)
		apperrors.HandleError(w, err)
		return
	}

	// Return the status as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		logger.Error("Failed to encode equipment status: %v", err)
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *InventoryController) getEquipmentStatus(ctx context.Context, inventoryID int64) (*EquipmentStatus, error) {
	logger.Debug("Getting equipment status for inventory ID: %d", inventoryID)

	// Get all equipped items
	equippedItems, err := c.inventoryRepo.GetEquippedItems(ctx, inventoryID)
	if err != nil {
		logger.Error("Failed to get equipped items: %v", err)
		return nil, err
	}

	// Initialize the status with empty slots and all available slots
	status := &EquipmentStatus{
		EquippedSlots: make(map[string]ItemSummary),
		AvailableSlots: []string{
			"head", "body", "main_hand", "off_hand", "ring_left", "ring_right",
			"neck", "back", "belt", "feet", "hands",
		},
	}

	// Track occupied slots
	occupiedSlots := make(map[string]bool)

	// Process each equipped item
	for _, item := range equippedItems {
		// Skip items without slots
		if item.Slot == "" {
			logger.Debug("Skipping equipped item %d without slot", item.ID)
			continue
		}

		// Mark this slot as occupied
		occupiedSlots[item.Slot] = true

		// Get item details for the summary
		var name string
		var twoHanded bool

		switch item.ItemType {
		case "weapon":
			if c.weaponRepo != nil {
				weapon, err := c.weaponRepo.GetWeapon(ctx, item.ItemID)
				if err == nil {
					name = weapon.Name
					twoHanded = models.IsTwoHanded(weapon.Properties)
				} else {
					logger.Error("Failed to get weapon details: %v", err)
					name = "Unknown Weapon"
				}
			}
		case "armor":
			if c.armorRepo != nil {
				armor, err := c.armorRepo.GetArmor(ctx, item.ItemID)
				if err == nil {
					name = armor.Name
				} else {
					logger.Error("Failed to get armor details: %v", err)
					name = "Unknown Armor"
				}
			}
		case "shield":
			if c.shieldRepo != nil {
				shield, err := c.shieldRepo.GetShield(ctx, item.ItemID)
				if err == nil {
					name = shield.Name
				} else {
					logger.Error("Failed to get shield details: %v", err)
					name = "Unknown Shield"
				}
			}
		case "ring":
			if c.ringRepo != nil {
				ring, err := c.ringRepo.GetRing(ctx, item.ItemID)
				if err == nil {
					name = ring.Name
				} else {
					logger.Error("Failed to get ring details: %v", err)
					name = "Unknown Ring"
				}
			}
		default:
			name = fmt.Sprintf("Unknown %s", item.ItemType)
		}

		// Add to equipped slots
		status.EquippedSlots[item.Slot] = ItemSummary{
			ID:        item.ID,
			ItemType:  item.ItemType,
			ItemID:    item.ItemID,
			Name:      name,
			TwoHanded: twoHanded,
		}

		// If it's a two-handed weapon, mark off_hand as occupied too
		if twoHanded && item.Slot == string(models.SlotMainHand) {
			occupiedSlots[string(models.SlotOffHand)] = true
		}
	}

	// Update available slots by excluding occupied ones
	var availableSlots []string
	for _, slot := range status.AvailableSlots {
		if !occupiedSlots[slot] {
			availableSlots = append(availableSlots, slot)
		}
	}
	status.AvailableSlots = availableSlots

	logger.Debug("Equipment status: occupied slots=%v, available slots=%v",
		occupiedSlots, status.AvailableSlots)

	return status, nil
}

func (c *InventoryController) GetCombatEquipment(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", chi.URLParam(r, "id"))))
		return
	}

	// Get character's inventory
	inventory, err := c.inventoryRepo.GetInventoryByCharacter(r.Context(), characterID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("inventory", fmt.Sprintf("character %d", characterID)))
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Prepare response structure
	response := struct {
		Weapons []map[string]interface{} `json:"weapons"`
		Armor   []map[string]interface{} `json:"armor"`
	}{
		Weapons: []map[string]interface{}{},
		Armor:   []map[string]interface{}{},
	}

	// Process equipped items
	for _, item := range inventory.Items {
		if !item.IsEquipped {
			continue
		}

		switch item.ItemType {
		case "weapon":
			weapon, err := c.weaponRepo.GetWeapon(r.Context(), item.ItemID)
			if err == nil {
				weaponInfo := map[string]interface{}{
					"inventory_item": item,
					"weapon":         weapon,
				}
				response.Weapons = append(response.Weapons, weaponInfo)
			}

		case "armor":
			armor, err := c.armorRepo.GetArmor(r.Context(), item.ItemID)
			if err == nil {
				armorInfo := map[string]interface{}{
					"inventory_item": item,
					"armor":          armor,
				}
				response.Armor = append(response.Armor, armorInfo)
			}

		case "shield":
			shield, err := c.shieldRepo.GetShield(r.Context(), item.ItemID)
			if err == nil {
				shieldInfo := map[string]interface{}{
					"inventory_item": item,
					"shield":         shield,
				}
				response.Armor = append(response.Armor, shieldInfo)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}
}
