document.addEventListener('DOMContentLoaded', function() {
    // Tab switching
    const tabButtons = document.querySelectorAll('.tab-button');
    const tabContents = document.querySelectorAll('.tab-content');

    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tabId = button.getAttribute('data-tab');
            
            // Remove active class from all buttons and contents
            tabButtons.forEach(btn => btn.classList.remove('active'));
            tabContents.forEach(content => content.classList.remove('active'));
            
            // Add active class to current button and content
            button.classList.add('active');
            document.getElementById(tabId).classList.add('active');
        });
    });

    // Get character ID from the URL
    const pathSegments = window.location.pathname.split('/');
    const characterId = pathSegments[pathSegments.indexOf('view') + 1];

    // HP Modal handling
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
                // Update the UI with new HP values
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

    // XP Modal handling
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
                // Get current XP from the display
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
                // Update the UI with new XP values
                xpElement.textContent = data.experience_points;
                
                // Check for level up
                const levelElement = document.querySelector('.character-meta');
                const currentLevel = parseInt(levelElement.textContent.match(/Level (\d+)/)[1]);
                
                if (data.level > currentLevel) {
                    alert(`Congratulations! You've reached level ${data.level}!`);
                    // Update level display without reloading
                    levelElement.textContent = levelElement.textContent.replace(
                        `Level ${currentLevel}`, 
                        `Level ${data.level}`
                    );
                    
                    // Refresh the page to update all level-dependent abilities
                    location.reload();
                }

                // Update XP needed for next level if displayed
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

    // Close modals when clicking outside
    window.addEventListener('click', (e) => {
        if (e.target === hpModal) {
            hpModal.style.display = 'none';
        }
        if (e.target === xpModal) {
            xpModal.style.display = 'none';
        }
    });
});