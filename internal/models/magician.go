package models

// MagicianClassData represents the level-specific data for magicians
type MagicianClassData struct {
	Level            int    `json:"level"`
	ExperiencePoints int    `json:"experience_points"`
	HitDice          string `json:"hit_dice"`
	SavingThrow      int    `json:"saving_throw"`
	FightingAbility  int    `json:"fighting_ability"`
	CastingAbility   int    `json:"casting_ability"`
	SpellSlotsLevel1 int    `json:"spell_slots_level_1"`
	SpellSlotsLevel2 int    `json:"spell_slots_level_2"`
	SpellSlotsLevel3 int    `json:"spell_slots_level_3"`
	SpellSlotsLevel4 int    `json:"spell_slots_level_4"`
	SpellSlotsLevel5 int    `json:"spell_slots_level_5"`
	SpellSlotsLevel6 int    `json:"spell_slots_level_6"`
}

// MagicianAbility represents a magician ability
type MagicianAbility struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinLevel    int    `json:"min_level"`
}

// GetMagicianAbilities returns the predefined magician abilities
func GetMagicianAbilities() []*MagicianAbility {
	return []*MagicianAbility{
		{
			ID:          1,
			Name:        "Alchemy",
			Description: "To practice the sorcery-science of alchemy. Apprentice magicians learn how to identify potions by taste alone; albeit the practice is not always safe. At 7th level, a magician may concoct potions with the assistance of an alchemist. By 11th level, the assistance of an alchemist is no longer required.",
			MinLevel:    1,
		},
		{
			ID:          2,
			Name:        "Familiar",
			Description: "To summon a small animal (bat, cat, owl, rat, raven, snake, etc.) of 1d3+1 hp to function as a familiar. Retaining a familiar provides the following benefits: Within range 120 (feet indoors, yards outdoors), the magician can see and hear through the animal; sight is narrowly focused, sounds reverberate metallically. The hit point total of the familiar is added to the magician's total. The magician can memorize one extra spell of each available spell level per day. These benefits are lost if the familiar is rendered dead, unconscious, or out of range. The familiar is an extraordinary example of the species, has a perfect morale score (ML 12), and always attends and abides the will of its master. To summon a familiar, the magician must perform a series of rites and rituals for 24 hours. If the familiar dies, the magician also sustains 3d6 hp damage. The magician cannot summon another familiar for 1d4+2 months.",
			MinLevel:    1,
		},
		{
			ID:          3,
			Name:        "Read Magic",
			Description: "The ability to decipher otherwise unintelligible magical inscriptions or symbols placed on weapons, armour, items, doors, walls, and other media by means of the sorcerer mark spell or other like methods.",
			MinLevel:    1,
		},
		{
			ID:          4,
			Name:        "Scroll Use",
			Description: "To decipher and invoke scrolls with spells from the Magician Spell List, unless the scroll was created by an ecclesiastical sorcerer (one who casts cleric or druid spells).",
			MinLevel:    1,
		},
		{
			ID:          5,
			Name:        "Scroll Writing",
			Description: "To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 500 gp + 100 gp per spell level. Materials may include the finest vellum, paper, or papyrus; a fresh quill; and sorcerer's ink, such as sepia. This involved process requires one week per spell level.",
			MinLevel:    1,
		},
		{
			ID:          6,
			Name:        "New Weapon Skill",
			Description: "At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list. This new proficiency is dependent upon training and practice.",
			MinLevel:    4,
		},
	}
}
