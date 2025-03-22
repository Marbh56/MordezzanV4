package models

type ThiefClassData struct {
	Level            int    `json:"level"`
	ExperiencePoints int    `json:"experience_points"`
	HitDice          string `json:"hit_dice"`
	SavingThrow      int    `json:"saving_throw"`
	FightingAbility  int    `json:"fighting_ability"`
}

type ThiefSkills struct {
	Climb              string `json:"climb"`               // Climb (DX)
	DecipherScript     string `json:"decipher_script"`     // Decipher Script (IN)
	DiscernNoise       string `json:"discern_noise"`       // Discern Noise (WS)
	Hide               string `json:"hide"`                // Hide (DX)
	ManipulateTraps    string `json:"manipulate_traps"`    // Manipulate Traps (DX)
	MoveSilently       string `json:"move_silently"`       // Move Silently (DX)
	OpenLocks          string `json:"open_locks"`          // Open Locks (DX)
	PickPockets        string `json:"pick_pockets"`        // Pick Pockets (DX)
	ReadScrolls        string `json:"read_scrolls"`        // Read Scrolls (IN)
	BackstabMultiplier int    `json:"backstab_multiplier"` // Backstab damage multiplier
}

type ThiefAbility struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinLevel    int    `json:"min_level"`
}

// GetThiefSkillsByLevel returns the thief skills for a given level
// using the X:12 system from Hyperborea 3E
func GetThiefSkillsByLevel(level int) ThiefSkills {
	skills := ThiefSkills{}

	// According to Table 16: Progressive Thief Abilities in Hyperborea 3E
	switch {
	case level <= 2:
		skills = ThiefSkills{
			Climb:              "8:12",
			DecipherScript:     "0:12",
			DiscernNoise:       "4:12",
			Hide:               "5:12",
			ManipulateTraps:    "3:12",
			MoveSilently:       "5:12",
			OpenLocks:          "3:12",
			PickPockets:        "4:12",
			ReadScrolls:        "—", // Not available at this level
			BackstabMultiplier: 2,
		}
	case level <= 4:
		skills = ThiefSkills{
			Climb:              "8:12",
			DecipherScript:     "1:12",
			DiscernNoise:       "5:12",
			Hide:               "6:12",
			ManipulateTraps:    "4:12",
			MoveSilently:       "6:12",
			OpenLocks:          "4:12",
			PickPockets:        "5:12",
			ReadScrolls:        "—", // Not available at this level
			BackstabMultiplier: 2,
		}
	case level <= 6:
		skills = ThiefSkills{
			Climb:              "9:12",
			DecipherScript:     "2:12",
			DiscernNoise:       "6:12",
			Hide:               "7:12",
			ManipulateTraps:    "5:12",
			MoveSilently:       "7:12",
			OpenLocks:          "5:12",
			PickPockets:        "6:12",
			ReadScrolls:        "0:12", // First available at level 5
			BackstabMultiplier: 3,
		}
	case level <= 8:
		skills = ThiefSkills{
			Climb:              "9:12",
			DecipherScript:     "3:12",
			DiscernNoise:       "7:12",
			Hide:               "8:12",
			ManipulateTraps:    "6:12",
			MoveSilently:       "8:12",
			OpenLocks:          "6:12",
			PickPockets:        "7:12",
			ReadScrolls:        "3:12",
			BackstabMultiplier: 3,
		}
	case level <= 10:
		skills = ThiefSkills{
			Climb:              "10:12",
			DecipherScript:     "4:12",
			DiscernNoise:       "8:12",
			Hide:               "9:12",
			ManipulateTraps:    "7:12",
			MoveSilently:       "9:12",
			OpenLocks:          "7:12",
			PickPockets:        "8:12",
			ReadScrolls:        "4:12",
			BackstabMultiplier: 4,
		}
	case level <= 12:
		skills = ThiefSkills{
			Climb:              "10:12",
			DecipherScript:     "5:12",
			DiscernNoise:       "9:12",
			Hide:               "10:12",
			ManipulateTraps:    "8:12",
			MoveSilently:       "10:12",
			OpenLocks:          "8:12",
			PickPockets:        "9:12",
			ReadScrolls:        "5:12",
			BackstabMultiplier: 5,
		}
	}

	return skills
}

// GetThiefAbilities returns the predefined thief abilities
func GetThiefAbilities() []*ThiefAbility {
	return []*ThiefAbility{
		{
			ID:          1,
			Name:        "Backstab",
			Description: "A backstab attempt with a class 1 or 2 melee weapon. The target must be unaware of the attack, which may be the result of hiding or moving silently. Also, the target must have vital organs (e.g., skeleton, zombie exempt) and a discernible \"back\" (e.g., green slime, purple worm exempt). If the requirements are met: 1) The attack roll is made at +4 \"to hit.\" 2) Additional weapon damage dice are rolled according to the thief's level: 1st-4th level = ×2, 5th-8th level = ×3, 9th-12th level = ×4. Other damage modifiers (strength, sorcery, etc.) are added afterwards.",
			MinLevel:    1,
		},
		{
			ID:          2,
			Name:        "Thieves' Cant",
			Description: "The secret language of thieves, a strange pidgin in which some words may be unintelligible to an ignorant listener, whereas others might be common yet of alternative meaning. This covert tongue is used in conjunction with specific body language, hand gestures, and facial expressions. Two major dialects of Thieves' Cant are used in Hyperborea: one by city thieves, the other by pirates; commonalities exist betwixt the two.",
			MinLevel:    1,
		},
		{
			ID:          3,
			Name:        "Agile",
			Description: "+1 AC bonus when unarmoured and unencumbered (small shield allowed).",
			MinLevel:    1,
		},
		{
			ID:          5,
			Name:        "New Weapon Skill",
			Description: "At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list. This new proficiency is dependent upon training and practice.",
			MinLevel:    4,
		},
		{
			ID:          6,
			Name:        "READ SCROLLS",
			Description: "Starting at 5th level, thieves can attempt to use magic-user scrolls with a chance of success that increases with level. Failure may result in a backfire effect.",
			MinLevel:    5,
		},
	}
}
