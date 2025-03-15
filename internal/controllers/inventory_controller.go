package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
)

type InventoryController struct {
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
	tmpl            *template.Template
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
	Notes       string      `json:"notes,omitempty"`
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
	tmpl *template.Template,
) *InventoryController {
	return &InventoryController{
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
		tmpl:            tmpl,
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

	// Check if character exists
	if c.characterRepo != nil {
		_, err := c.characterRepo.GetCharacter(r.Context(), characterID)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}
	}

	inventory, err := c.inventoryRepo.GetInventoryByCharacter(r.Context(), characterID)
	if err != nil {
		if apperrors.IsNotFound(err) {
			// If inventory doesn't exist, create a default one
			input := &models.CreateInventoryInput{
				CharacterID: characterID,
				MaxWeight:   100.0, // Default max weight
			}
			inventoryID, err := c.inventoryRepo.CreateInventory(r.Context(), input)
			if err != nil {
				apperrors.HandleError(w, err)
				return
			}
			inventory, err = c.inventoryRepo.GetInventory(r.Context(), inventoryID)
			if err != nil {
				apperrors.HandleError(w, err)
				return
			}
		} else {
			apperrors.HandleError(w, err)
			return
		}
	}

	// Enrich inventory items with details
	enrichedItems, err := c.enrichInventoryItems(r.Context(), inventory.Items)
	if err != nil {
		logger.Error("Failed to enrich inventory items: %v", err)
	}

	// Include any treasure for this character
	if c.treasureRepo != nil {
		treasure, err := c.treasureRepo.GetTreasureByCharacter(r.Context(), characterID)
		if err == nil && treasure != nil {
			// Only add if we successfully found treasure
			inventory.Treasure = treasure
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"inventory": inventory,
		"items":     enrichedItems,
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

	id, err := c.inventoryRepo.AddInventoryItem(r.Context(), inventoryID, &input)
	if err != nil {
		apperrors.HandleError(w, err)
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
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(enrichedItem); err != nil {
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

	if err := c.inventoryRepo.UpdateInventoryItem(r.Context(), itemID, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}

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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(enrichedItem); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *InventoryController) RemoveInventoryItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid inventory item ID format"))
		return
	}

	if err := c.inventoryRepo.RemoveInventoryItem(r.Context(), itemID); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		return EnrichedInventoryItem{}, apperrors.NewBadRequest("Invalid item type: " + item.ItemType)
	}

	if err != nil {
		return EnrichedInventoryItem{}, err
	}

	return EnrichedInventoryItem{
		ID:          item.ID,
		InventoryID: item.InventoryID,
		ItemType:    item.ItemType,
		ItemID:      item.ItemID,
		ItemDetails: details,
		Quantity:    item.Quantity,
		IsEquipped:  item.IsEquipped,
		Notes:       item.Notes,
	}, nil
}
