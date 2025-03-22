let currentHP = 0;
let currentXP = 0;
let currentClass = '';
let fighterLevels = [];
let currentLevel = 1;
let characterId = '';

document.addEventListener('DOMContentLoaded', function () {
    const token = localStorage.getItem('authToken');
    const userId = localStorage.getItem('userId');
    const username = localStorage.getItem('username');
    if (!token || !userId) {
        window.location.href = '/auth/login-page';
        return;
    }
    document.getElementById('username-display').textContent = username;
    document.getElementById('logout-link').addEventListener('click', function (e) {
        e.preventDefault();
        localStorage.removeItem('authToken');
        localStorage.removeItem('userId');
        localStorage.removeItem('username');
        window.location.href = '/auth/login-page';
    });
    characterId = getCharacterIdFromURL();
    if (characterId) {
        fetchCharacterDetails(characterId).then(() => {
            if (currentClass) {
                fetchClassData();
            }
        }).catch(error => {
            console.error("Error in character loading sequence:", error);
        });
        const editButton = document.getElementById('edit-character-btn');
        if (editButton) {
            editButton.href = `/characters/${characterId}/edit`;
        }
    } else {
        window.location.href = '/';
    }
    setupHPControls();
    setupXPControls();
    setupTabs();
});

function getCharacterIdFromURL() {
    const pathParts = window.location.pathname.split('/');
    return pathParts[pathParts.length - 1];
}

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
        toggleSpellsTab(character.class);
    } catch (error) {
        console.error('Error fetching character details:', error);
        alert('Failed to load character details: ' + error.message);
    }
}

function displayCharacterDetails(character) {
    document.title = `${character.name} - Mordezzan`;
    document.getElementById('character-name').textContent = character.name;
    document.getElementById('character-level-class').textContent = `Level ${character.level} ${character.class}`;
    document.getElementById('info-name').textContent = character.name;
    document.getElementById('info-class').textContent = character.class;
    document.getElementById('info-level').textContent = character.level;
    document.getElementById('info-hit-dice').textContent = character.hit_dice || '-';
    currentHP = character.hit_points;
    currentXP = character.experience_points || 0;
    currentClass = character.class;
    currentLevel = character.level;
    updateHPDisplay(character.hit_points);
    document.getElementById('hp-set').value = character.hit_points;
    updateXPDisplay(character.experience_points || 0, character.level, character.class);
    const savingThrowsSection = document.getElementById('saving-throws-section');
    if (character.saving_throw) {
        updateSavingThrows(character);
        savingThrowsSection.style.display = 'block';
    } else {
        savingThrowsSection.style.display = 'none';
    }
    updateAbilityScores(character);
    displayClassAbilities(character);
    setupNotesTab();
    if (character.class === 'Fighter') {
        fetchFighterLevels();
    }
}

function updateHPDisplay(hp) {
    document.getElementById('info-hp').textContent = hp;
    document.getElementById('hp-set').value = hp;
    const hpElement = document.getElementById('info-hp');
    if (hp <= 0) {
        hpElement.style.color = '#e74c3c';
    } else if (hp <= 5) {
        hpElement.style.color = '#e67e22';
    } else {
        hpElement.style.color = ''; // Default color
    }
}

// Update XP display in the UI
function updateXPDisplay(xp, level, characterClass) {
    document.getElementById('info-experience').textContent = xp;
    document.getElementById('xp-set').value = xp;
    const fighterLevelInfo = document.getElementById('fighter-level-info');
    if (characterClass === 'Fighter' && fighterLevelInfo) {
        fighterLevelInfo.style.display = 'block';
        document.getElementById('current-level').textContent = level;
        let nextLevelXP = 0;
        let xpNeeded = 0;
        if (fighterLevels && fighterLevels.length > 0) {
            for (let i = 0; i < fighterLevels.length; i++) {
                if (fighterLevels[i].level > level) {
                    nextLevelXP = fighterLevels[i].experience_points;
                    break;
                }
            }
            xpNeeded = Math.max(0, nextLevelXP - xp);
            if (xpNeeded > 0) {
                document.getElementById('next-level-xp').textContent = nextLevelXP;
                document.getElementById('xp-needed').textContent = xpNeeded;
            } else {
                document.getElementById('next-level-xp').innerHTML = '<i>Maximum level reached</i>';
                document.getElementById('xp-needed').textContent = '0';
            }
        }
    } else if (fighterLevelInfo) {
        fighterLevelInfo.style.display = 'none';
    }
}

async function fetchFighterLevels() {
    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`/api/characters/${characterId}/class-data`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (!response.ok) {
            throw new Error('Failed to fetch fighter levels');
        }
        const data = await response.json();
        if (data.class_type === 'Fighter' && data.level_data) {
            fighterLevels = data.level_data;
            updateXPDisplay(currentXP, currentLevel, currentClass);
        }
    } catch (error) {
        console.error('Error fetching fighter levels:', error);
    }
}

