package services

import (
	"context"
	"fmt"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
)

// SpellPreparationService handles business logic for preparing spells
type SpellPreparationService struct {
	preparedSpellRepo repositories.PreparedSpellRepository
	spellRepo         repositories.SpellRepository
	characterRepo     repositories.CharacterRepository
	inventoryRepo     repositories.InventoryRepository
	spellbookRepo     repositories.SpellbookRepository
	magicianService   *MagicianService
}

// NewSpellPreparationService creates a new spell preparation service
func NewSpellPreparationService(
	preparedSpellRepo repositories.PreparedSpellRepository,
	spellRepo repositories.SpellRepository,
	characterRepo repositories.CharacterRepository,
	inventoryRepo repositories.InventoryRepository,
	spellbookRepo repositories.SpellbookRepository,
	magicianService *MagicianService,
) *SpellPreparationService {
	return &SpellPreparationService{
		preparedSpellRepo: preparedSpellRepo,
		spellRepo:         spellRepo,
		characterRepo:     characterRepo,
		inventoryRepo:     inventoryRepo,
		spellbookRepo:     spellbookRepo,
		magicianService:   magicianService,
	}
}

// GetPreparedSpells gets all prepared spells for a character with spell details
func (s *SpellPreparationService) GetPreparedSpells(ctx context.Context, characterID int64) ([]*models.PreparedSpell, error) {
	// Verify character exists
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Check if character is a spellcaster
	if character.Class != "Magician" && character.Class != "Cleric" && character.Class != "Druid" {
		return nil, apperrors.NewBadRequest(fmt.Sprintf("Character class '%s' cannot prepare spells", character.Class))
	}

	// Get prepared spells
	preparedSpells, err := s.preparedSpellRepo.GetPreparedSpellsByCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Enrich with spell details
	for i, ps := range preparedSpells {
		spell, err := s.spellRepo.GetSpell(ctx, ps.SpellID)
		if err == nil {
			preparedSpells[i].SpellDetails = spell
		}
	}

	return preparedSpells, nil
}

func (s *SpellPreparationService) GetAvailableSpellSlots(ctx context.Context, characterID int64) (*models.SpellSlots, error) {
	// Get character
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Currently only Magicians are implemented
	if character.Class != "Magician" {
		return nil, apperrors.NewBadRequest(fmt.Sprintf("Spell preparation not implemented for class '%s'", character.Class))
	}

	return s.getMagicianSpellSlots(ctx, character)
}

// getMagicianSpellSlots calculates available spell slots for Magicians
func (s *SpellPreparationService) getMagicianSpellSlots(ctx context.Context, character *models.Character) (*models.SpellSlots, error) {
	// Get magician data
	magicianData, err := s.magicianService.GetAllMagicianLevelData(ctx)
	if err != nil {
		return nil, err
	}

	// Find character's level data
	var levelData *models.MagicianClassData
	for _, data := range magicianData {
		if data.Level == character.Level {
			levelData = data
			break
		}
	}

	if levelData == nil {
		return nil, apperrors.NewInternalError(fmt.Errorf("could not find magician data for level %d", character.Level))
	}

	// Calculate bonus slots from intelligence
	intBonusLevel := 0
	intBonusSlots := 0

	if character.Intelligence >= 13 && character.Intelligence <= 14 {
		intBonusLevel = 1
		intBonusSlots = 1
	} else if character.Intelligence >= 15 && character.Intelligence <= 16 {
		intBonusLevel = 2
		intBonusSlots = 1
	} else if character.Intelligence == 17 {
		intBonusLevel = 3
		intBonusSlots = 1
	} else if character.Intelligence == 18 {
		intBonusLevel = 4
		intBonusSlots = 1
	}

	// Get counts of prepared spells by level
	level1Used, err := s.preparedSpellRepo.CountPreparedSpellsByLevel(ctx, character.ID, 1)
	if err != nil {
		level1Used = 0
	}
	level2Used, err := s.preparedSpellRepo.CountPreparedSpellsByLevel(ctx, character.ID, 2)
	if err != nil {
		level2Used = 0
	}
	level3Used, err := s.preparedSpellRepo.CountPreparedSpellsByLevel(ctx, character.ID, 3)
	if err != nil {
		level3Used = 0
	}
	level4Used, err := s.preparedSpellRepo.CountPreparedSpellsByLevel(ctx, character.ID, 4)
	if err != nil {
		level4Used = 0
	}
	level5Used, err := s.preparedSpellRepo.CountPreparedSpellsByLevel(ctx, character.ID, 5)
	if err != nil {
		level5Used = 0
	}
	level6Used, err := s.preparedSpellRepo.CountPreparedSpellsByLevel(ctx, character.ID, 6)
	if err != nil {
		level6Used = 0
	}

	// Calculate total slots available
	level1Total := levelData.SpellSlotsLevel1
	level2Total := levelData.SpellSlotsLevel2
	level3Total := levelData.SpellSlotsLevel3
	level4Total := levelData.SpellSlotsLevel4
	level5Total := levelData.SpellSlotsLevel5
	level6Total := levelData.SpellSlotsLevel6

	// Add intelligence bonus slots to appropriate level
	if intBonusLevel == 1 {
		level1Total += intBonusSlots
	} else if intBonusLevel == 2 {
		level2Total += intBonusSlots
	} else if intBonusLevel == 3 {
		level3Total += intBonusSlots
	} else if intBonusLevel == 4 {
		level4Total += intBonusSlots
	}

	// Create spell slots info
	slots := &models.SpellSlots{
		Level1: level1Total - level1Used,
		Level2: level2Total - level2Used,
		Level3: level3Total - level3Used,
		Level4: level4Total - level4Used,
		Level5: level5Total - level5Used,
		Level6: level6Total - level6Used,
		Level7: 0,
		Level8: 0,
		Level9: 0,
	}

	return slots, nil
}

