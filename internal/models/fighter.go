package models

type FighterClassData struct {
	Level            int    `json:"level"`
	ExperiencePoints int    `json:"experience_points"`
	HitDice          string `json:"hit_dice"`
	SavingThrow      int    `json:"saving_throw"`
	FightingAbility  int    `json:"fighting_ability"`
}

type FighterAbility struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinLevel    int    `json:"min_level"`
}

// GetFighterAbilities returns the predefined fighter abilities
func GetFighterAbilities() []*FighterAbility {
	return []*FighterAbility{
		{
			ID:          1,
			Name:        "Heroic Fighting",
			Description: "To smite multiple foes. When combatting opponents of 1 HD or less, double normal melee attacks per round (2/1, or 3/1 if wielding a mastered weapon). This dramatic attack could be effected as a single, devastating swing or lunge that bursts through multiple foes. At 7th level, when combating foes of 2 HD or less, double normal melee attacks per round (3/1, or 4/1 if wielding a mastered weapon).",
			MinLevel:    1,
		},
		{
			ID:          2,
			Name:        "Weapon Mastery",
			Description: "Mastery of two weapons (+1 \"to hit\" and +1 damage). Additional weapons may be mastered at 4th, 8th, and 12th levels; however, see grand mastery below, for another option. As noted in Chapter 6: Equipment, the attack rate for melee weapons and the rates of fire for most missile weapons improve through weapon mastery. For more information on weapon mastery, see Chapter 9: Combat, weapon skill.",
			MinLevel:    1,
		},
		{
			ID:          3,
			Name:        "Grand Mastery",
			Description: "At 4th, 8th, or 12th level (player's choice), when a new weapon mastery is gained, fighters may elect to intensify their training with an already mastered weapon. With this weapon the fighter becomes a grand master (+2 \"to hit\" and +2 damage, increased attack rate, etc.). A fighter may achieve grand mastery with but one weapon. For more information, see Chapter 9: Combat, weapon skill.",
			MinLevel:    4,
		},
	}
}
