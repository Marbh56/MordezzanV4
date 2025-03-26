document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('characterForm');
    const characterIdInput = document.getElementById('characterId');
    const isEditMode = characterIdInput !== null;
    
    if (form) {
        form.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const formData = {
                user_id: document.body.getAttribute('data-user-id') || 1,
                name: document.getElementById('name').value,
                class: document.getElementById('class').value,
                level: parseInt(document.getElementById('level').value, 10),
                experience_points: parseInt(document.getElementById('experience_points').value, 10) || 0,
                strength: parseInt(document.getElementById('strength').value, 10),
                dexterity: parseInt(document.getElementById('dexterity').value, 10),
                constitution: parseInt(document.getElementById('constitution').value, 10),
                intelligence: parseInt(document.getElementById('intelligence').value, 10),
                wisdom: parseInt(document.getElementById('wisdom').value, 10),
                charisma: parseInt(document.getElementById('charisma').value, 10),
                max_hit_points: parseInt(document.getElementById('max_hit_points').value, 10),
                current_hit_points: parseInt(document.getElementById('current_hit_points').value, 10),
                temporary_hit_points: parseInt(document.getElementById('temporary_hit_points').value, 10) || 0
            };
            
            try {
                let response;
                let characterId;
                
                if (isEditMode) {
                    // Edit existing character
                    characterId = characterIdInput.value;
                    response = await fetch(`/api/characters/${characterId}`, {
                        method: 'PUT',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(formData)
                    });
                } else {
                    // Create new character
                    response = await fetch('/api/characters', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(formData)
                    });
                }
                
                let data;
                try {
                    data = await response.json();
                } catch (parseError) {
                    console.error('Failed to parse response as JSON:', parseError);
                    throw new Error('Server returned an invalid response');
                }
                
                if (!response.ok) {
                    console.error('Server error details:', data);
                    
                    // Check if we have validation errors
                    if (data.fields) {
                        const errorMessages = Object.entries(data.fields)
                            .map(([field, message]) => `${field}: ${message}`)
                            .join('\n');
                        throw new Error(`Validation errors:\n${errorMessages}`);
                    } else {
                        throw new Error(data.message || 'Failed to save character');
                    }
                }
                
                // Redirect to character view page
                window.location.href = `/characters/view/${isEditMode ? characterId : data.id}`;
                
            } catch (error) {
                console.error('Error saving character:', error);
                showError(error.message);
            }
        });
    }
    
    function showError(message) {
        let errorDiv = document.querySelector('.error-message');
        if (!errorDiv) {
            errorDiv = document.createElement('div');
            errorDiv.className = 'error-message';
            form.insertBefore(errorDiv, form.firstChild);
        }
        errorDiv.textContent = message;
        errorDiv.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
});