package models

type ClericClassData struct {
	Level            int    `json:"level"`
	ExperiencePoints int    `json:"experience_points"`
	HitDice          string `json:"hit_dice"`
	SavingThrow      int    `json:"saving_throw"`
	FightingAbility  int    `json:"fighting_ability"`
	TurningAbility   int    `json:"turning_ability"`
	SpellSlotsLevel1 int    `json:"spell_slots_level1"`
	SpellSlotsLevel2 int    `json:"spell_slots_level2"`
	SpellSlotsLevel3 int    `json:"spell_slots_level3"`
	SpellSlotsLevel4 int    `json:"spell_slots_level4"`
	SpellSlotsLevel5 int    `json:"spell_slots_level5"`
}

type ClericAbility struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinLevel    int    `json:"min_level"`
}

// GetClericAbilities returns the predefined cleric abilities
func GetClericAbilities() []*ClericAbility {
	return []*ClericAbility{
		{
			ID:          1,
			Name:        "Turn Undead",
			Description: "The ability to repel and destroy undead creatures by channeling divine energy. The effectiveness increases with the cleric's level and is affected by Charisma.",
			MinLevel:    1,
		},
		{
			ID:          2,
			Name:        "Divine Spellcasting",
			Description: "Starting at 2nd level, clerics can cast divine spells. The number of spells available increases with level, and additional spells may be granted based on high Wisdom scores.",
			MinLevel:    2,
		},
		{
			ID:          3,
			Name:        "Healing Hands",
			Description: "Once per day per three levels (rounded up), the cleric can heal 1d6+1 hit points with a touch. This healing increases to 2d6+2 at 6th level and 3d6+3 at 9th level.",
			MinLevel:    3,
		},
		{
			ID:          4,
			Name:        "Divine Blessing",
			Description: "At 6th level, the cleric can bestow a blessing on a person or object once per day. This blessing grants a +1 bonus to saves, attacks, or defense for 1 hour.",
			MinLevel:    6,
		},
		{
			ID:          5,
			Name:        "Divine Intervention",
			Description: "At 9th level, once per week, the cleric may call upon their deity for direct intervention in dire circumstances. Success is determined by the Referee based on the situation and the cleric's devotion.",
			MinLevel:    9,
		},
	}
}