document.addEventListener('DOMContentLoaded', function() {
    const pathSegments = window.location.pathname.split('/');
    const characterId = pathSegments[pathSegments.indexOf('view') + 1];
    
    // Get references to the tab and the weapon container
    const combatTab = document.getElementById('combat');
    const weaponsContainer = document.getElementById('equipped-weapons-container');
    
    if (combatTab && weaponsContainer) {
        // Add a listener for the combat tab to load weapon stats
        const combatButton = document.querySelector('.tab-button[data-tab="combat"]');
        if (combatButton) {
            combatButton.addEventListener('click', function() {
                loadWeaponStats();
            });
            
            // If combat tab is initially active, load the stats
            if (combatButton.classList.contains('active')) {
                loadWeaponStats();
            }
        }
    }
    
    function loadWeaponStats() {
        if (!weaponsContainer) return;
        
        weaponsContainer.innerHTML = `
            <div class="loading-container">
                <div class="loading-spinner"></div>
                <p>Loading weapon stats...</p>
            </div>
        `;
        
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
            renderWeaponStats(data.weapon_stats || []);
        })
        .catch(error => {
            console.error("Error fetching weapon stats:", error);
            weaponsContainer.innerHTML = `
                <div class="error-message">
                    <h3>Error Loading Weapon Stats</h3>
                    <p>${error.message}</p>
                    <p>Check the browser console for more details.</p>
                </div>
            `;
        });
    }
    
    function renderWeaponStats(weaponStats) {
        if (!weaponStats || weaponStats.length === 0) {
            weaponsContainer.innerHTML = `
                <div class="empty-state">
                    <h3 class="empty-title">No Weapons Equipped</h3>
                    <p class="empty-description">This character doesn't have any weapons equipped.</p>
                </div>
            `;
            return;
        }
        
        let weaponsHTML = '<div class="weapons-grid">';
        
        weaponStats.forEach(stats => {
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
        weaponsContainer.innerHTML = weaponsHTML;
        
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