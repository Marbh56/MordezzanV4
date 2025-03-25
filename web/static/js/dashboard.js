document.addEventListener('DOMContentLoaded', function() {
    // Add hover animations to buttons
    const buttons = document.querySelectorAll('.btn');
    buttons.forEach(button => {
        button.addEventListener('mouseenter', function() {
            this.style.transition = 'all 0.3s ease';
        });
    });
    
    // Optional: Fetch the latest character data periodically
    function refreshCharacterData() {
        fetch('/api/characters')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                if (data && data.length > 0) {
                    // If we need to update the UI with fresh data
                    // This could be implemented later when needed
                }
            })
            .catch(error => {
                console.error('Error fetching character data:', error);
            });
    }
    
    // Refresh character data every 5 minutes (if needed)
    // Uncomment the following line to enable automatic refreshing
    // setInterval(refreshCharacterData, 5 * 60 * 1000);
});