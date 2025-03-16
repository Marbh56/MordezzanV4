package models

import (
	"time"
)

type Character struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	Name         string `json:"name"`
	Class        string `json:"class"`
	Level        int    `json:"level"`
	Strength     int    `json:"strength"`
	Dexterity    int    `json:"dexterity"`
	Constitution int    `json:"constitution"`
	Wisdom       int    `json:"wisdom"`
	Intelligence int    `json:"intelligence"`
	Charisma     int    `json:"charisma"`
	HitPoints    int    `json:"hit_points"`

	// Derived stats
	// Str
	MeleeModifier     int    `json:"Melee_modifier"`
	DamageAdjustment  int    `json:"damage_adjustment"`
	StrengthTest      string `json:"strength_test"`
	ExtraStrengthFeat string `json:"extra_strength_feat"`
	// Dex
	RangedModifier     int    `json:"Ranged_modifier"`
	DefenceAdjustment  int    `json:"defence_adjustment"`
	DexterityTest      string `json:"dexterity_test"`
	ExtraDexterityFeat string `json:"extra_dexterity_feat"`
	// Con
	HPModifier            int    `json:"HP_modifier"`
	PoisonRadModifier     int    `json:"poison_rad_modifier"`
	TraumaSurvival        string `json:"trauma_survival"`
	ConstitutionTest      string `json:"constitution_test"`
	ExtraConstitutionFeat string `json:"extra_constitution_feat"`

	// Int
	LanguageModifier string `json:"language_modifier"`
	MagiciansBonus   string `json:"magicians_bonus"`
	MagiciansChance  string `json:"magicians_chance"`

	// Wis
	WillpowerModifier int    `json:"willpower_modifier"`
	ClericBonus       string `json:"cleric_bonus"`
	ClericChance      string `json:"cleric_chance"`

	// Cha
	ReactionModifier      int `json:"reaction_modifier"`
	MaxFollowers          int `json:"max_followers"`
	UndeadTurningModifier int `json:"undead_turning_modifier"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCharacterInput struct {
	UserID       int64  `json:"user_id"`
	Name         string `json:"name"`
	Class        string `json:"class"`
	Level        int    `json:"level"`
	Strength     int    `json:"strength"`
	Dexterity    int    `json:"dexterity"`
	Constitution int    `json:"constitution"`
	Wisdom       int    `json:"wisdom"`
	Intelligence int    `json:"intelligence"`
	Charisma     int    `json:"charisma"`
	HitPoints    int    `json:"hit_points"`
}

type UpdateCharacterInput struct {
	Name         string `json:"name"`
	Class        string `json:"class"`
	Level        int    `json:"level"`
	Strength     int    `json:"strength"`
	Dexterity    int    `json:"dexterity"`
	Constitution int    `json:"constitution"`
	Wisdom       int    `json:"wisdom"`
	Intelligence int    `json:"intelligence"`
	Charisma     int    `json:"charisma"`
	HitPoints    int    `json:"hit_points"`
}

func (i *CreateCharacterInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Class == "" {
		return NewValidationError("class", "Class cannot be empty")
	}
	if i.UserID <= 0 {
		return NewValidationError("user_id", "Invalid user ID")
	}
	if i.Level < 1 {
		return NewValidationError("level", "Level must be at least 1") // Added validation for level
	}
	if i.Strength < 3 || i.Strength > 18 {
		return NewValidationError("strength", "Strength must be between 3 and 18")
	}
	if i.Dexterity < 3 || i.Dexterity > 18 {
		return NewValidationError("dexterity", "Dexterity must be between 3 and 18")
	}
	if i.Constitution < 3 || i.Constitution > 18 {
		return NewValidationError("constitution", "Constitution must be between 3 and 18")
	}
	if i.Wisdom < 3 || i.Wisdom > 18 {
		return NewValidationError("wisdom", "Wisdom must be between 3 and 18")
	}
	if i.Intelligence < 3 || i.Intelligence > 18 {
		return NewValidationError("intelligence", "Intelligence must be between 3 and 18")
	}
	if i.Charisma < 3 || i.Charisma > 18 {
		return NewValidationError("charisma", "Charisma must be between 3 and 18")
	}
	if i.HitPoints < 1 {
		return NewValidationError("hit_points", "Hit points must be positive")
	}
	return nil
}

func (i *UpdateCharacterInput) Validate() error {
	if i.Name == "" {
		return NewValidationError("name", "Name cannot be empty")
	}
	if i.Class == "" {
		return NewValidationError("class", "Class cannot be empty")
	}
	if i.Level < 1 {
		return NewValidationError("level", "Level must be at least 1") // Added validation for level
	}
	if i.Strength < 3 || i.Strength > 18 {
		return NewValidationError("strength", "Strength must be between 3 and 18")
	}
	if i.Dexterity < 3 || i.Dexterity > 18 {
		return NewValidationError("dexterity", "Dexterity must be between 3 and 18")
	}
	if i.Constitution < 3 || i.Constitution > 18 {
		return NewValidationError("constitution", "Constitution must be between 3 and 18")
	}
	if i.Wisdom < 3 || i.Wisdom > 18 {
		return NewValidationError("wisdom", "Wisdom must be between 3 and 18")
	}
	if i.Intelligence < 3 || i.Intelligence > 18 {
		return NewValidationError("intelligence", "Intelligence must be between 3 and 18")
	}
	if i.Charisma < 3 || i.Charisma > 18 {
		return NewValidationError("charisma", "Charisma must be between 3 and 18")
	}
	if i.HitPoints < 1 {
		return NewValidationError("hit_points", "Hit points must be positive")
	}
	return nil
}

func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (c *Character) CalculateDerivedStats() {
	c.calculateStrengthModifiers()
	c.calculateDexterityModifiers()
	c.calculateConstitutionModifiers()
	c.calculateIntelligenceModifiers()
	c.calculateWisdomModifiers()
	c.calculateCharismaModifiers()

	// Add other stat calculations as needed
}

func (c *Character) calculateStrengthModifiers() {
	switch {
	case c.Strength == 3:
		c.MeleeModifier = -2
		c.DamageAdjustment = -2
		c.StrengthTest = "1:6"
		c.ExtraStrengthFeat = "0%"
	case c.Strength >= 4 && c.Strength <= 6:
		c.MeleeModifier = -1
		c.DamageAdjustment = -1
		c.StrengthTest = "1:6"
		c.ExtraStrengthFeat = "1%"
	case c.Strength >= 7 && c.Strength <= 8:
		c.MeleeModifier = 0
		c.DamageAdjustment = -1
		c.StrengthTest = "2:6"
		c.ExtraStrengthFeat = "2%"
	case c.Strength >= 9 && c.Strength <= 12:
		c.MeleeModifier = 0
		c.DamageAdjustment = 0
		c.StrengthTest = "2:6"
		c.ExtraStrengthFeat = "4%"
	case c.Strength >= 13 && c.Strength <= 14:
		c.MeleeModifier = 0
		c.DamageAdjustment = 1
		c.StrengthTest = "3:6"
		c.ExtraStrengthFeat = "8"
	case c.Strength >= 15 && c.Strength <= 16:
		c.MeleeModifier = 1
		c.DamageAdjustment = 1
		c.StrengthTest = "3:6"
		c.ExtraStrengthFeat = "16%"
	case c.Strength == 17:
		c.MeleeModifier = 1
		c.DamageAdjustment = 2
		c.StrengthTest = "4:6"
		c.ExtraStrengthFeat = "24%"
	case c.Strength == 18:
		c.MeleeModifier = 2
		c.DamageAdjustment = 3
		c.StrengthTest = "5:6"
		c.ExtraStrengthFeat = "32%"
	}
}

func (c *Character) calculateDexterityModifiers() {
	switch {
	case c.Dexterity == 3:
		c.RangedModifier = -2
		c.DefenceAdjustment = -2
		c.DexterityTest = "1:6"
		c.ExtraDexterityFeat = "0%"
	case c.Dexterity >= 4 && c.Dexterity <= 6:
		c.RangedModifier = -1
		c.DefenceAdjustment = -1
		c.DexterityTest = "1:6"
		c.ExtraDexterityFeat = "1%"
	case c.Dexterity >= 7 && c.Dexterity <= 8:
		c.RangedModifier = -1
		c.DefenceAdjustment = 0
		c.DexterityTest = "2:6"
		c.ExtraDexterityFeat = "2%"
	case c.Dexterity >= 9 && c.Dexterity <= 12:
		c.RangedModifier = 0
		c.DefenceAdjustment = 0
		c.DexterityTest = "2:6"
		c.ExtraDexterityFeat = "4%"
	case c.Dexterity >= 13 && c.Dexterity <= 14:
		c.RangedModifier = 1
		c.DefenceAdjustment = 0
		c.DexterityTest = "3:6"
		c.ExtraDexterityFeat = "8%"
	case c.Dexterity >= 15 && c.Dexterity <= 16:
		c.RangedModifier = 1
		c.DefenceAdjustment = 1
		c.DexterityTest = "3:6"
		c.ExtraDexterityFeat = "16%"
	case c.Dexterity == 17:
		c.RangedModifier = 2
		c.DefenceAdjustment = 1
		c.DexterityTest = "4:6"
		c.ExtraDexterityFeat = "24%"
	case c.Dexterity == 18:
		c.RangedModifier = 3
		c.DefenceAdjustment = 2
		c.DexterityTest = "5:6"
		c.ExtraDexterityFeat = "32%"
	}
}

func (c *Character) calculateConstitutionModifiers() {
	switch {
	case c.Constitution == 3:
		c.HPModifier = -1
		c.PoisonRadModifier = -2
		c.TraumaSurvival = "45%"
		c.ConstitutionTest = "1:6"
		c.ExtraConstitutionFeat = "0%"
	case c.Constitution >= 4 && c.Constitution <= 6:
		c.HPModifier = -1
		c.PoisonRadModifier = -1
		c.TraumaSurvival = "55%"
		c.ConstitutionTest = "1:6"
		c.ExtraConstitutionFeat = "1%"
	case c.Constitution >= 7 && c.Constitution <= 8:
		c.HPModifier = 0
		c.PoisonRadModifier = 0
		c.TraumaSurvival = "65%"
		c.ConstitutionTest = "2:6"
		c.ExtraConstitutionFeat = "2%"
	case c.Constitution >= 9 && c.Constitution <= 12:
		c.HPModifier = 0
		c.PoisonRadModifier = 0
		c.TraumaSurvival = "75%"
		c.ConstitutionTest = "2:6"
		c.ExtraConstitutionFeat = "4%"
	case c.Constitution >= 13 && c.Constitution <= 14:
		c.HPModifier = 1
		c.PoisonRadModifier = 0
		c.TraumaSurvival = "80%"
		c.ConstitutionTest = "3:6"
		c.ExtraConstitutionFeat = "8%"
	case c.Constitution >= 15 && c.Constitution <= 16:
		c.HPModifier = 1
		c.PoisonRadModifier = 1
		c.TraumaSurvival = "85%"
		c.ConstitutionTest = "3:6"
		c.ExtraConstitutionFeat = "16%"
	case c.Constitution == 17:
		c.HPModifier = 2
		c.PoisonRadModifier = 1
		c.TraumaSurvival = "90%"
		c.ConstitutionTest = "4:6"
		c.ExtraConstitutionFeat = "24%"
	case c.Constitution == 18:
		c.HPModifier = 3
		c.PoisonRadModifier = 2
		c.TraumaSurvival = "95%"
		c.ConstitutionTest = "5:6"
		c.ExtraConstitutionFeat = "32%"
	}
}

func (c *Character) calculateIntelligenceModifiers() {
	switch {
	case c.Intelligence == 3:
		c.LanguageModifier = "Illiterate"
		c.MagiciansBonus = "N/A"
		c.MagiciansChance = "N/A"
	case c.Intelligence >= 4 && c.Intelligence <= 6:
		c.LanguageModifier = "Illiterate"
		c.MagiciansBonus = "N/A"
		c.MagiciansChance = "N/A"
	case c.Intelligence >= 7 && c.Intelligence <= 8:
		c.LanguageModifier = "0"
		c.MagiciansBonus = "N/A"
		c.MagiciansChance = "N/A"
	case c.Intelligence >= 9 && c.Intelligence <= 12:
		c.LanguageModifier = "0"
		c.MagiciansBonus = "-"
		c.MagiciansChance = "50%"
	case c.Intelligence >= 13 && c.Intelligence <= 14:
		c.LanguageModifier = "+1"
		c.MagiciansBonus = "One Level 1"
		c.MagiciansChance = "65%"
	case c.Intelligence >= 15 && c.Intelligence <= 16:
		c.LanguageModifier = "+1"
		c.MagiciansBonus = "One Level 2"
		c.MagiciansChance = "75%"
	case c.Intelligence == 17:
		c.LanguageModifier = "+2"
		c.MagiciansBonus = "One Level 3"
		c.MagiciansChance = "85%"
	case c.Intelligence == 18:
		c.LanguageModifier = "+3"
		c.MagiciansBonus = "One Level 4"
		c.MagiciansChance = "95%"
	}
}

func (c *Character) calculateWisdomModifiers() {
	switch {
	case c.Wisdom == 3:
		c.WillpowerModifier = -2
		c.ClericBonus = "N/A"
		c.ClericChance = "N/A"
	case c.Wisdom >= 4 && c.Wisdom <= 6:
		c.WillpowerModifier = -1
		c.ClericBonus = "N/A"
		c.ClericChance = "N/A"
	case c.Wisdom >= 7 && c.Wisdom <= 8:
		c.WillpowerModifier = 0
		c.ClericBonus = "-"
		c.ClericChance = "50%"
	case c.Wisdom >= 9 && c.Wisdom <= 12:
		c.WillpowerModifier = 0
		c.ClericBonus = "-"
		c.ClericChance = "50%"
	case c.Wisdom >= 13 && c.Wisdom <= 14:
		c.WillpowerModifier = 1
		c.ClericBonus = "One Level 1"
		c.ClericChance = "65%"
	case c.Wisdom >= 15 && c.Wisdom <= 16:
		c.WillpowerModifier = 1
		c.ClericBonus = "One Level 2"
		c.ClericChance = "75%"
	case c.Wisdom == 17:
		c.WillpowerModifier = 2
		c.ClericBonus = "One Level 3"
		c.ClericChance = "85%"
	case c.Wisdom == 18:
		c.WillpowerModifier = 3
		c.ClericBonus = "One Level 4"
		c.ClericChance = "95%"
	}
}

func (c *Character) calculateCharismaModifiers() {
	switch {
	case c.Charisma == 3:
		c.ReactionModifier = -3
		c.MaxFollowers = 1
		c.UndeadTurningModifier = -1
	case c.Charisma >= 4 && c.Charisma <= 6:
		c.ReactionModifier = -2
		c.MaxFollowers = 2
		c.UndeadTurningModifier = -1
	case c.Charisma >= 7 && c.Charisma <= 8:
		c.ReactionModifier = -1
		c.MaxFollowers = 3
		c.UndeadTurningModifier = 0
	case c.Charisma >= 9 && c.Charisma <= 12:
		c.ReactionModifier = 0
		c.MaxFollowers = 4
		c.UndeadTurningModifier = 0
	case c.Charisma >= 13 && c.Charisma <= 14:
		c.ReactionModifier = 1
		c.MaxFollowers = 6
		c.UndeadTurningModifier = +1
	case c.Charisma >= 15 && c.Charisma <= 16:
		c.ReactionModifier = 1
		c.MaxFollowers = 8
		c.UndeadTurningModifier = +1
	case c.Charisma == 17:
		c.ReactionModifier = 2
		c.MaxFollowers = 10
		c.UndeadTurningModifier = +1
	case c.Charisma == 18:
		c.ReactionModifier = 3
		c.MaxFollowers = 12
		c.UndeadTurningModifier = +1
	}
}