async function fetchClassData() {
    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`/api/characters/${characterId}/class-data`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (!response.ok) {
            throw new Error('Failed to fetch class data');
        }
        const data = await response.json();
        console.log("Class data loaded:", data);
        if (data.class_type === 'Fighter') {
            fighterLevels = data.level_data;
        } else if (data.class_type === 'Magician') {
            window.magicianLevels = data.level_data;
        }
        updateClassSpecificUI(data);
        updateXPDisplay(currentXP, currentLevel, currentClass);
        return data;
    } catch (error) {
        console.error('Error fetching class data:', error);
    }
}

function updateClassSpecificUI(classData) {
    if (classData.current_level_data) {
        const levelData = classData.current_level_data;
        const hitDiceElement = document.getElementById('info-hit-dice');
        if (hitDiceElement && levelData.hit_dice) {
            hitDiceElement.textContent = levelData.hit_dice;
        }
        const baseSaveElement = document.getElementById('base_save');
        if (baseSaveElement && levelData.saving_throw) {
            baseSaveElement.textContent = levelData.saving_throw;
        }
        if (classData.class_type === 'Magician' && levelData.spell_slots) {
            updateSpellSlots(levelData.spell_slots);
        }
        if (levelData.abilities) {
            displayClassAbilities({
                class: classData.class_type,
                abilities: levelData.abilities
            });
        }
    }
}

function updateSpellSlots(spellSlots) {
    const spellSlotsContainer = document.getElementById('spell-slots-container');
    if (!spellSlotsContainer) return;
    spellSlotsContainer.innerHTML = '';
    // Add a card for each spell level
    for (const [level, count] of Object.entries(spellSlots)) {
        if (count <= 0) continue;
        const levelNumber = level.replace('level', '');
        const slotCard = document.createElement('div');
        slotCard.className = 'spell-slot-card';
        slotCard.innerHTML = `
            <div class="spell-slot-level">Level ${levelNumber}</div>
            <div class="spell-slot-count">${count}</div>
            <div class="spell-slot-used">Available</div>
        `;
        spellSlotsContainer.appendChild(slotCard);
    }
    spellSlotsContainer.style.display = 'flex';
}

function updateSavingThrows(character) {
    console.log("Updating saving throws for character:", character.name);
    const baseSave = character.saving_throw || 15;
    document.getElementById('base_save').textContent = baseSave;
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
    document.getElementById('ability-str').textContent = character.strength;
    document.getElementById('str-attack-mod').textContent = formatModifier(character.melee_modifier || 0);
    document.getElementById('str-damage-adj').textContent = formatModifier(character.damage_adjustment || 0);
    document.getElementById('str-test').textContent = character.strength_test || '-';
    document.getElementById('str-feat').textContent = character.extra_strength_feat || '-';
    document.getElementById('ability-dex').textContent = character.dexterity;
    document.getElementById('dex-ranged-mod').textContent = formatModifier(character.ranged_modifier || 0);
    document.getElementById('dex-defence-adj').textContent = formatModifier(character.defence_adjustment || 0);
    document.getElementById('dex-test').textContent = character.dexterity_test || '-';
    document.getElementById('dex-feat').textContent = character.extra_dexterity_feat || '-';
    document.getElementById('ability-con').textContent = character.constitution;
    document.getElementById('con-hp-mod').textContent = formatModifier(character.hp_modifier || 0);
    document.getElementById('con-poison-mod').textContent = formatModifier(character.poison_rad_modifier || 0);
    document.getElementById('con-trauma').textContent = character.trauma_survival || '-';
    document.getElementById('con-test').textContent = character.constitution_test || '-';
    document.getElementById('ability-int').textContent = character.intelligence;
    document.getElementById('int-language').textContent = character.language_modifier || '-';
    document.getElementById('int-magician-bonus').textContent = character.magicians_bonus || '-';
    document.getElementById('int-magician-chance').textContent = character.magicians_chance || '-';
    document.getElementById('ability-wis').textContent = character.wisdom;
    document.getElementById('wis-willpower').textContent = formatModifier(character.willpower_modifier || 0);
    document.getElementById('wis-cleric-bonus').textContent = character.cleric_bonus || '-';
    document.getElementById('wis-cleric-chance').textContent = character.cleric_chance || '-';
    document.getElementById('ability-cha').textContent = character.charisma;
    document.getElementById('cha-reaction').textContent = formatModifier(character.reaction_modifier || 0);
    document.getElementById('cha-followers').textContent = character.max_followers || '-';
    document.getElementById('cha-turning').textContent = formatModifier(character.undead_turning_modifier || 0);
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
        classAbilitiesTitle.textContent = character.class;
        classAbilitiesContainer.innerHTML = '';
        // Add each ability
        character.abilities.forEach(ability => {
            const abilityCard = document.createElement('div');
            abilityCard.className = 'ability-card';
            abilityCard.innerHTML = `
                <div class="ability-name">${ability.name}
                    ${ability.min_level ? `<span class="ability-level">Level ${ability.min_level}</span>` : ''}
                </div>
                <div class="ability-description">${ability.description}</div>
            `;
            classAbilitiesContainer.appendChild(abilityCard);
        });
        classAbilitiesSection.style.display = 'block';
    } else {
        classAbilitiesSection.style.display = 'none';
    }
}

