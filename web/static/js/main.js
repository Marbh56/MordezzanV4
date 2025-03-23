document.addEventListener('DOMContentLoaded', function() {
    // Add animation class to fade in content
    document.body.classList.add('fade-in');
    
    // Check for error messages and display them
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const error = urlParams.get('error');
    
    if (error) {
        displayAlert(error, 'danger');
    }
    
    // Check for success messages
    const success = urlParams.get('success');
    if (success) {
        displayAlert(success, 'success');
    }
    
    // Initialize any interactive elements
    initializeTooltips();
    
    // Handle JWT token expiration
    checkTokenExpiration();
    
    // Add event listeners to all forms to handle submission errors
    setupFormErrorHandling();
});

/**
 * Display an alert message at the top of the page
 * @param {string} message - The message to display
 * @param {string} type - The type of alert (success, danger, warning, info)
 */
function displayAlert(message, type = 'info') {
    // Find the container to insert the alert
    const container = document.querySelector('.container') || document.body;
    const form = document.querySelector('form');
    
    // Create the alert element
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type}`;
    alertDiv.textContent = message;
    
    // Add a close button
    const closeButton = document.createElement('button');
    closeButton.innerHTML = '&times;';
    closeButton.className = 'alert-close';
    closeButton.onclick = function() {
        alertDiv.remove();
    };
    alertDiv.appendChild(closeButton);
    
    // Insert before the form or at the top of the container
    if (form) {
        form.parentNode.insertBefore(alertDiv, form);
    } else {
        const firstChild = container.firstChild;
        container.insertBefore(alertDiv, firstChild);
    }
    
    // Auto remove after 5 seconds for success/info messages
    if (type === 'success' || type === 'info') {
        setTimeout(() => {
            alertDiv.remove();
        }, 5000);
    }
}

/**
 * Initialize tooltips on elements with the 'data-tooltip' attribute
 */
function initializeTooltips() {
    const tooltipElements = document.querySelectorAll('[data-tooltip]');
    tooltipElements.forEach(element => {
        // Tooltip implementation can be added here
        // This is a placeholder for future enhancement
    });
}

/**
 * Check if the JWT token is expired and redirect to login if needed
 */
function checkTokenExpiration() {
    const token = localStorage.getItem('jwt_token');
    if (!token) return;
    
    try {
        // Parse the token (it's in format header.payload.signature)
        const payload = token.split('.')[1];
        const decodedPayload = atob(payload);
        const payloadObj = JSON.parse(decodedPayload);
        
        // Check if token is expired
        const currentTime = Math.floor(Date.now() / 1000);
        if (payloadObj.exp && payloadObj.exp < currentTime) {
            // Token is expired, redirect to login
            localStorage.removeItem('jwt_token');
            window.location.href = '/auth/login-page?error=Your session has expired. Please log in again.';
        }
    } catch (error) {
        console.error('Error checking token expiration:', error);
    }
}

/**
 * Set up error handling for all forms
 */
function setupFormErrorHandling() {
    const forms = document.querySelectorAll('form');
    forms.forEach(form => {
        form.addEventListener('submit', function(event) {
            // Only apply this to forms that submit to API endpoints
            if (form.action.includes('/api/')) {
                event.preventDefault();
                
                const formData = new FormData(form);
                const formObject = {};
                
                formData.forEach((value, key) => {
                    formObject[key] = value;
                });
                
                fetch(form.action, {
                    method: form.method,
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem('jwt_token')}`
                    },
                    body: JSON.stringify(formObject)
                })
                .then(response => {
                    if (!response.ok) {
                        return response.json().then(data => {
                            throw new Error(data.error || 'An error occurred');
                        });
                    }
                    return response.json();
                })
                .then(data => {
                    // Success handler
                    if (data.redirect) {
                        window.location.href = data.redirect;
                    } else {
                        displayAlert(data.message || 'Operation successful', 'success');
                    }
                })
                .catch(error => {
                    displayAlert(error.message, 'danger');
                });
            }
        });
    });
}