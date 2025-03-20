async function fetchSpells() {
    try {
        const token = localStorage.getItem('authToken');
        const spellsContainer = document.getElementById('spells-container');
        const spellsLoading = document.getElementById('spells-loading');
        const spellsEmpty = document.getElementById('spells-empty');

        spellsLoading.style.display = 'block';
        spellsEmpty.style.display = 'none';

        const response = await fetch(`/api/characters/${characterId}/spells`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to fetch spells');
        }

        const spells = await response.json();
        spellsLoading.style.display = 'none';

        if (spells.length === 0) {
            spellsEmpty.style.display = 'block';
            return;
        }

        // Clear previous spells
        spellsContainer.innerHTML = '';

        // Add each spell
        spells.forEach(spell => {
            // Implementation from original file
        });
    } catch (error) {
        console.error('Error fetching spells:', error);
        document.getElementById('spells-loading').style.display = 'none';
        document.getElementById('spells-empty').style.display = 'block';
        document.getElementById('spells-empty').textContent = 'Failed to load spells: ' + error.message;
    }
}