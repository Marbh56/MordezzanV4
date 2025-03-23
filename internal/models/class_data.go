package models

// ClassData represents a character class's level-specific data
type ClassData struct {
	ID               int64          `json:"id"`
	ClassName        string         `json:"class_name"`
	Level            int            `json:"level"`
	ExperiencePoints int            `json:"experience_points"`
	HitDice          string         `json:"hit_dice"`
	SavingThrow      int            `json:"saving_throw"`
	FightingAbility  int            `json:"fighting_ability"`
	CastingAbility   int            `json:"casting_ability,omitempty"`
	SpellSlots       map[string]int `json:"spell_slots,omitempty"`
}

// ClassAbility represents a class-specific ability
type ClassAbility struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinLevel    int    `json:"min_level"`
}


type SpellSlots struct {
	Level1 int `json:"level1"`
	Level2 int `json:"level2"`
	Level3 int `json:"level3"`
	Level4 int `json:"level4"`
	Level5 int `json:"level5"`
	Level6 int `json:"level6"`
	Level7 int `json:"level7"`
	Level8 int `json:"level8"`
	Level9 int `json:"level9"`
}
