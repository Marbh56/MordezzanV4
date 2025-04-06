document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('characterForm');
    const characterIdInput = document.getElementById('characterId');
    const isEditMode = characterIdInput !== null;
    
    // Get user ID from meta tag instead of body attribute
    function getUserId() {
        const userIdMeta = document.querySelector('meta[name="user-id"]');
        if (userIdMeta) {
            return userIdMeta.getAttribute('content');
        }
        console.error('User ID meta tag not found!');
        return null;
    }
    
    if (form) {
        form.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const userId = getUserId();
            if (!userId) {
                showError('User ID not found. Please try logging out and back in.');
                return;
            }
            
            const formData = {
                user_id: parseInt(userId, 10),
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
                    characterId = characterIdInput.value;
                    response = await fetch(`/api/characters/${characterId}`, {
                        method: 'PUT',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(formData)
                    });
                } else {
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
                    if (data.fields) {
                        const errorMessages = Object.entries(data.fields)
                            .map(([field, message]) => `${field}: ${message}`)
                            .join('\n');
                        throw new Error(`Validation errors:\n${errorMessages}`);
                    } else {
                        throw new Error(data.message || 'Failed to save character');
                    }
                }
                
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