// getClericSpellSlots calculates available spell slots for Clerics
func (s *SpellPreparationService) getClericSpellSlots(ctx context.Context, character *models.Character) (*models.SpellSlots, error) {
	// Placeholder - implement cleric spell slot calculation
	return &models.SpellSlots{}, nil
}

// getDruidSpellSlots calculates available spell slots for Druids
func (s *SpellPreparationService) getDruidSpellSlots(ctx context.Context, character *models.Character) (*models.SpellSlots, error) {
	// Placeholder - implement druid spell slot calculation
	return &models.SpellSlots{}, nil
}

// PrepareSpell prepares a spell for a character
func (s *SpellPreparationService) PrepareSpell(ctx context.Context, characterID int64, input *models.PrepareSpellInput) error {
	// Validate input
	if err := input.Validate(); err != nil {
		return err
	}

	// Ensure characterID matches input
	if characterID != input.CharacterID {
		return apperrors.NewBadRequest("Character ID mismatch")
	}

	// Get character
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return err
	}

	// Check if character is a spellcaster
	if character.Class != "Magician" {
		return apperrors.NewBadRequest(fmt.Sprintf("Spell preparation not implemented for class '%s'", character.Class))
	}

	// Get spell to verify it exists and check its level
	spell, err := s.spellRepo.GetSpell(ctx, input.SpellID)
	if err != nil {
		return err
	}

	// Check if the character has access to this spell
	if !s.characterHasSpellInSpellbook(ctx, character.ID, spell.ID) {
		return apperrors.NewBadRequest(fmt.Sprintf("Character does not have spell '%s' in any spellbook", spell.Name))
	}

	// Check if the slot level is valid for this spell
	if input.SlotLevel < spell.MagLevel {
		return apperrors.NewBadRequest(fmt.Sprintf("Spell '%s' requires at least a level %d slot (requested: %d)",
			spell.Name, spell.MagLevel, input.SlotLevel))
	}

	// Get available slots
	slots, err := s.GetAvailableSpellSlots(ctx, character.ID)
	if err != nil {
		return err
	}

	// Check if there are enough slots available
	var slotsAvailable int
	switch input.SlotLevel {
	case 1:
		slotsAvailable = slots.Level1
	case 2:
		slotsAvailable = slots.Level2
	case 3:
		slotsAvailable = slots.Level3
	case 4:
		slotsAvailable = slots.Level4
	case 5:
		slotsAvailable = slots.Level5
	case 6:
		slotsAvailable = slots.Level6
	default:
		return apperrors.NewBadRequest(fmt.Sprintf("Invalid slot level: %d", input.SlotLevel))
	}

	if slotsAvailable <= 0 {
		return apperrors.NewBadRequest(fmt.Sprintf("No level %d spell slots remaining", input.SlotLevel))
	}

	// Check if the spell is already prepared
	isPrepared, err := s.preparedSpellRepo.IsSpellPrepared(ctx, input.CharacterID, input.SpellID)
	if err != nil {
		return err
	}
	if isPrepared {
		return apperrors.NewBadRequest(fmt.Sprintf("Spell '%s' is already prepared", spell.Name))
	}

	// Prepare the spell
	return s.preparedSpellRepo.PrepareSpell(ctx, input.CharacterID, input.SpellID, input.SlotLevel)
}

