function setupNotesTab() {
    const saveButton = document.getElementById('btn-save-notes');
    const notesTextarea = document.getElementById('character-notes');
    const notesMessage = document.getElementById('notes-message');

    // Load existing notes
    loadCharacterNotes();

    // Save notes button
    saveButton.addEventListener('click', async function () {
        const notes = notesTextarea.value;
        try {
            // In the future, implement an API endpoint for saving notes
            // For now, just store in localStorage
            localStorage.setItem(`character_notes_${characterId}`, notes);

            // Show success message
            notesMessage.textContent = 'Notes saved successfully!';
            notesMessage.style.color = '#2ecc71';
            notesMessage.style.display = 'block';

            // Hide message after 3 seconds
            setTimeout(() => {
                notesMessage.style.display = 'none';
            }, 3000);
        } catch (error) {
            console.error('Error saving notes:', error);
            notesMessage.textContent = 'Failed to save notes. Please try again.';
            notesMessage.style.color = '#e74c3c';
            notesMessage.style.display = 'block';
        }
    });
}

function loadCharacterNotes() {
    const notesTextarea = document.getElementById('character-notes');
    // For now, just retrieve from localStorage
    const savedNotes = localStorage.getItem(`character_notes_${characterId}`);
    if (savedNotes) {
        notesTextarea.value = savedNotes;
    }
}