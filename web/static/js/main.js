// Main JavaScript file
document.addEventListener('DOMContentLoaded', function() {
    console.log('User profile page loaded');
    
    // Add timestamp to show when the page was loaded
    const timestamp = document.createElement('div');
    timestamp.className = 'timestamp';
    timestamp.textContent = 'Page loaded at: ' + new Date().toLocaleTimeString();
    timestamp.style.textAlign = 'center';
    timestamp.style.marginTop = '20px';
    timestamp.style.fontSize = '12px';
    timestamp.style.color = '#999';
    
    document.querySelector('.user-profile').appendChild(timestamp);
});