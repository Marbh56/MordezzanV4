document.addEventListener('DOMContentLoaded', function() {
    const pathSegments = window.location.pathname.split('/');
    const characterId = pathSegments[pathSegments.indexOf('view') + 1];
    const masteryButton = document.querySelector('.tab-button[data-tab="weapon-mastery"]');
    
    if (masteryButton) {
        masteryButton.addEventListener('click', function() {
            loadWeaponMasteries();
        });
        
        if (masteryButton.classList.contains('active')) {
            loadWeaponMasteries();
        }
    }
    
    function loadWeaponMasteries() {
        const masteryTab = document.getElementById('weapon-mastery');
        
        masteryTab.innerHTML = `
            <div class="loading-container">
                <div class="loading-spinner"></div>
                <p>Loading weapon masteries...</p>
            </div>
        `;
        
        // Fetch available weapons and current masteries
        fetch(`/api/characters/${characterId}/weapon-masteries/available`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch weapon mastery data');
                }
                return response.json();
            })
            .then(data => renderMasteryUI(data))
            .catch(error => {
                console.error('Error loading weapon masteries:', error);
                masteryTab.innerHTML = `
                    <div class="error-message">
                        <h3>Error Loading Weapon Masteries</h3>
                        <p>${error.message}</p>
                    </div>
                `;
            });
    }
    
    function renderMasteryUI(data) {
        const masteryTab = document.getElementById('weapon-mastery');
        const { available_weapons, current_masteries, total_slots, used_slots, can_grand_master, character_level } = data;
        
        // Check if character class has weapon mastery
        if (total_slots === 0) {
            masteryTab.innerHTML = `
                <div class="empty-state">
                    <h3 class="empty-title">Weapon Mastery Not Available</h3>
                    <p class="empty-description">This character class does not have weapon mastery abilities.</p>
                </div>
            `;
            return;
        }
        
        let html = `
            <div class="mastery-header">
                <h2>Weapon Mastery</h2>
                <div class="mastery-slots">
                    <span class="slot-count">${used_slots}/${total_slots} Slots Used</span>
                    ${used_slots < total_slots ? `<button class="btn btn-primary add-mastery-btn" id="addMasteryBtn">Add Weapon Mastery</button>` : ''}
                </div>
            </div>
            <div class="mastery-description">
                <p>Weapon mastery provides +1 to hit and +1 damage with the chosen weapon. Grand mastery provides +2 to hit and +2 damage and improved attack rate.</p>
                <p>Characters gain additional mastery slots at levels 4, 8, and 12. At these levels, they may also choose to upgrade one mastered weapon to grand mastery instead of learning a new weapon.</p>
            </div>
        `;
        
        if (current_masteries.length === 0) {
            html += `
                <div class="empty-state">
                    <h3 class="empty-title">No Weapons Mastered</h3>
                    <p class="empty-description">This character hasn't mastered any weapons yet.</p>
                </div>
            `;
        } else {
            html += `<div class="masteries-grid">`;
            
            current_masteries.forEach(mastery => {
                const isMastered = mastery.mastery_level === 'mastered';
                const canUpgrade = isMastered && can_grand_master && character_level >= 4;
                
                html += `
                    <div class="mastery-card ${mastery.mastery_level}">
                        <div class="mastery-weapon-name">${mastery.weapon_base_name}</div>
                        <div class="mastery-level-badge">
                            ${mastery.mastery_level === 'grand_mastery' ? 'Grand Master' : 'Mastered'}
                        </div>
                        <div class="mastery-bonuses">
                            ${mastery.mastery_level === 'grand_mastery' ? 
                                '<div class="bonus">+2 to hit</div><div class="bonus">+2 damage</div><div class="bonus">Improved attack rate</div>' : 
                                '<div class="bonus">+1 to hit</div><div class="bonus">+1 damage</div>'}
                        </div>
                        <div class="mastery-actions">
                            ${canUpgrade ? 
                                `<button class="btn btn-secondary upgrade-btn" data-weapon-base-name="${mastery.weapon_base_name}">Upgrade to Grand Master</button>` : ''}
                            <button class="btn btn-danger remove-btn" data-weapon-base-name="${mastery.weapon_base_name}">Remove</button>
                        </div>
                    </div>
                `;
            });
            
            html += `</div>`;
        }
        
        masteryTab.innerHTML = html;
        
        // Add event listener for "Add Mastery" button
        const addBtn = document.getElementById('addMasteryBtn');
        if (addBtn) {
            addBtn.addEventListener('click', () => showAddMasteryModal(available_weapons));
        }
        
        // Add event listeners for upgrade buttons
        document.querySelectorAll('.upgrade-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const weaponBaseName = e.target.getAttribute('data-weapon-base-name');
                upgradeWeaponMastery(weaponBaseName);
            });
        });
        
        // Add event listeners for remove buttons
        document.querySelectorAll('.remove-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const weaponBaseName = e.target.getAttribute('data-weapon-base-name');
                removeWeaponMastery(weaponBaseName);
            });
        });
    }
    
    function showAddMasteryModal(availableWeapons) {
        // Create modal if it doesn't exist
        let modal = document.getElementById('masteryModal');
        if (!modal) {
            modal = document.createElement('div');
            modal.id = 'masteryModal';
            modal.className = 'modal';
            modal.innerHTML = `
                <div class="modal-content">
                    <h2 class="modal-title">Add Weapon Mastery</h2>
                    <form id="masteryForm">
                        <div class="modal-form-group">
                            <label for="weaponBaseName">Select Weapon</label>
                            <select id="weaponBaseName" name="weaponBaseName" required>
                                <option value="" disabled selected>Choose a weapon</option>
                            </select>
                        </div>
                        <div class="modal-actions">
                            <button type="button" class="btn btn-secondary" id="cancelMasteryBtn">Cancel</button>
                            <button type="submit" class="btn btn-primary">Add Mastery</button>
                        </div>
                    </form>
                </div>
            `;
            document.body.appendChild(modal);
            
            // Add event listeners
            document.getElementById('cancelMasteryBtn').addEventListener('click', () => {
                modal.style.display = 'none';
            });
            
            document.getElementById('masteryForm').addEventListener('submit', (e) => {
                e.preventDefault();
                const weaponBaseName = document.getElementById('weaponBaseName').value;
                addWeaponMastery(weaponBaseName);
            });
            
            // Close when clicking outside
            modal.addEventListener('click', (e) => {
                if (e.target === modal) {
                    modal.style.display = 'none';
                }
            });
        }
        
        // Populate select with available weapons
        const select = document.getElementById('weaponBaseName');
        select.innerHTML = '<option value="" disabled selected>Choose a weapon</option>';
        
        availableWeapons.forEach(weapon => {
            const option = document.createElement('option');
            option.value = weapon.name;
            option.textContent = weapon.name;
            select.appendChild(option);
        });
        
        // Show modal
        modal.style.display = 'flex';
    }
    
    function addWeaponMastery(weaponBaseName) {
        fetch(`/api/characters/${characterId}/weapon-masteries`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                weapon_base_name: weaponBaseName,
                mastery_level: 'mastered'
            })
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => Promise.reject(err));
            }
            return response.json();
        })
        .then(() => {
            // Hide modal
            document.getElementById('masteryModal').style.display = 'none';
            // Reload masteries
            loadWeaponMasteries();
        })
        .catch(error => {
            console.error('Error adding weapon mastery:', error);
            alert(`Failed to add weapon mastery: ${error.message || 'Unknown error'}`);
        });
    }
    
    function upgradeWeaponMastery(weaponBaseName) {
        if (!confirm('Are you sure you want to upgrade this weapon to Grand Mastery? You can only have one Grand Mastery weapon.')) {
            return;
        }
        
        fetch(`/api/characters/${characterId}/weapon-masteries/${encodeURIComponent(weaponBaseName)}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                mastery_level: 'grand_mastery'
            })
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => Promise.reject(err));
            }
            return response.json();
        })
        .then(() => {
            // Reload masteries
            loadWeaponMasteries();
        })
        .catch(error => {
            console.error('Error upgrading weapon mastery:', error);
            alert(`Failed to upgrade weapon mastery: ${error.message || 'Unknown error'}`);
        });
    }
    
    function removeWeaponMastery(weaponBaseName) {
        if (!confirm('Are you sure you want to remove this weapon mastery?')) {
            return;
        }
        
        fetch(`/api/characters/${characterId}/weapon-masteries/${encodeURIComponent(weaponBaseName)}`, {
            method: 'DELETE'
        })
        .then(response => {
            if (!response.ok) {
                const contentType = response.headers.get('content-type');
                if (contentType && contentType.includes('application/json')) {
                    return response.json().then(err => Promise.reject(err));
                } else {
                    return Promise.reject(new Error('Server error'));
                }
            }
            // Reload masteries
            loadWeaponMasteries();
        })
        .catch(error => {
            console.error('Error removing weapon mastery:', error);
            alert(`Failed to remove weapon mastery: ${error.message || 'Unknown error'}`);
        });
    }
});