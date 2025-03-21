async function fetchCharacterDetails(charId) {
    try {
        const token = localStorage.getItem('authToken');

        const response = await fetch(`/api/characters/${charId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to fetch character details');
        }

        const character = await response.json();
        displayCharacterDetails(character);

        // Show/hide spells tab based on character class
        toggleSpellsTab(character.class);
    } catch (error) {
        console.error('Error fetching character details:', error);
        alert('Failed to load character details: ' + error.message);
    }
}

function displayCharacterDetails(character) {
    // Set page title
    document.title = `${character.name} - Mordezzan`;

    // Update header
    document.getElementById('character-name').textContent = character.name;
    document.getElementById('character-level-class').textContent = `Level ${character.level} ${character.class}`;

    // Update basic info
    document.getElementById('info-name').textContent = character.name;
    document.getElementById('info-class').textContent = character.class;
    document.getElementById('info-level').textContent = character.level;
    document.getElementById('info-hit-dice').textContent = character.hit_dice || '-';

    window.characterData = character;

    // Update global variables for current character state
    currentHP = character.hit_points;
    currentXP = character.experience_points || 0;
    currentClass = character.class;
    currentLevel = character.level;

    // Update HP display and set field
    updateHPDisplay(character.hit_points);
    document.getElementById('hp-set').value = character.hit_points;

    // Update XP display and set field
    updateXPDisplay(character.experience_points || 0, character.level, character.class);

    // Update saving throws if applicable
    const savingThrowsSection = document.getElementById('saving-throws-section');
    if (character.saving_throw) {
        updateSavingThrows(character);
        savingThrowsSection.style.display = 'block';
    } else {
        savingThrowsSection.style.display = 'none';
    }

    // Update ability scores
    updateAbilityScores(character);

    // Display class abilities
    displayClassAbilities(character);

    // Setup notes tab
    setupNotesTab();
    
    // If character is a Fighter, fetch the fighter levels data
    if (character.class === 'Fighter') {
        fetchFighterLevels();
    }
}

// Implementing updateHPDisplay function as it was defined in character_detail.js
function updateHPDisplay(hp) {
    document.getElementById('info-hp').textContent = hp;
    document.getElementById('hp-set').value = hp;
    
    // Optional: Add visual indicators for low HP
    const hpElement = document.getElementById('info-hp');
    
    if (hp <= 0) {
        hpElement.style.color = '#e74c3c'; // Red for 0 HP
    } else if (hp <= 5) {
        hpElement.style.color = '#e67e22'; // Orange for low HP
    } else {
        hpElement.style.color = ''; // Default color
    }
}

// Implementing updateXPDisplay function as it was defined in character_detail.js
function updateXPDisplay(xp, level, characterClass) {
    document.getElementById('info-experience').textContent = xp;
    document.getElementById('xp-set').value = xp;
    
    // Update fighter level info if applicable
    const fighterLevelInfo = document.getElementById('fighter-level-info');
    
    if (characterClass === 'Fighter' && fighterLevelInfo) {
        // Show fighter-specific level information
        fighterLevelInfo.style.display = 'block';
        
        // Update level displays
        document.getElementById('current-level').textContent = level;
        
        // Find XP needed for next level
        let nextLevelXP = 0;
        let xpNeeded = 0;
        
        if (fighterLevels && fighterLevels.length > 0) {
            // Find the next level threshold
            for (let i = 0; i < fighterLevels.length; i++) {
                if (fighterLevels[i].level > level) {
                    nextLevelXP = fighterLevels[i].xp;
                    break;
                }
            }
            
            // Calculate XP needed
            xpNeeded = Math.max(0, nextLevelXP - xp);
            
            // Update the display
            if (xpNeeded > 0) {
                document.getElementById('next-level-xp').textContent = nextLevelXP;
                document.getElementById('xp-needed').textContent = xpNeeded;
            } else {
                document.getElementById('next-level-xp').innerHTML = '<i>Maximum level reached</i>';
                document.getElementById('xp-needed').textContent = '0';
            }
        }
    } else if (fighterLevelInfo) {
        // Hide fighter level info for non-fighters
        fighterLevelInfo.style.display = 'none';
    }
}

function updateSavingThrows(character) {
    // Implementation from original file
    // This is a placeholder - you'll need to implement based on your actual requirements
    console.log("Updating saving throws for character:", character.name);
    
    // Example implementation:
    const baseSave = character.saving_throw || 15;
    document.getElementById('base_save').textContent = baseSave;
    
    // Set modifiers if provided
    if (character.saving_throw_modifiers) {
        const mods = character.saving_throw_modifiers;
        document.getElementById('death-mod').textContent = formatModifier(mods.death || 0);
        document.getElementById('transform-mod').textContent = formatModifier(mods.transform || 0);
        document.getElementById('device-mod').textContent = formatModifier(mods.device || 0);
        document.getElementById('avoidance-mod').textContent = formatModifier(mods.avoidance || 0);
        document.getElementById('sorcery-mod').textContent = formatModifier(mods.sorcery || 0);
    }
}

function updateAbilityScores(character) {
    // Strength
    document.getElementById('ability-str').textContent = character.strength;
    document.getElementById('str-attack-mod').textContent = formatModifier(character.melee_modifier || 0);
    document.getElementById('str-damage-adj').textContent = formatModifier(character.damage_adjustment || 0);
    document.getElementById('str-test').textContent = character.strength_test || '-';
    document.getElementById('str-feat').textContent = character.extra_strength_feat || '-';

    // Dexterity
    if (document.getElementById('ability-dex')) {
        document.getElementById('ability-dex').textContent = character.dexterity;
        document.getElementById('dex-ranged-mod').textContent = formatModifier(character.ranged_modifier || 0);
        document.getElementById('dex-defense-adj').textContent = formatModifier(character.defence_adjustment || 0);
        document.getElementById('dex-test').textContent = character.dexterity_test || '-';
        document.getElementById('dex-feat').textContent = character.extra_dexterity_feat || '-';
    }
    
    // Constitution
    if (document.getElementById('ability-con')) {
        document.getElementById('ability-con').textContent = character.constitution;
        document.getElementById('con-hp-mod').textContent = formatModifier(character.hp_modifier || 0);
        document.getElementById('con-poison-mod').textContent = formatModifier(character.poison_rad_modifier || 0);
        document.getElementById('con-trauma').textContent = character.trauma_survival || '-';
        document.getElementById('con-test').textContent = character.constitution_test || '-';
        document.getElementById('con-feat').textContent = character.extra_constitution_feat || '-';
    }
    
    // Intelligence
    if (document.getElementById('ability-int')) {
        document.getElementById('ability-int').textContent = character.intelligence;
        document.getElementById('int-language').textContent = character.language_modifier || '-';
        document.getElementById('int-magician-bonus').textContent = character.magicians_bonus || '-';
        document.getElementById('int-magician-chance').textContent = character.magicians_chance || '-';
    }
    
    // Wisdom
    if (document.getElementById('ability-wis')) {
        document.getElementById('ability-wis').textContent = character.wisdom;
        document.getElementById('wis-willpower').textContent = formatModifier(character.willpower_modifier || 0);
        document.getElementById('wis-cleric-bonus').textContent = character.cleric_bonus || '-';
        document.getElementById('wis-cleric-chance').textContent = character.cleric_chance || '-';
    }
    
    // Charisma
    if (document.getElementById('ability-cha')) {
        document.getElementById('ability-cha').textContent = character.charisma;
        document.getElementById('cha-reaction').textContent = formatModifier(character.reaction_modifier || 0);
        document.getElementById('cha-followers').textContent = character.max_followers || '-';
        document.getElementById('cha-turning').textContent = formatModifier(character.undead_turning_modifier || 0);
    }
}

function displayClassAbilities(character) {
    const classAbilitiesSection = document.getElementById('class-abilities-section');
    const classAbilitiesTitle = document.getElementById('class-abilities-title');
    const classAbilitiesContainer = document.getElementById('class-abilities-container');

    if (!classAbilitiesSection || !classAbilitiesTitle || !classAbilitiesContainer) {
        console.warn('Class abilities DOM elements not found');
        return;
    }

    if (character.abilities && character.abilities.length > 0) {
        // Set the class name in the title
        classAbilitiesTitle.textContent = character.class;
        
        // Clear previous abilities
        classAbilitiesContainer.innerHTML = '';

        // Add each ability
        character.abilities.forEach(ability => {
            // Example implementation - adjust based on your actual ability data structure
            const abilityCard = document.createElement('div');
            abilityCard.className = 'ability-card';
            
            abilityCard.innerHTML = `
                <div class="ability-name">${ability.name}
                    ${ability.level ? `<span class="ability-level">Level ${ability.level}</span>` : ''}
                </div>
                <div class="ability-description">${ability.description}</div>
            `;
            
            classAbilitiesContainer.appendChild(abilityCard);
        });

        // Show the section
        classAbilitiesSection.style.display = 'block';
    } else {
        // Hide the section if no abilities
        classAbilitiesSection.style.display = 'none';
    }
}

// Helper functions
function formatModifier(value) {
    return value >= 0 ? `+${value}` : value.toString();
}

function toggleSpellsTab(characterClass) {
    const spellsTab = document.getElementById('spells-tab-nav');
    if (!spellsTab) return;
    
    // Only show spells tab for magic-using classes
    if (['Wizard', 'Cleric', 'Druid', 'Bard', 'Paladin'].includes(characterClass)) {
        spellsTab.style.display = 'block';
    } else {
        spellsTab.style.display = 'none';
    }
}