document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('createCharacterForm');
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
                current_hit_points: parseInt(document.getElementById('max_hit_points').value, 10),
                temporary_hit_points: 0
            };

            try {
                console.log('Sending character data:', formData);
                
                const response = await fetch('/api/characters', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData)
                });
                
                // Try to parse response as JSON
                let data;
                try {
                    data = await response.json();
                } catch (parseError) {
                    console.error('Failed to parse response as JSON:', parseError);
                    throw new Error('Server returned an invalid response');
                }
                
                if (!response.ok) {
                    console.error('Server error details:', data);
                    throw new Error(data.message || 'Failed to create character');
                }
                
                window.location.href = `/characters/view/${data.id}`;
            } catch (error) {
                console.error('Error creating character:', error);
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
        
        // Scroll to error message
        errorDiv.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
});