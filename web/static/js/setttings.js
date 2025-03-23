document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('accountSettingsForm');
    const changePasswordCheckbox = document.getElementById('change_password');
    const passwordFields = document.getElementById('passwordFields');
    const currentPassword = document.getElementById('current_password');
    const newPassword = document.getElementById('new_password');
    const confirmPassword = document.getElementById('confirm_password');
    const cancelButton = document.getElementById('cancelButton');

    // Toggle password fields visibility
    changePasswordCheckbox.addEventListener('change', function() {
        if (this.checked) {
            passwordFields.classList.add('active');
            currentPassword.setAttribute('required', '');
            newPassword.setAttribute('required', '');
            confirmPassword.setAttribute('required', '');
        } else {
            passwordFields.classList.remove('active');
            currentPassword.removeAttribute('required');
            newPassword.removeAttribute('required');
            confirmPassword.removeAttribute('required');
            currentPassword.value = '';
            newPassword.value = '';
            confirmPassword.value = '';
        }
    });

    // Validate new password and confirmation match
    confirmPassword.addEventListener('input', function() {
        if (newPassword.value !== confirmPassword.value) {
            confirmPassword.setCustomValidity('Passwords do not match');
        } else {
            confirmPassword.setCustomValidity('');
        }
    });

    newPassword.addEventListener('input', function() {
        if (confirmPassword.value && newPassword.value !== confirmPassword.value) {
            confirmPassword.setCustomValidity('Passwords do not match');
        } else {
            confirmPassword.setCustomValidity('');
        }
    });

    // Handle form submission
    form.addEventListener('submit', function(event) {
        event.preventDefault();
        removeAllAlerts();
        
        // Basic client-side validation
        let hasError = false;
        let errorMessage = '';
        
        const username = document.getElementById('username');
        if (username.value.trim().length < 3) {
            errorMessage = 'Username must be at least 3 characters long';
            hasError = true;
        }
        
        const email = document.getElementById('email');
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email.value)) {
            errorMessage = 'Please enter a valid email address';
            hasError = true;
        }
        
        if (changePasswordCheckbox.checked) {
            if (newPassword.value.length < 8) {
                errorMessage = 'New password must be at least 8 characters long';
                hasError = true;
            }
            
            if (newPassword.value !== confirmPassword.value) {
                errorMessage = 'New passwords do not match';
                hasError = true;
            }
        }
        
        if (hasError) {
            displayAlert(errorMessage, 'danger');
            window.scrollTo(0, 0);
        } else {
            submitSettings();
        }
    });
    
    cancelButton.addEventListener('click', function() {
        window.history.back();
    });
    
    function submitSettings() {
        displayAlert('Saving changes...', 'info');
        
        const data = {
            username: document.getElementById('username').value,
            email: document.getElementById('email').value
        };
        
        if (changePasswordCheckbox.checked) {
            data.current_password = currentPassword.value;
            data.new_password = newPassword.value;
        }
        
        // Use the API endpoint for updating user settings
        fetch('/api/user/settings', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('jwt_token')}`
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.message || data.fields ? Object.values(data.fields)[0] : 'Failed to update settings');
                });
            }
            return response.json();
        })
        .then(data => {
            removeAllAlerts();
            displayAlert('Settings updated successfully', 'success');
            
            // Update UI elements if username changed
            if (data.username && document.querySelector('.user-name')) {
                document.querySelector('.user-name').textContent = data.username;
                if (document.querySelector('.user-avatar')) {
                    document.querySelector('.user-avatar').textContent = data.username.charAt(0).toUpperCase();
                }
            }
        })
        .catch(error => {
            removeAllAlerts();
            displayAlert(error.message, 'danger');
        });
    }
    
    function displayAlert(message, type) {
        const alertContainer = document.querySelector('.settings-container');
        const settingsCard = document.querySelector('.settings-card');
        
        if (!alertContainer || !settingsCard) {
            console.error('Alert container or settings card not found');
            return;
        }
        
        const alertDiv = document.createElement('div');
        alertDiv.className = `alert alert-${type}`;
        alertDiv.textContent = message;
        
        const closeButton = document.createElement('button');
        closeButton.innerHTML = '&times;';
        closeButton.className = 'alert-close';
        closeButton.onclick = function() {
            alertDiv.remove();
        };
        
        alertDiv.appendChild(closeButton);
        alertContainer.insertBefore(alertDiv, settingsCard);
        
        // Auto-remove success and info alerts after 5 seconds
        if (type === 'success' || type === 'info') {
            setTimeout(() => {
                if (alertDiv.parentNode) {
                    alertDiv.remove();
                }
            }, 5000);
        }
    }
    
    function removeAllAlerts() {
        const alerts = document.querySelectorAll('.alert');
        alerts.forEach(alert => {
            alert.remove();
        });
    }
});