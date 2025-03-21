document.addEventListener('DOMContentLoaded', function() {
    if (document.getElementById('combat-tab')) {
        console.log('Combat tab initialized');
        const combatTab = document.querySelector('.tab-item[data-tab="combat-tab"]');
        if (combatTab) {
            combatTab.addEventListener('click', function() {
                console.log('Combat tab clicked');
                loadCombatData();
            });
        }
        if (document.getElementById('combat-tab').classList.contains('active')) {
            console.log('Combat tab is active by default');
            loadCombatData();
        }
    }
});

function loadCombatData() {
    const characterId = getCharacterIdFromURL();
    if (!characterId) {
        console.error('No character ID found in URL');
        return;
    }
    
    try {
        console.log('Character data at beginning of loadCombatData:', window.characterData);
        
        if (!window.characterData || !window.characterData.fighting_ability) {
            console.log('Fetching character details because fighting_ability is missing');
            fetchCharacterAndClassData(characterId);
        } else {
            console.log('Using existing character data with fighting ability:', window.characterData.fighting_ability);
            updateFightingAbility();
            fetchInventoryForCombat(characterId);
            generateCombatMatrix();
        }
    } catch (error) {
        console.error('Error loading combat data:', error);
    }
}

function getCharacterIdFromURL() {
    const pathParts = window.location.pathname.split('/');
    return pathParts[pathParts.length - 1];
}

