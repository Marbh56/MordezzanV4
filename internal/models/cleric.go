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

func GetClericAbilities() []*ClericAbility {
	return []*ClericAbility{
		{
			ID:          1,
			Name:        "Scroll Use",
			Description: "The ability to decipher and invoke scrolls with spells from the Cleric Spell List, unless the scroll was created by a thaumaturgical sorcerer (one who casts magician or magician subclass spells).",
			MinLevel:    1,
		},
		{
			ID:          2,
			Name:        "Scroll Writing",
			Description: "To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 500 gp + 100 gp per spell level. Materials vary: Some clerics engrave thin tablets of stone, whereas others use vellum or parchment, a fine quill, and sorcerer's ink, such as sepia. This involved process requires one week per spell level and must be completed on consecrated ground, such as a shrine, fane, or temple.",
			MinLevel:    1,
		},
		{
			ID:          3,
			Name:        "Sorcery",
			Description: "Clerics memorize and cast spells without maintaining spell books; they might bear scriptures in prayer books, sacred scrolls, or graven tablets. Clerics begin with three level 1 spells from the Cleric Spell List, revealed through initiation. They gain three new spells at each level, of castable levels, through spiritual revelation, otherworldly favor, or theological study. These are learned automatically without qualification rolls. Additional spells may be learned outside level training through more arduous processes.",
			MinLevel:    1,
		},
		{
			ID:          4,
			Name:        "Turn Undead",
			Description: "Clerics can exert control over undead and some dæmonic beings within 30 feet by displaying a holy symbol and speaking a commandment. Evil clerics can instead compel their service. Success is determined by cross-referencing the cleric's turning ability (TA) with the Undead Type on a d12 roll. Results: T (turned): 2d6 undead cower/flee for 1 turn; D (destroyed): 2d6 undead are exorcized; UD (ultimate destruction): 1d6+6 undead are exorcized. Clerics with 15+ charisma gain +1 on turning rolls. This ability can be used daily a number of times equal to the character's TA, with one attempt per encounter.",
			MinLevel:    1,
		},
		{
			ID:          5,
			Name:        "New Weapon Skill",
			Description: "At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list. This new proficiency is dependent upon training and practice.",
			MinLevel:    4,
		},
	}
}