// UnprepareSpell removes a prepared spell
func (s *SpellPreparationService) UnprepareSpell(ctx context.Context, characterID int64, spellID int64) error {
	// First check if the spell is actually prepared
	isPrepared, err := s.preparedSpellRepo.IsSpellPrepared(ctx, characterID, spellID)
	if err != nil {
		return err
	}
	if !isPrepared {
		return apperrors.NewBadRequest("Spell is not prepared")
	}

	return s.preparedSpellRepo.UnprepareSpell(ctx, characterID, spellID)
}

// ClearPreparedSpells removes all prepared spells for a character
func (s *SpellPreparationService) ClearPreparedSpells(ctx context.Context, characterID int64) error {
	// First verify the character exists
	_, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return err
	}

	return s.preparedSpellRepo.ClearPreparedSpells(ctx, characterID)
}

// getSlotsAvailableForLevel returns the number of slots available for a given level
func (s *SpellPreparationService) getSlotsAvailableForLevel(slots *models.SpellSlots, level int) int {
	switch level {
	case 1:
		return slots.Level1
	case 2:
		return slots.Level2
	case 3:
		return slots.Level3
	case 4:
		return slots.Level4
	case 5:
		return slots.Level5
	case 6:
		return slots.Level6
	case 7:
		return slots.Level7
	case 8:
		return slots.Level8
	case 9:
		return slots.Level9
	default:
		return 0
	}
}

// getMinLevelForSpell returns the minimum spell level for a class
func (s *SpellPreparationService) getMinLevelForSpell(class string, spell *models.Spell) int {
	switch class {
	case "Magician":
		return spell.MagLevel
	case "Cleric":
		return spell.ClrLevel
	case "Druid":
		return spell.DrdLevel
	default:
		return 0
	}
}

// canCharacterPrepareSpell checks if a character can prepare a specific spell
func (s *SpellPreparationService) canCharacterPrepareSpell(ctx context.Context, character *models.Character, spell *models.Spell) bool {
	// Check if the spell is available for the character's class
	spellLevel := s.getMinLevelForSpell(character.Class, spell)
	if spellLevel == 0 {
		return false
	}

	// For magicians, check if they have the spell in a spellbook
	if character.Class == "Magician" {
		return s.characterHasSpellInSpellbook(ctx, character.ID, spell.ID)
	}

	// Clerics and Druids have access to all spells of their class
	return true
}

// characterHasSpellInSpellbook checks if a character has a spell in any of their spellbooks
func (s *SpellPreparationService) characterHasSpellInSpellbook(ctx context.Context, characterID int64, spellID int64) bool {
	// 1. Get the character's inventory
	inventory, err := s.inventoryRepo.GetInventoryByCharacter(ctx, characterID)
	if err != nil {
		return false
	}

	// 2. Find all spellbooks in the inventory
	for _, item := range inventory.Items {
		if item.ItemType == "spellbook" {
			// Get the spellbook using spellbookRepo
			spellbook, err := s.spellbookRepo.GetSpellbook(ctx, item.ItemID)
			if err != nil {
				continue
			}

			// Get spells in this spellbook
			spellIDs, err := s.spellbookRepo.GetSpellsInSpellbook(ctx, spellbook.ID)
			if err != nil {
				continue
			}

			// Check if the spell is in this spellbook
			for _, id := range spellIDs {
				if id == spellID {
					return true
				}
			}
		}
	}

	return false
}
