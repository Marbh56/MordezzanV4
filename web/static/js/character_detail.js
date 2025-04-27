document.addEventListener('DOMContentLoaded', function() {
    const weaponMasteryClasses = [
        'Fighter', 'Barbarian', 'Berserker', 'Cataphract',
        'Huntsman', 'Paladin', 'Ranger', 'Assassin'
    ];
    
    const thiefSkillClasses = [
        'Thief', 'Assassin', 'Bard', 'Legerdemainist', 'Scout'
    ];
    
    // Handle visibility for class-specific tabs
    const characterClassElement = document.querySelector('.character-meta');
    if (characterClassElement) {
        const classText = characterClassElement.textContent;
        const characterClass = classText.match(/Level \d+ (.+)/)?.[1];
        
        // Handle weapon mastery tab visibility
        const hasWeaponMastery = weaponMasteryClasses.includes(characterClass);
        const masteryButton = document.querySelector('.tab-button[data-tab="weapon-mastery"]');
        if (masteryButton && !hasWeaponMastery) {
            masteryButton.style.display = 'none';
        }
        
        // Handle thief skills tab visibility
        const hasThiefSkills = thiefSkillClasses.includes(characterClass);
        const thiefSkillsButton = document.querySelector('.tab-button[data-tab="thief-skills"]');
        if (thiefSkillsButton && !hasThiefSkills) {
            thiefSkillsButton.style.display = 'none';
        }
    }
    
    // Tab switching functionality
    const tabButtons = document.querySelectorAll('.tab-button');
    const tabContents = document.querySelectorAll('.tab-content');
    
    // Verify all tabs have corresponding content
    tabButtons.forEach(button => {
        const tabId = button.getAttribute('data-tab');
        const tabContent = document.getElementById(tabId);
        
        if (!tabContent) {
            console.error(`Tab content not found for tab: ${tabId}`);
        }
    });
    
    // Add click handlers for tab switching
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabId = button.getAttribute('data-tab');
            
            // Verify the target tab content exists
            const targetTab = document.getElementById(tabId);
            if (!targetTab) {
                console.error(`Cannot switch to tab: Tab content #${tabId} not found`);
                return;
            }
            
            // Remove active class from all tabs
            tabButtons.forEach(btn => btn.classList.remove('active'));
            tabContents.forEach(content => content.classList.remove('active'));
            
            // Add active class to selected tab
            button.classList.add('active');
            targetTab.classList.add('active');
            
            // Console log for debugging
            console.log(`Switched to tab: ${tabId}`);
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

    // Add fetch thief skills functionality
    const fetchThiefSkills = async () => {
        console.log('Attempting to fetch thief skills data');
        const thiefSkillsTab = document.getElementById('thief-skills');
        if (!thiefSkillsTab) {
            console.error('Thief skills tab element not found in DOM');
            return;
        }
        
        try {
            console.log(`Fetching thief skills for character ID: ${characterId}`);
            const response = await fetch(`/api/characters/${characterId}/thief-skills`);
            console.log('Thief skills API response:', response);
            
            if (!response.ok) {
                throw new Error(`Failed to fetch thief skills: ${response.status} ${response.statusText}`);
            }
            
            const data = await response.json();
            console.log('Thief skills data received:', data);
            renderThiefSkills(data);
        } catch (error) {
            console.error('Error fetching thief skills:', error);
            thiefSkillsTab.innerHTML = `
                <div class="error-message">
                    Failed to load thief skills: ${error.message}
                </div>
            `;
        }
    };
    
    const renderThiefSkills = (skillsData) => {
        console.log('Rendering thief skills with data:', skillsData);
        const thiefSkillsTab = document.getElementById('thief-skills');
        
        // Check if skillsData is an array or if it's an object with a skills property
        const skills = Array.isArray(skillsData) ? skillsData : 
                      (skillsData && skillsData.skills ? skillsData.skills : []);
        
        if (skills.length > 0) {
            // Log the first skill to see its structure
            console.log('Sample skill object structure:', skills[0]);
        }
        
        if (!skills || skills.length === 0) {
            thiefSkillsTab.innerHTML = `
                <div class="empty-state">
                    <h3 class="empty-title">No Thief Skills Found</h3>
                    <p class="empty-description">This character doesn't have any thief skills at the current level.</p>
                </div>
            `;
            return;
        }
        
        // Render the skills using the actual data structure
        let skillsHTML = `
            <h2>Thief Skills</h2>
            <div class="thief-skills-container">
        `;
        
        skills.forEach(skill => {
            // Try to access success chance from various possible property names
            const successChance = 
                skill.success_chance || 
                skill.successChance || 
                skill.SuccessChance || 
                (typeof skill.Value !== 'undefined' ? skill.Value : 
                    (typeof skill.value !== 'undefined' ? skill.value : 'Unknown'));
            
            skillsHTML += `
                <div class="skill-card">
                    <h3 class="skill-name">${skill.Name || skill.name}</h3>
                    <div class="skill-value">${successChance}</div>
                    <div class="skill-description">${skill.Description || skill.description || ''}</div>
                </div>
            `;
        });
        
        skillsHTML += `</div>`;
        thiefSkillsTab.innerHTML = skillsHTML;
    };

    // Add event listener for thief skills tab
    const thiefSkillsButton = document.querySelector('.tab-button[data-tab="thief-skills"]');
    if (thiefSkillsButton) {
        console.log('Thief skills button found in DOM');
        thiefSkillsButton.addEventListener('click', () => {
            console.log('Thief skills tab clicked');
            const thiefSkillsTab = document.getElementById('thief-skills');
            if (thiefSkillsTab) {
                if (thiefSkillsTab.querySelector('.loading-container') || thiefSkillsTab.innerHTML.trim() === '') {
                    console.log('Thief skills tab is empty or loading, fetching data');
                    fetchThiefSkills();
                } else {
                    console.log('Thief skills data already loaded');
                }
            } else {
                console.error('Thief skills tab container not found');
            }
        });
        
        // Load thief skills if this tab is active by default
        if (thiefSkillsButton.classList.contains('active')) {
            console.log('Thief skills tab is active by default, fetching data');
            fetchThiefSkills();
        }
    } else {
        console.warn('Thief skills button not found in DOM');
    }
});