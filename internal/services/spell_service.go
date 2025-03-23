package services

import (
	"context"
	"fmt"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"strconv"
)

// SpellService handles business logic for spells and spell casting
type SpellService struct {
	spellRepo          repositories.SpellRepository
	spellCastingRepo   repositories.SpellCastingRepository
	characterRepo      repositories.CharacterRepository
	classRepo          repositories.ClassRepository
	classService       *ClassService
	encumbranceService *EncumbranceService
}

// NewSpellService creates a new spell service
func NewSpellService(
	spellRepo repositories.SpellRepository,
	spellCastingRepo repositories.SpellCastingRepository,
	characterRepo repositories.CharacterRepository,
	classRepo repositories.ClassRepository,
	classService *ClassService,
	encumbranceService *EncumbranceService,
) *SpellService {
	return &SpellService{
		spellRepo:          spellRepo,
		spellCastingRepo:   spellCastingRepo,
		characterRepo:      characterRepo,
		classRepo:          classRepo,
		classService:       classService,
		encumbranceService: encumbranceService,
	}
}

// GetCharacterSpellsInfo retrieves all spell-related information for a character
func (s *SpellService) GetCharacterSpellsInfo(ctx context.Context, characterID int64) (*models.CharacterSpellsInfo, error) {
	// Get character details
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	// Get known spells
	knownSpells, err := s.spellCastingRepo.GetKnownSpells(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get known spells: %v", err)
	}

	// Get prepared spells
	preparedSpells, err := s.spellCastingRepo.GetPreparedSpells(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get prepared spells: %v", err)
	}

	// Get base spell slots
	spellSlots, err := s.spellCastingRepo.GetSpellSlots(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spell slots: %v", err)
	}

	// Get bonus spells
	bonusSpells, err := s.spellCastingRepo.CalculateBonusSpells(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate bonus spells: %v", err)
	}

	// Get spell limits based on class and level
	classSpellLimits := make(map[string]map[string][]int)

	// Determine primary casting class
	primaryCastingClass := getPrimaryCastingClass(character.Class)
	if primaryCastingClass != "" {
		maxKnownSpells, err := s.spellCastingRepo.GetMaxKnownSpells(ctx, characterID, primaryCastingClass)
		if err != nil {
			return nil, fmt.Errorf("failed to get max known spells: %v", err)
		}

		classSpellLimits[primaryCastingClass] = make(map[string][]int)
		for level, max := range maxKnownSpells {
			classSpellLimits[primaryCastingClass][level] = []int{max, 0}
		}
	}

	// Count current known spells by class and level
	for _, spell := range knownSpells {
		spellLevelStr := fmt.Sprintf("level%d", spell.SpellLevel)
		if limits, ok := classSpellLimits[spell.SpellClass]; ok {
			if levelCounts, ok := limits[spellLevelStr]; ok {
				levelCounts[1]++
				classSpellLimits[spell.SpellClass][spellLevelStr] = levelCounts
			}
		}
	}

	// Calculate available prepared slots
	availablePreparedSlots := make(map[string]int)
	for level, count := range spellSlots {
		// Add base slots
		availablePreparedSlots[level] = count

		// Add bonus slots from attributes if the character has the casting ability for this level
		levelNum, _ := strconv.Atoi(level[5:])

		// For divine casters
		if isDivineCaster(character.Class) && bonusSpells["divine"] != nil {
			if character.Level >= (2*levelNum - 1) { // Check if character can cast this level
				if bonus, ok := bonusSpells["divine"][level]; ok {
					availablePreparedSlots[level] += bonus
				}
			}
		}

		// For arcane casters
		if isArcaneCaster(character.Class) && bonusSpells["arcane"] != nil {
			if character.Level >= (2*levelNum - 1) { // Check if character can cast this level
				if bonus, ok := bonusSpells["arcane"][level]; ok {
					availablePreparedSlots[level] += bonus
				}
			}
		}
	}

	// Count prepared spells and subtract from available slots
	for _, spell := range preparedSpells {
		levelKey := fmt.Sprintf("level%d", spell.SpellLevel)
		if count, ok := availablePreparedSlots[levelKey]; ok && count > 0 {
			availablePreparedSlots[levelKey]--
		}
	}

	return &models.CharacterSpellsInfo{
		KnownSpells:            knownSpells,
		PreparedSpells:         preparedSpells,
		SpellSlots:             spellSlots,
		AvailablePreparedSlots: availablePreparedSlots,
		BonusSpells:            bonusSpells,
		ClassSpellLimits:       classSpellLimits,
	}, nil
}

