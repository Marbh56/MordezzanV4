package models

// ThiefSkillProgression represents a progression of a thief skill at different level ranges
type ThiefSkillProgression struct {
	SkillID       int64  `json:"skill_id"`
	SkillName     string `json:"skill_name"`
	Attribute     string `json:"attribute"`
	LevelRange    string `json:"level_range"`
	SuccessChance string `json:"success_chance"`
}

// ThiefSkillWithChance represents a thief skill with its success chance for a specific level
type ThiefSkillWithChance struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Attribute     string `json:"attribute"`
	SuccessChance string `json:"success_chance"`
}