async function fetchCharacterAndClassData(characterId) {
    try {
        const token = localStorage.getItem('authToken');
        
        // First fetch basic character data
        const charResponse = await fetch(`/api/characters/${characterId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!charResponse.ok) {
            throw new Error('Failed to fetch character details');
        }
        
        window.characterData = await charResponse.json();
        
        // Then fetch class-specific data
        const classResponse = await fetch(`/api/characters/${characterId}/class-data`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!classResponse.ok) {
            throw new Error('Failed to fetch class data');
        }
        
        const classData = await classResponse.json();
        console.log('Class data loaded:', classData);
        
        // Update character data with class-specific info if needed
        if (classData.current_level_data && classData.current_level_data.fighting_ability) {
            window.characterData.fighting_ability = classData.current_level_data.fighting_ability;
            console.log('Updated fighting ability:', window.characterData.fighting_ability);
        }
        
        // Now update the UI
        updateFightingAbility();
        
        // Important: Make sure we're still fetching the inventory
        fetchInventoryForCombat(characterId);
        
        generateCombatMatrix();
    } catch (error) {
        console.error('Error fetching character and class data:', error);
    }
}

function updateFightingAbility() {
    const faDisplay = document.getElementById('fighting-ability');
    if (faDisplay && window.characterData) {
        const fightingAbility = window.characterData.fighting_ability || 0;
        console.log('Updating fighting ability display to:', fightingAbility);
        faDisplay.textContent = fightingAbility;
    }
}

async function fetchInventoryForCombat(characterId) {
    try {
        const token = localStorage.getItem('authToken');
        const response = await fetch(`/api/inventories/character/${characterId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to fetch inventory');
        }
        
        const data = await response.json();
        console.log('Inventory data for combat:', data);
        processEquippedItems(data.items);
    } catch (error) {
        console.error('Error fetching inventory for combat:', error);
        // Important: Still call these functions even if fetching inventory fails
        // to display default values
        calculateArmorClass(null, null);
        calculateMovementRate(null);
        displayEquippedWeapons([]);
    }
}

function processEquippedItems(items) {
    if (!items || !Array.isArray(items)) {
        console.warn('No items data available');
        calculateArmorClass(null, null);
        calculateMovementRate(null);
        displayEquippedWeapons([]);
        return;
    }
    
    const equippedItems = items.filter(item => item.is_equipped === true);
    console.log('Equipped items:', equippedItems);
    
    const equippedArmor = equippedItems.find(item => 
        item.item_type === 'armor' && item.is_equipped === true
    );
    
    const equippedShield = equippedItems.find(item => 
        item.item_type === 'shield' && item.is_equipped === true
    );
    
    const equippedWeapons = equippedItems.filter(item => 
        item.item_type === 'weapon' && item.is_equipped === true
    );
    
    calculateArmorClass(equippedArmor, equippedShield);
    calculateMovementRate(equippedArmor);
    displayEquippedWeapons(equippedWeapons);
    
    // If no equipped items, still display default values
    if (equippedItems.length === 0) {
        console.log('No equipped items found - using defaults');
        calculateArmorClass(null, null);
        calculateMovementRate(null);
        displayEquippedWeapons([]);
    }
}

function calculateArmorClass(armor, shield) {
    let baseAC = 9; // Default unarmored AC
    let acSource = "Unarmored";
    
    if (armor && armor.item_details) {
        baseAC = armor.item_details.ac || 9;
        acSource = armor.item_details.name || "Armor";
    }
    
    let shieldMod = 0;
    if (shield && shield.item_details) {
        shieldMod = shield.item_details.defense_modifier || 0;
        acSource += ` + ${shield.item_details.name}`;
    }
    
    let defenceAdjustment = 0;
    if (window.characterData && window.characterData.defence_adjustment !== undefined) {
        defenceAdjustment = window.characterData.defence_adjustment;
    }
    
    const finalAC = baseAC - shieldMod - defenceAdjustment;
    
    const acDisplay = document.getElementById('armor-class');
    if (acDisplay) {
        acDisplay.textContent = finalAC;
        acDisplay.title = `${acSource} (Defence adjustment: ${defenceAdjustment >= 0 ? '+' : ''}${defenceAdjustment})`;
    }
}

function calculateMovementRate(armor) {
    let movementRate = 40; // Default unarmored movement rate
    
    if (armor && armor.item_details) {
        if (typeof armor.item_details.movement_rate === 'number') {
            movementRate = armor.item_details.movement_rate;
        }
    }
    
    const mvDisplay = document.getElementById('movement-rate');
    if (mvDisplay) {
        mvDisplay.textContent = movementRate.toString();
    }
}

function displayEquippedWeapons(weapons) {
    const weaponsTable = document.getElementById('weapons-table');
    const weaponsList = document.getElementById('weapons-list');
    const noWeapons = document.getElementById('no-weapons');
    const loading = document.getElementById('weapons-loading');
    
    if (loading) {
        loading.style.display = 'none';
    }
    
    if (!weapons || weapons.length === 0) {
        if (noWeapons) noWeapons.style.display = 'block';
        if (weaponsTable) weaponsTable.style.display = 'none';
        return;
    }
    
    if (weaponsList) {
        weaponsList.innerHTML = '';
    }
    
    // Get strength modifiers from global character data
    let meleeModifier = 0;
    let damageModifier = 0;
    let rangedModifier = 0;
    
    if (window.characterData) {
        meleeModifier = window.characterData.melee_modifier || 0;
        damageModifier = window.characterData.damage_adjustment || 0;
        rangedModifier = window.characterData.ranged_modifier || 0;
    }
    
    // Display each weapon
    weapons.forEach(weapon => {
        const details = weapon.item_details || {};
        const row = document.createElement('tr');
        
        const isRanged = details.category === 'Ranged' ||
                         (details.type && details.type.toLowerCase().includes('ranged'));
        
        const hitModifier = isRanged ? rangedModifier : meleeModifier;
        const toHitStr = formatModifier(hitModifier);
        const damageStr = `${details.damage || '1d6'} ${isRanged ? '' : formatModifier(damageModifier)}`;
        
        row.innerHTML = `
            <td>${details.name || 'Unknown Weapon'}</td>
            <td>${damageStr}</td>
            <td>${toHitStr}</td>
            <td>${details.special_properties || weapon.notes || ''}</td>
        `;
        
        weaponsList.appendChild(row);
    });
    
    if (weaponsTable) {
        weaponsTable.style.display = 'table';
    }
    
    if (noWeapons) {
        noWeapons.style.display = 'none';
    }
}

function generateCombatMatrix() {
    const matrixBody = document.getElementById('combat-matrix-body');
    if (!matrixBody) {
        console.warn('Combat matrix body element not found');
        return;
    }
    
    let fightingAbility = 0;
    if (window.characterData) {
        fightingAbility = window.characterData.fighting_ability || 0;
        const faDisplay = document.getElementById('fighting-ability');
        if (faDisplay) {
            faDisplay.textContent = fightingAbility;
        }
    }
    
    matrixBody.innerHTML = '';
    
    // Combat matrix values (rows are FA, columns are AC)
    const matrixData = [
        [11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29],
        [10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28],
        [9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27],
        [8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26],
        [7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25],
        [6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24],
        [5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23],
        [4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22],
        [3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21],
        [2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20],
        [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19],
        [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18],
        [-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17]
    ];
    
    // Generate the matrix rows
    for (let fa = 0; fa <= 12; fa++) {
        const row = document.createElement('tr');
        const rowData = matrixData[fa];
        
        const faCell = document.createElement('th');
        faCell.textContent = fa;
        row.appendChild(faCell);
        
        rowData.forEach((value, index) => {
            const cell = document.createElement('td');
            cell.textContent = value;
            
            if (fa === fightingAbility) {
                cell.classList.add('highlight');
            }
            
            row.appendChild(cell);
        });
        
        matrixBody.appendChild(row);
    }
}

function formatModifier(value) {
    return value >= 0 ? `+${value}` : value.toString();
}