func (s *SpellService) ClearPreparedSpells(ctx context.Context, characterID int64) error {
	return s.spellCastingRepo.ClearPreparedSpells(ctx, characterID)
}

func (s *SpellService) UnprepareSpell(ctx context.Context, characterID, spellID int64) error {
	return s.spellCastingRepo.UnprepareSpell(ctx, characterID, spellID)
}

func (s *SpellService) PrepareSpell(ctx context.Context, input *models.PrepareSpellInput) (int64, error) {
	return s.spellCastingRepo.PrepareSpell(ctx, input)
}

func (s *SpellService) RemoveKnownSpell(ctx context.Context, characterID, spellID int64) error {
	return s.spellCastingRepo.RemoveKnownSpell(ctx, characterID, spellID)
}

// AddInitialSpellsForNewCharacter adds the starting spells for a new character based on class
func (s *SpellService) AddInitialSpellsForNewCharacter(ctx context.Context, characterID int64) error {
	// Get character details
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return fmt.Errorf("failed to get character: %v", err)
	}

	// Determine if this character can cast spells
	primaryCastingClass := getPrimaryCastingClass(character.Class)
	if primaryCastingClass == "" {
		// No spellcasting for this class
		return nil
	}

	// Get spells available for this class
	spells, err := s.spellRepo.GetSpellsByClass(ctx, primaryCastingClass)
	if err != nil {
		return fmt.Errorf("failed to get spells for class %s: %v", primaryCastingClass, err)
	}

	// Filter to level 1 spells
	var level1Spells []*models.Spell
	for _, spell := range spells {
		spellLevel := spell.GetLevel(primaryCastingClass)
		if spellLevel == 1 {
			level1Spells = append(level1Spells, spell)
		}
	}

	// Characters start with 3 known spells of level 1
	initialSpellCount := 3
	if len(level1Spells) < initialSpellCount {
		initialSpellCount = len(level1Spells)
	}

	// Add the initial spells
	for i := 0; i < initialSpellCount; i++ {
		spell := level1Spells[i]
		input := &models.AddKnownSpellInput{
			CharacterID: characterID,
			SpellID:     spell.ID,
			SpellClass:  primaryCastingClass,
			Notes:       "Initial spell",
		}

		_, err := s.spellCastingRepo.AddKnownSpell(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to add initial spell %s: %v", spell.Name, err)
		}
	}

	return nil
}

// GetSpellsLearnableOnLevelUp gets spells a character can learn when leveling up
func (s *SpellService) GetSpellsLearnableOnLevelUp(ctx context.Context, characterID int64, newLevel int) ([]*models.Spell, error) {
	// Get character details
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %v", err)
	}

	// Determine if this character can cast spells
	primaryCastingClass := getPrimaryCastingClass(character.Class)
	if primaryCastingClass == "" {
		// No spellcasting for this class
		return nil, nil
	}

	// Calculate the highest spell level this character can cast at the new level
	maxSpellLevel := calculateMaxSpellLevel(primaryCastingClass, newLevel)

	// Get spells available for this class up to the max level
	allSpells, err := s.spellRepo.GetSpellsByClass(ctx, primaryCastingClass)
	if err != nil {
		return nil, fmt.Errorf("failed to get spells for class %s: %v", primaryCastingClass, err)
	}

	// Get the character's known spells
	knownSpells, err := s.spellCastingRepo.GetKnownSpells(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get known spells: %v", err)
	}

	// Create a map of known spell IDs for easy lookup
	knownSpellIDs := make(map[int64]bool)
	for _, spell := range knownSpells {
		knownSpellIDs[spell.SpellID] = true
	}

	// Filter to spells of appropriate level that aren't already known
	var learnableSpells []*models.Spell
	for _, spell := range allSpells {
		spellLevel := spell.GetLevel(primaryCastingClass)
		if spellLevel <= maxSpellLevel && !knownSpellIDs[spell.ID] {
			learnableSpells = append(learnableSpells, spell)
		}
	}

	return learnableSpells, nil
}

