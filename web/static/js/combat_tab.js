document.addEventListener('DOMContentLoaded', function() {
    const pathSegments = window.location.pathname.split('/');
    const characterId = pathSegments[pathSegments.indexOf('view') + 1];
    console.log("Character ID:", characterId);

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
        const weaponsContainer = document.getElementById('equipped-weapons-container');
        const armorContainer = document.getElementById('equipped-armor-container');
        
        if (weaponsContainer) {
            weaponsContainer.innerHTML = `
                <div class="loading-container">
                    <div class="loading-spinner"></div>
                    <p>Loading equipped weapons...</p>
                </div>
            `;
        }
        
        if (armorContainer) {
            armorContainer.innerHTML = `
                <div class="loading-container">
                    <div class="loading-spinner"></div>
                    <p>Loading equipped armor...</p>
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
                renderWeapons(data.weapons || [], weaponsContainer);
                renderArmor(data.armor || [], armorContainer);
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
                if (weaponsContainer) weaponsContainer.innerHTML = errorMessage;
                if (armorContainer) armorContainer.innerHTML = errorMessage;
            });
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
        
        console.log("Fetching AC for character ID:", characterId);
        fetch(`/api/characters/${characterId}/ac`, {
            headers: {
                'Accept': 'application/json'
            }
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
                console.error("Error fetching armor class:", error);
                acContainer.innerHTML = `
                    <div class="error-message">
                        <h3>Error Loading Armor Class</h3>
                        <p>${error.message}</p>
                    </div>
                `;
            });
    }
    
    function renderArmorClass(acData, container) {
        // Create a detailed breakdown of AC calculation
        let acHTML = `
            <div class="ac-breakdown">
                <div class="ac-total">
                    <h3>Armor Class</h3>
                    <div class="stat-value">${acData.final_ac}</div>
                </div>
                <div class="ac-components">
                    <div class="ac-component-title">AC Calculation</div>
                    <div class="ac-component">
                        <span class="component-label">Base AC:</span>
                        <span class="component-value">${acData.base_ac}</span>
                    </div>
        `;
        
        // Only show armor if equipped
        if (acData.armor_equipped) {
            acHTML += `
                <div class="ac-component">
                    <span class="component-label">Armor (${acData.armor_equipped}):</span>
                    <span class="component-value">${acData.armor_ac}</span>
                </div>
            `;
        }
        
        // Add shield bonus if any
        if (acData.shield_bonus > 0) {
            acHTML += `
                <div class="ac-component">
                    <span class="component-label">Shield (${acData.shield_equipped}):</span>
                    <span class="component-value">-${acData.shield_bonus}</span>
                </div>
            `;
        }
        
        // Add dexterity modifier if any
        if (acData.dexterity_mod !== 0) {
            acHTML += `
                <div class="ac-component">
                    <span class="component-label">Dexterity Modifier:</span>
                    <span class="component-value">${acData.dexterity_mod > 0 ? '-' + acData.dexterity_mod : '+' + Math.abs(acData.dexterity_mod)}</span>
                </div>
            `;
        }
        
        // Add natural armor if any
        if (acData.natural_ac > 0) {
            acHTML += `
                <div class="ac-component">
                    <span class="component-label">Natural Armor:</span>
                    <span class="component-value">-${acData.natural_ac}</span>
                </div>
            `;
        }
        
        // Add other bonuses if any
        if (acData.other_bonuses > 0) {
            acHTML += `
                <div class="ac-component">
                    <span class="component-label">Other Bonuses:</span>
                    <span class="component-value">-${acData.other_bonuses}</span>
                </div>
            `;
        }
        
        // Complete the HTML
        acHTML += `
                    <div class="ac-component ac-final">
                        <span class="component-label">Final AC:</span>
                        <span class="component-value">${acData.final_ac}</span>
                    </div>
                </div>
            </div>
        `;
        
        container.innerHTML = acHTML;
    }

    function renderWeapons(weapons, container) {
        if (!container) return;
        if (!weapons || weapons.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <h3 class="empty-title">No Weapons Equipped</h3>
                    <p class="empty-description">This character doesn't have any weapons equipped.</p>
                </div>
            `;
            return;
        }

        let weaponsHTML = '<div class="weapons-grid">';
        weapons.forEach(item => {
            const weapon = item.weapon;
            const inventoryItem = item.inventory_item;
            weaponsHTML += `
                <div class="weapon-card">
                    <h3 class="weapon-name">${weapon.name}</h3>
                    <div class="weapon-category">${weapon.category}</div>
                    <div class="weapon-stat">
                        <span class="stat-label">Damage</span>
                        <span class="stat-value">${weapon.damage}</span>
                    </div>
                    ${weapon.weapon_class ? `
                    <div class="weapon-stat">
                        <span class="stat-label">Class</span>
                        <span class="stat-value">${weapon.weapon_class}</span>
                    </div>` : ''}
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

        let armorHTML = '<div class="armor-grid">';
        armorItems.forEach(item => {
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
});