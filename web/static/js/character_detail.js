document.addEventListener('DOMContentLoaded', function() {
    // Classes that have weapon mastery abilities
    const weaponMasteryClasses = [
        'Fighter', 'Barbarian', 'Berserker', 'Cataphract', 
        'Huntsman', 'Paladin', 'Ranger', 'Assassin'
    ];
    
    // Get the character class from the meta info
    const characterClassElement = document.querySelector('.character-meta');
    if (characterClassElement) {
        const classText = characterClassElement.textContent;
        const characterClass = classText.match(/Level \d+ (.+)/)?.[1];
        
        // Check if the class has weapon mastery
        const hasWeaponMastery = weaponMasteryClasses.includes(characterClass);
        
        // Hide the tab if the class doesn't have weapon mastery
        const masteryButton = document.querySelector('.tab-button[data-tab="weapon-mastery"]');
        if (masteryButton && !hasWeaponMastery) {
            masteryButton.style.display = 'none';
        }
    }

    const tabButtons = document.querySelectorAll('.tab-button');
    const tabContents = document.querySelectorAll('.tab-content');
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabId = button.getAttribute('data-tab');
            tabButtons.forEach(btn => btn.classList.remove('active'));
            tabContents.forEach(content => content.classList.remove('active'));
            button.classList.add('active');
            document.getElementById(tabId).classList.add('active');
        });
    });
    const pathSegments = window.location.pathname.split('/');
    const characterId = pathSegments[pathSegments.indexOf('view') + 1];
    const hpModal = document.getElementById('hpModal');
    const takeDamageBtn = document.getElementById('takeDamageBtn');
    const healBtn = document.getElementById('healBtn');
    const cancelHpBtn = document.getElementById('cancelHpBtn');
    const hpForm = document.getElementById('hpForm');
    const hpModalTitle = document.getElementById('hpModalTitle');
    let isHealing = false;
    function openHpModal(healing) {
        isHealing = healing;
        hpModalTitle.textContent = healing ? 'Heal Character' : 'Apply Damage';
        document.getElementById('tempHP').parentElement.style.display = healing ? 'block' : 'none';
        document.getElementById('confirmHpBtn').className = healing ? 'btn btn-primary' : 'btn damage-btn';
        hpModal.style.display = 'flex';
    }
    if (takeDamageBtn) {
        takeDamageBtn.addEventListener('click', () => openHpModal(false));
    }
    if (healBtn) {
        healBtn.addEventListener('click', () => openHpModal(true));
    }
    if (cancelHpBtn) {
        cancelHpBtn.addEventListener('click', () => {
            hpModal.style.display = 'none';
        });
    }
    if (hpForm) {
        hpForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const amount = parseInt(document.getElementById('hpAmount').value);
            const temp = document.getElementById('tempHP').checked;
            if (amount <= 0) {
                alert('Please enter a positive number');
                return;
            }
            try {
                const response = await fetch(`/api/characters/${characterId}/modify-hp`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        delta: isHealing ? amount : -amount,
                        temp: temp
                    })
                });
                if (!response.ok) {
                    throw new Error('Failed to update hit points');
                }
                const data = await response.json();
                const hpStatValue = document.querySelector('.hp-stat .stat-value');
                hpStatValue.textContent = `${data.current_hit_points}/${data.max_hit_points}`;
                const tempHpElement = document.querySelector('.temp-hp');
                if (data.temporary_hit_points > 0) {
                    if (tempHpElement) {
                        tempHpElement.textContent = `+${data.temporary_hit_points} temp`;
                    } else {
                        const newTempHp = document.createElement('div');
                        newTempHp.className = 'temp-hp';
                        newTempHp.textContent = `+${data.temporary_hit_points} temp`;
                        hpStatValue.after(newTempHp);
                    }
                } else if (tempHpElement) {
                    tempHpElement.remove();
                }
                hpModal.style.display = 'none';
            } catch (error) {
                console.error('Error:', error);
                alert('Failed to update hit points: ' + error.message);
            }
        });
    }
    const xpModal = document.getElementById('xpModal');
    const addXpBtn = document.getElementById('addXpBtn');
    const cancelXpBtn = document.getElementById('cancelXpBtn');
    const xpForm = document.getElementById('xpForm');
    if (addXpBtn) {
        addXpBtn.addEventListener('click', () => {
            xpModal.style.display = 'flex';
        });
    }
    if (cancelXpBtn) {
        cancelXpBtn.addEventListener('click', () => {
            xpModal.style.display = 'none';
        });
    }
    if (xpForm) {
        xpForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const amount = parseInt(document.getElementById('xpAmount').value);
            if (amount <= 0) {
                alert('Please enter a positive number');
                return;
            }
            try {
                const xpElement = document.querySelector('.xp-stat .stat-value');
                const currentXp = parseInt(xpElement.textContent);
                const newXp = currentXp + amount;
                const response = await fetch(`/api/characters/${characterId}/xp`, {
                    method: 'PATCH',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        experience_points: newXp
                    })
                });
                if (!response.ok) {
                    throw new Error('Failed to update experience points');
                }
                const data = await response.json();
                xpElement.textContent = data.experience_points;
                const levelElement = document.querySelector('.character-meta');
                const currentLevel = parseInt(levelElement.textContent.match(/Level (\d+)/)[1]);
                if (data.level > currentLevel) {
                    alert(`Congratulations! You've reached level ${data.level}!`);
                    levelElement.textContent = levelElement.textContent.replace(
                        `Level ${currentLevel}`,
                        `Level ${data.level}`
                    );
                    location.reload();
                }
                const xpNeededElement = document.querySelector('.xp-stat div:not(.stat-value):not(.stat-actions)');
                if (xpNeededElement && data.next_level_experience) {
                    const xpNeeded = data.next_level_experience - data.experience_points;
                    xpNeededElement.textContent = `${xpNeeded} to next level`;
                }
                xpModal.style.display = 'none';
            } catch (error) {
                console.error('Error:', error);
                alert('Failed to update experience points: ' + error.message);
            }
        });
    }
    window.addEventListener('click', (e) => {
        if (e.target === hpModal) {
            hpModal.style.display = 'none';
        }
        if (e.target === xpModal) {
            xpModal.style.display = 'none';
        }
    });
    const fetchAbilities = async () => {
        const abilitiesTab = document.getElementById('abilities');
        try {
            const response = await fetch(`/api/characters/${characterId}/class-data`);
            if (!response.ok) {
                throw new Error('Failed to fetch character abilities');
            }
            const data = await response.json();
            renderAbilities(data);
        } catch (error) {
            console.error('Error:', error);
            abilitiesTab.innerHTML = `
                <div class="error-message">
                    Failed to load character abilities: ${error.message}
                </div>
            `;
        }
    };
    const renderAbilities = (classData) => {
        const abilitiesTab = document.getElementById('abilities');
        if (!classData || !classData.current_level_data || !classData.current_level_data.abilities ||
            classData.current_level_data.abilities.length === 0) {
            abilitiesTab.innerHTML = `
                <div class="empty-state">
                    <h3 class="empty-title">No Abilities Found</h3>
                    <p class="empty-description">This character class doesn't have any special abilities at the current level.</p>
                </div>
            `;
            return;
        }
        const abilities = classData.current_level_data.abilities;
        let abilitiesHTML = `
            <h2>Class Abilities</h2>
            <div class="abilities-container">
        `;
        abilities.forEach(ability => {
            abilitiesHTML += `
                <div class="ability-card">
                    <h3 class="ability-name">${ability.name}</h3>
                    <div class="ability-description">${ability.description}</div>
                    <div class="ability-level">Available from Level ${ability.min_level}</div>
                </div>
            `;
        });
        abilitiesHTML += `</div>`;
        abilitiesTab.innerHTML = abilitiesHTML;
    };
    const abilitiesButton = document.querySelector('.tab-button[data-tab="abilities"]');
    if (abilitiesButton) {
        abilitiesButton.addEventListener('click', () => {
            const abilitiesTab = document.getElementById('abilities');
            if (abilitiesTab.querySelector('.loading-container')) {
                fetchAbilities();
            }
        });
        if (abilitiesButton.classList.contains('active')) {
            fetchAbilities();
        }
    }
});