// LearnSpellOnLevelUp adds a spell to a character's known spells when leveling up
func (s *SpellService) LearnSpellOnLevelUp(ctx context.Context, characterID int64, spellID int64) error {
	// Get character details
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return fmt.Errorf("failed to get character: %v", err)
	}

	// Determine primary casting class
	primaryCastingClass := getPrimaryCastingClass(character.Class)
	if primaryCastingClass == "" {
		return fmt.Errorf("character class %s cannot cast spells", character.Class)
	}

	// Get the spell details
	spell, err := s.spellRepo.GetSpell(ctx, spellID)
	if err != nil {
		return fmt.Errorf("failed to get spell: %v", err)
	}

	// Get spell level for this class
	spellLevel := spell.GetLevel(primaryCastingClass)

	// Calculate the highest spell level this character can cast
	maxSpellLevel := calculateMaxSpellLevel(primaryCastingClass, character.Level)
	if spellLevel > maxSpellLevel {
		return fmt.Errorf("character cannot learn spells of level %d yet", spellLevel)
	}

	// Check if the character already knows the maximum number of spells of this level
	maxKnownSpells, err := s.spellCastingRepo.GetMaxKnownSpells(ctx, characterID, primaryCastingClass)
	if err != nil {
		return fmt.Errorf("failed to get max known spells: %v", err)
	}

	// Get the character's known spells of this level
	knownSpellsByClass, err := s.spellCastingRepo.GetKnownSpellsByClass(ctx, characterID, primaryCastingClass)
	if err != nil {
		return fmt.Errorf("failed to get known spells: %v", err)
	}

	// Count spells of this level
	levelSpellCount := 0
	for _, knownSpell := range knownSpellsByClass {
		if knownSpell.SpellLevel == spellLevel {
			levelSpellCount++
		}
	}

	// Check if the character can learn more spells of this level
	levelKey := fmt.Sprintf("level%d", spellLevel)
	if levelSpellCount >= maxKnownSpells[levelKey] {
		return fmt.Errorf("character cannot learn more level %d spells for class %s", spellLevel, primaryCastingClass)
	}

	// Add the spell to the character's known spells
	input := &models.AddKnownSpellInput{
		CharacterID: characterID,
		SpellID:     spellID,
		SpellClass:  primaryCastingClass,
		Notes:       fmt.Sprintf("Learned at level %d", character.Level),
	}

	_, err = s.spellCastingRepo.AddKnownSpell(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to add spell to known spells: %v", err)
	}

	return nil
}

