document.addEventListener('DOMContentLoaded', function() {
    const pathSegments = window.location.pathname.split('/');
    const characterId = pathSegments[pathSegments.indexOf('view') + 1];
    console.log("Character ID:", characterId);
    
    // Make sure the combat tab's style is applied
    const style = document.createElement('style');
    style.textContent = `
        .combat-grid {
            display: flex;
            flex-direction: column;
            gap: 1.5rem;
            margin-bottom: 2rem;
        }
        
        .combat-stats-row {
            display: flex;
            justify-content: space-between;
            gap: 1.5rem;
        }
    `;
    document.head.appendChild(style);
    
    const combatButton = document.querySelector('.tab-button[data-tab="combat"]');
    if (combatButton) {
        combatButton.addEventListener('click', function() {
            console.log("Combat tab clicked");
            loadCombatEquipment();
            loadArmorClass();
        });
        
        if (combatButton.classList.contains('active')) {
            console.log("Combat tab is initially active");
            loadCombatEquipment();
            loadArmorClass();
        }
    }

    function loadCombatEquipment() {
        const armorContainer = document.getElementById('equipped-armor-container');
        const weaponsContainer = document.getElementById('equipped-weapons-container');
        
        if (armorContainer) {
            armorContainer.innerHTML = `
                <div class="loading-container">
                    <div class="loading-spinner"></div>
                    <p>Loading equipped armor...</p>
                </div>
            `;
        }
        
        if (weaponsContainer) {
            weaponsContainer.innerHTML = `
                <div class="loading-container">
                    <div class="loading-spinner"></div>
                    <p>Loading equipped weapons...</p>
                </div>
            `;
        }
        
        console.log("Fetching combat equipment for character ID:", characterId);
        fetch(`/api/characters/${characterId}/combat-equipment`, {
            headers: {
                'Accept': 'application/json'
            }
        })
        .then(response => {
            console.log("Combat equipment response status:", response.status);
            if (!response.ok) {
                return response.text().then(text => {
                    console.error("Response body:", text);
                    throw new Error(`Failed to fetch combat equipment: ${response.status}`);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("Combat equipment data:", data);
            if (armorContainer) {
                renderArmor(data.armor || [], armorContainer);
            }
            
            // Also load weapon stats in separate request
            loadWeaponStats(weaponsContainer);
            
            // Calculate movement rate based on armor
            updateMovementRate(data.armor || []);
        })
        .catch(error => {
            console.error("Error fetching combat equipment:", error);
            const errorMessage = `
                <div class="error-message">
                    <h3>Error Loading Equipment</h3>
                    <p>${error.message}</p>
                    <p>Check the browser console for more details.</p>
                </div>
            `;
            if (armorContainer) armorContainer.innerHTML = errorMessage;
        });
    }
    
    function loadWeaponStats(container) {
        if (!container) return;
        
        fetch(`/api/characters/${characterId}/weapon-stats`, {
            headers: {
                'Accept': 'application/json'
            }
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => {
                    console.error("Response body:", text);
                    throw new Error(`Failed to fetch weapon stats: ${response.status}`);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("Weapon stats data:", data);
            renderWeaponStats(data.weapon_stats || [], container);
        })
        .catch(error => {
            console.error("Error fetching weapon stats:", error);
            container.innerHTML = `
                <div class="error-message">
                    <h3>Error Loading Weapon Stats</h3>
                    <p>${error.message}</p>
                </div>
            `;
        });
    }
    
    function updateMovementRate(armorItems) {
        const movementRateElement = document.getElementById('movement-rate');
        if (!movementRateElement) return;
        
        // Default movement rate is 40 feet per round
        let movementRate = 40;
        
        // Check if any equipped armor affects movement
        if (armorItems.length > 0) {
            for (const item of armorItems) {
                if (item.inventory_item && item.inventory_item.is_equipped) {
                    const armorObj = item.armor;
                    
                    // If armor has movement_rate property and is less than current rate
                    if (armorObj && armorObj.movement_rate && armorObj.movement_rate < movementRate) {
                        movementRate = armorObj.movement_rate;
                    }
                    
                    // Check if armor is heavy - some heavy armors might not have the movement_rate property explicitly set
                    if (armorObj && armorObj.weight_class === "Heavy" && movementRate > 30) {
                        movementRate = 30;
                    }
                }
            }
        }
        
        // Update the movement rate display
        movementRateElement.textContent = movementRate;
    }

    function loadArmorClass() {
        const acContainer = document.getElementById('armor-class-container');
        if (!acContainer) return;
        
        acContainer.innerHTML = `
            <div class="loading-container">
                <div class="loading-spinner"></div>
                <p>Loading armor class details...</p>
            </div>
        `;
        
        console.log("Fetching character data for ID:", characterId);
        fetch(`/api/characters/${characterId}`, {
            headers: {
                'Accept': 'application/json'
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Failed to fetch character data: ${response.status}`);
            }
            return response.json();
        })
        .then(characterData => {
            console.log("Character data:", characterData);
            const fightingAbility = characterData.fighting_ability || 0;
            highlightCombatMatrix(fightingAbility);
            
            return fetch(`/api/characters/${characterId}/ac`, {
                headers: {
                    'Accept': 'application/json'
                }
            });
        })
        .then(response => {
            console.log("AC response status:", response.status);
            if (!response.ok) {
                return response.text().then(text => {
                    console.error("Response body:", text);
                    throw new Error(`Failed to fetch armor class: ${response.status}`);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("AC data:", data);
            renderArmorClass(data, acContainer);
        })
        .catch(error => {
            console.error("Error fetching armor class or character data:", error);
            acContainer.innerHTML = `
                <div class="error-message">
                    <h3>Error Loading Armor Class</h3>
                    <p>${error.message}</p>
                </div>
            `;
        });
    }
    
    function highlightCombatMatrix(fightingAbility) {
        const combatMatrix = document.querySelector('.combat-matrix');
        if (!combatMatrix) return;
        
        const highlightedRows = combatMatrix.querySelectorAll('tr.highlighted-row');
        highlightedRows.forEach(row => {
            row.classList.remove('highlighted-row');
        });
        
        const matrixNote = document.querySelector('.combat-matrix-note');
        if (matrixNote) {
            matrixNote.classList.remove('highlighted-row-exists');
        }
        
        const rows = combatMatrix.querySelectorAll('tbody tr');
        if (fightingAbility < 0 || fightingAbility > 12) {
            console.log("Fighting Ability outside matrix range:", fightingAbility);
            return;
        }
        
        if (rows[fightingAbility]) {
            rows[fightingAbility].classList.add('highlighted-row');
            if (matrixNote) {
                matrixNote.classList.add('highlighted-row-exists');
                let highlightedNote = matrixNote.querySelector('.highlighted-note');
                if (!highlightedNote) {
                    highlightedNote = document.createElement('div');
                    highlightedNote.className = 'highlighted-note';
                    matrixNote.appendChild(highlightedNote);
                }
                highlightedNote.textContent = `Your Fighting Ability (FA) is ${fightingAbility}, highlighted above.`;
            }
        }
    }
    
    function renderArmorClass(acData, container) {
        // Log the incoming data to help with debugging
        console.log("AC Data:", acData);
        
        let acHTML = `
            <div class="ac-breakdown">
                <div class="ac-total">
                    <div class="stat-value">${acData.final_ac}</div>
                </div>
                <div class="ac-components">
        `;
        
        if (acData.armor_equipped) {
            acHTML += `
                <div class="ac-component">
                    <span class="component-label">Armor (${acData.armor_equipped}):</span>
                    <span class="component-value">${acData.armor_ac}</span>
                </div>
            `;
        }
        
        if (acData.shield_bonus > 0) {
            acHTML += `
                <div class="ac-component">
                    <span class="component-label">Shield (${acData.shield_equipped}):</span>
                    <span class="component-value">-${acData.shield_bonus}</span>
                </div>
            `;
        }
        
        // Handle dexterity modifier - in Hyperborea, positive Dex modifier lowers AC (improves it)
        // and negative Dex modifier increases AC (worsens it)
        const dexMod = parseFloat(acData.dexterity_mod);
        if (!isNaN(dexMod)) {
            let dexDisplay;
            if (dexMod > 0) {
                // Positive dex mod improves AC (lowers it) so show with minus sign
                dexDisplay = `-${dexMod}`;
            } else if (dexMod < 0) {
                // Negative dex mod worsens AC (raises it) so show with plus sign
                dexDisplay = `+${Math.abs(dexMod)}`;
            } else {
                // Zero modifier has no effect
                dexDisplay = "+0";
            }
            
            acHTML += `
                <div class="ac-component">
                    <span class="component-label">Dexterity:</span>
                    <span class="component-value">${dexDisplay}</span>
                </div>
            `;
        }
        
        if (acData.natural_ac > 0) {
            acHTML += `
                <div class="ac-component">
                    <span class="component-label">Natural:</span>
                    <span class="component-value">-${acData.natural_ac}</span>
                </div>
            `;
        }
        
        if (acData.other_bonuses > 0) {
            acHTML += `
                <div class="ac-component">
                    <span class="component-label">Other:</span>
                    <span class="component-value">-${acData.other_bonuses}</span>
                </div>
            `;
        }
        
        acHTML += `
                </div>
            </div>
        `;
        
        container.innerHTML = acHTML;
    }
    
    function renderArmor(armorItems, container) {
        if (!container) return;
        
        if (!armorItems || armorItems.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <h3 class="empty-title">No Armor Equipped</h3>
                    <p class="empty-description">This character doesn't have any armor or shields equipped.</p>
                </div>
            `;
            return;
        }
        
        // Filter to only show equipped items
        const equippedArmorItems = armorItems.filter(item => 
            item.inventory_item && item.inventory_item.is_equipped
        );
        
        if (equippedArmorItems.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <h3 class="empty-title">No Armor Equipped</h3>
                    <p class="empty-description">This character has armor items in inventory but none are equipped.</p>
                </div>
            `;
            return;
        }
        
        let armorHTML = '<div class="armor-grid">';
        equippedArmorItems.forEach(item => {
            const inventoryItem = item.inventory_item;
            const isShield = inventoryItem.item_type === 'shield';
            const armorObj = isShield ? item.shield : item.armor;
            
            let mainStatLabel, mainStatValue;
            if (isShield) {
                mainStatLabel = "Defense";
                mainStatValue = `+${armorObj.defense_modifier}`;
            } else {
                mainStatLabel = "AC";
                mainStatValue = armorObj.AC || armorObj.ac;
            }
            
            armorHTML += `
                <div class="armor-card">
                    <h3 class="armor-name">${armorObj.name}</h3>
                    <div class="armor-type">${isShield ? 'Shield' : (armorObj.weight_class || 'Armor')}</div>
                    <div class="armor-stat">
                        <span class="stat-label">${mainStatLabel}</span>
                        <span class="stat-value">${mainStatValue}</span>
                    </div>
                    ${!isShield && armorObj.damage_reduction ? `
                    <div class="armor-stat">
                        <span class="stat-label">DR</span>
                        <span class="stat-value">${armorObj.damage_reduction}</span>
                    </div>` : ''}
                    <div class="armor-stat">
                        <span class="stat-label">Weight</span>
                        <span class="stat-value">${armorObj.weight}</span>
                    </div>
                    ${!isShield && armorObj.movement_rate ? `
                    <div class="armor-properties">
                        <span class="properties-label">Movement</span>
                        <span class="properties-value">${armorObj.movement_rate} ft</span>
                    </div>` : ''}
                </div>
            `;
        });
        armorHTML += '</div>';
        container.innerHTML = armorHTML;
    }
    
    function renderWeaponStats(weaponStats, container) {
        if (!container) return;
        
        if (!weaponStats || weaponStats.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <h3 class="empty-title">No Weapons Equipped</h3>
                    <p class="empty-description">This character doesn't have any weapons equipped.</p>
                </div>
            `;
            return;
        }
        
        // Filter to only show equipped weapons
        const equippedWeapons = weaponStats.filter(stats => 
            stats.inventory_item && stats.inventory_item.is_equipped
        );
        
        if (equippedWeapons.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <h3 class="empty-title">No Weapons Equipped</h3>
                    <p class="empty-description">This character has weapons in inventory but none are equipped.</p>
                </div>
            `;
            return;
        }
        
        let weaponsHTML = '<div class="weapons-grid">';
        equippedWeapons.forEach(stats => {
            const weapon = stats.weapon;
            const inventoryItem = stats.inventory_item;
            const isMastered = stats.is_mastered;
            const masteryClass = isMastered ?
                (stats.mastery_level === 'grand_mastery' ? 'grand-mastery' : 'mastered') : '';
            
            weaponsHTML += `
                <div class="weapon-card ${masteryClass}">
                    <h3 class="weapon-name">${weapon.name}</h3>
                    <div class="weapon-category">${weapon.category}</div>
                    ${isMastered ? `
                    <div class="mastery-badge ${stats.mastery_level}">
                        ${stats.mastery_level === 'grand_mastery' ? 'Grand Master' : 'Mastered'}
                    </div>` : ''}
                    <div class="weapon-stat">
                        <span class="stat-label">Damage</span>
                        <span class="stat-value">${stats.final_damage}</span>
                        ${stats.damage_bonus !== 0 ? `<span class="stat-bonus">(Base: ${stats.base_damage})</span>` : ''}
                    </div>
                    <div class="weapon-stat">
                        <span class="stat-label">To Hit</span>
                        <span class="stat-value">${stats.final_to_hit >= 0 ? '+' : ''}${stats.final_to_hit}</span>
                        ${stats.to_hit_bonus !== 0 ? `<span class="stat-bonus">(Base: ${stats.base_to_hit})</span>` : ''}
                    </div>
                    <!-- Determine if this is a ranged/hurled weapon -->
                    ${(() => {
                        const isMissileWeapon = weapon.category === 'Ranged' || weapon.category === 'Hurled';
                        const attackRateLabel = isMissileWeapon ? "Rate of Fire" : "Attack Rate";
                        return `
                        <div class="weapon-stat">
                            <span class="stat-label">${attackRateLabel}</span>
                            <span class="stat-value">${stats.final_attack_rate}</span>
                            ${stats.improved_attack_rate ? `<span class="stat-bonus">(Base: ${stats.base_attack_rate})</span>` : ''}
                        </div>`;
                    })()}
                    ${(weapon.range_short && weapon.range_medium && weapon.range_long) ? `
                    <div class="weapon-stat">
                        <span class="stat-label">Range</span>
                        <span class="stat-value">${weapon.range_short}/${weapon.range_medium}/${weapon.range_long}</span>
                    </div>` : ''}
                    ${weapon.properties ? `
                    <div class="weapon-properties">
                        <span class="properties-label">Properties</span>
                        <span class="properties-value">${weapon.properties}</span>
                    </div>` : ''}
                    ${inventoryItem.quantity > 1 ? `
                    <div class="weapon-properties">
                        <span class="properties-label">Quantity</span>
                        <span class="properties-value">${inventoryItem.quantity}</span>
                    </div>` : ''}
                </div>
            `;
        });
        weaponsHTML += '</div>';
        container.innerHTML = weaponsHTML;
        
        // Add CSS for mastery styling if not already present
        if (!document.getElementById('weapon-stats-style')) {
            const style = document.createElement('style');
            style.id = 'weapon-stats-style';
            style.textContent = `
                .weapon-card.mastered {
                    border-left: 3px solid var(--primary-color);
                }
                .weapon-card.grand-mastery {
                    border-left: 3px solid #FFD700;
                    background-color: rgba(255, 215, 0, 0.05);
                }
                .mastery-badge {
                    display: inline-block;
                    margin: 0.5rem 0;
                    padding: 0.2rem 0.5rem;
                    border-radius: 4px;
                    font-size: 0.8rem;
                    background-color: var(--primary-color);
                    color: #000;
                }
                .mastery-badge.grand_mastery {
                    background-color: #FFD700;
                    color: #000;
                }
                .stat-bonus {
                    display: block;
                    font-size: 0.8rem;
                    color: #777;
                    margin-top: 0.2rem;
                }
            `;
            document.head.appendChild(style);
        }
    }
});