function formatModifier(value) {
    return value >= 0 ? `+${value}` : value.toString();
}

function setupHPControls() {
    const hpControl = document.getElementById('hp-control-panel');
    const toggleButton = document.getElementById('toggle-hp-panel');
    const hpHeader = document.querySelector('.hp-header');
    function toggleHPPanel() {
        hpControl.classList.toggle('collapsed');
    }
    toggleButton.addEventListener('click', function (e) {
        e.stopPropagation();
        toggleHPPanel();
    });
    hpHeader.addEventListener('click', toggleHPPanel);
    document.getElementById('btn-damage').addEventListener('click', function (e) {
        e.stopPropagation();
        const damageAmount = parseInt(document.getElementById('hp-change').value) || 0;
        if (damageAmount <= 0) return;
        const newHP = Math.max(0, currentHP - damageAmount);
        updateCharacterHP(newHP);
    });
    document.getElementById('btn-heal').addEventListener('click', function (e) {
        e.stopPropagation();
        const healAmount = parseInt(document.getElementById('hp-change').value) || 0;
        if (healAmount <= 0) return;
        const newHP = currentHP + healAmount;
        updateCharacterHP(newHP);
    });
    document.getElementById('btn-set-hp').addEventListener('click', function (e) {
        e.stopPropagation();
        const newHP = parseInt(document.getElementById('hp-set').value);
        if (isNaN(newHP) || newHP < 0) return;
        updateCharacterHP(newHP);
    });
}

async function updateCharacterHP(newHP) {
    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`/api/characters/${characterId}/hp`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ hit_points: newHP })
        });
        if (!response.ok) {
            throw new Error('Failed to update HP');
        }
        currentHP = newHP;
        updateHPDisplay(newHP);
    } catch (error) {
        console.error('Error updating character HP:', error);
        alert('Failed to update HP: ' + error.message);
    }
}

function setupXPControls() {
    const xpControl = document.getElementById('xp-control-panel');
    const toggleButton = document.getElementById('toggle-xp-panel');
    const xpHeader = document.querySelector('.xp-header');
    function toggleXPPanel() {
        xpControl.classList.toggle('collapsed');
    }
    toggleButton.addEventListener('click', function (e) {
        e.stopPropagation();
        toggleXPPanel();
    });
    xpHeader.addEventListener('click', toggleXPPanel);
    document.getElementById('btn-reward').addEventListener('click', function (e) {
        e.stopPropagation();
        const xpAmount = parseInt(document.getElementById('xp-change').value) || 0;
        if (xpAmount <= 0) return;
        const newXP = currentXP + xpAmount;
        updateCharacterXP(newXP);
    });
    document.getElementById('btn-set-xp').addEventListener('click', function (e) {
        e.stopPropagation();
        const newXP = parseInt(document.getElementById('xp-set').value);
        if (isNaN(newXP) || newXP < 0) return;
        updateCharacterXP(newXP);
    });
}

async function updateCharacterXP(newXP) {
    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`/api/characters/${characterId}/xp`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ experience_points: newXP })
        });
        if (!response.ok) {
            throw new Error('Failed to update XP');
        }
        const data = await response.json();
        currentXP = newXP;
        if (data.level) {
            currentLevel = data.level;
        }
        updateXPDisplay(newXP, currentLevel, currentClass);
    } catch (error) {
        console.error('Error updating character XP:', error);
        alert('Failed to update XP: ' + error.message);
    }
}

function setupTabs() {
    const tabItems = document.querySelectorAll('.tab-item');
    const tabContents = document.querySelectorAll('.tab-content');
    tabItems.forEach(tab => {
        tab.addEventListener('click', () => {
            tabItems.forEach(item => item.classList.remove('active'));
            tabContents.forEach(content => content.classList.remove('active'));
            tab.classList.add('active');
            const tabId = tab.getAttribute('data-tab');
            document.getElementById(tabId).classList.add('active');
            if (tabId === 'spells-tab') {
                if (typeof fetchSpells === 'function') {
                    fetchSpells();
                }
            } else if (tabId === 'inventory-tab') {
                if (typeof fetchInventory === 'function') {
                    fetchInventory();
                }
            }
        });
    });
}

function toggleSpellsTab(characterClass) {
    const spellsTab = document.getElementById('spells-tab-nav');
    if (!spellsTab) return;
    if (['Wizard', 'Magician', 'Cleric', 'Druid', 'Bard', 'Paladin'].includes(characterClass)) {
        spellsTab.style.display = 'block';
    } else {
        spellsTab.style.display = 'none';
    }
}

function setupNotesTab() {
    console.log("Notes tab setup");
}