// PrepareAllSpells prepares all spells a character can prepare
func (s *SpellService) PrepareAllSpells(ctx context.Context, characterID int64) error {
	// Get character details and spell info
	spellInfo, err := s.GetCharacterSpellsInfo(ctx, characterID)
	if err != nil {
		return fmt.Errorf("failed to get character spell info: %v", err)
	}

	// Clear current prepared spells
	err = s.spellCastingRepo.ClearPreparedSpells(ctx, characterID)
	if err != nil {
		return fmt.Errorf("failed to clear prepared spells: %v", err)
	}

	// Group known spells by level and class
	spellsByLevelAndClass := make(map[string]map[int][]models.KnownSpell)
	for _, spell := range spellInfo.KnownSpells {
		if spellsByLevelAndClass[spell.SpellClass] == nil {
			spellsByLevelAndClass[spell.SpellClass] = make(map[int][]models.KnownSpell)
		}
		spellsByLevelAndClass[spell.SpellClass][spell.SpellLevel] = append(
			spellsByLevelAndClass[spell.SpellClass][spell.SpellLevel],
			spell,
		)
	}

	// For each spell level, prepare as many spells as possible
	for classKey, spellsByLevel := range spellsByLevelAndClass {
		for level, spells := range spellsByLevel {
			// Get available slots for this level
			slotKey := fmt.Sprintf("level%d", level)
			availableSlots := spellInfo.AvailablePreparedSlots[slotKey]

			// Prepare as many spells as we have slots for
			for i := 0; i < len(spells) && i < availableSlots; i++ {
				spell := spells[i]
				input := &models.PrepareSpellInput{
					CharacterID: characterID,
					SpellID:     spell.SpellID,
					SpellLevel:  spell.SpellLevel,
					SpellClass:  classKey,
				}

				_, err := s.spellCastingRepo.PrepareSpell(ctx, input)
				if err != nil {
					return fmt.Errorf("failed to prepare spell %s: %v", spell.SpellName, err)
				}
			}
		}
	}

	return nil
}

// Private helper functions

// getPrimaryCastingClass determines the primary spellcasting class for a character
func getPrimaryCastingClass(characterClass string) string {
	// Map character classes to their primary spellcasting class
	castingClassMap := map[string]string{
		// Pure divine casters
		"Cleric": "Cleric",
		"Druid":  "Druid",
		"Priest": "Priest",

		// Pure arcane casters
		"Magician":    "Magician",
		"Illusionist": "Illusionist",
		"Necromancer": "Necromancer",
		"Pyromancer":  "Pyromancer",
		"Cryomancer":  "Cryomancer",
		"Warlock":     "Warlock",
		"Witch":       "Witch",

		// Hybrid classes - could have different primary casting classes
		"Paladin": "Cleric",   // Divine
		"Ranger":  "Druid",    // Divine
		"Bard":    "Magician", // Arcane
		"Shaman":  "Druid",    // Divine

		// Non-casting classes
		"Fighter":        "",
		"Barbarian":      "",
		"Berserker":      "",
		"Cataphract":     "",
		"Huntsman":       "",
		"Monk":           "",
		"Thief":          "",
		"Assassin":       "",
		"Legerdemainist": "",
		"Purloiner":      "",
		"Scout":          "",
	}

	if castingClass, ok := castingClassMap[characterClass]; ok {
		return castingClass
	}

	return ""
}

// isDivineCaster checks if a class is a divine spellcaster
func isDivineCaster(class string) bool {
	divineCasters := []string{"Cleric", "Druid", "Priest", "Paladin", "Shaman"}
	for _, c := range divineCasters {
		if c == class {
			return true
		}
	}
	return false
}

// isArcaneCaster checks if a class is an arcane spellcaster
func isArcaneCaster(class string) bool {
	arcaneCasters := []string{"Magician", "Illusionist", "Necromancer", "Pyromancer", "Cryomancer", "Warlock", "Witch", "Bard"}
	for _, c := range arcaneCasters {
		if c == class {
			return true
		}
	}
	return false
}

// calculateMaxSpellLevel determines the highest spell level a character can cast
func calculateMaxSpellLevel(castingClass string, characterLevel int) int {
	if isDivineCaster(castingClass) {
		// Divine casters (clerics, etc.)
		switch {
		case characterLevel >= 17:
			return 6
		case characterLevel >= 15:
			return 5
		case characterLevel >= 9:
			return 4
		case characterLevel >= 5:
			return 3
		case characterLevel >= 3:
			return 2
		case characterLevel >= 1:
			return 1
		default:
			return 0
		}
	} else if isArcaneCaster(castingClass) {
		// Arcane casters (magicians, etc.)
		switch {
		case characterLevel >= 17:
			return 6
		case characterLevel >= 14:
			return 5
		case characterLevel >= 11:
			return 4
		case characterLevel >= 7:
			return 3
		case characterLevel >= 3:
			return 2
		case characterLevel >= 1:
			return 1
		default:
			return 0
		}
	}

	return 0
}
