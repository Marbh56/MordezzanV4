document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.querySelector('.auth-form');
    
    if (loginForm) {
        loginForm.addEventListener('submit', function(event) {
            event.preventDefault();
            
            // Clear any existing alerts
            const existingAlerts = document.querySelectorAll('.alert');
            existingAlerts.forEach(alert => alert.remove());
            
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const remember = document.getElementById('remember').checked;
            
            // Validate inputs
            if (!username || !password) {
                displayAlert('Username and password are required', 'danger');
                return;
            }
            
            // Create loading indicator
            const loadingAlert = document.createElement('div');
            loadingAlert.className = 'alert alert-info';
            loadingAlert.textContent = 'Logging in...';
            loginForm.parentNode.insertBefore(loadingAlert, loginForm);
            
            // Disable form elements during submission
            toggleFormElements(loginForm, false);
            
            // Prepare login data (controller expects email field)
            const loginData = {
                email: username, // Using the username input as email
                password: password,
                remember: remember
            };
            
            // Send login request
            fetch('/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(loginData)
            })
            .then(response => {
                if (!response.ok) {
                    return response.json().then(data => {
                        throw new Error(data.message || 'Login failed');
                    });
                }
                return response.json();
            })
            .then(data => {
                // Remove loading indicator
                loadingAlert.remove();
                
                // Store JWT token in localStorage
                localStorage.setItem('jwt_token', data.token);
                
                // Show success message
                displayAlert('Login successful! Redirecting...', 'success');
                
                // Redirect to dashboard
                setTimeout(() => {
                    window.location.href = '/';
                }, 1000);
            })
            .catch(error => {
                // Remove loading indicator
                loadingAlert.remove();
                
                // Re-enable form
                toggleFormElements(loginForm, true);
                
                // Display error - improved error message
                let errorMessage = error.message;
                if (errorMessage === 'Invalid credentials') {
                    errorMessage = 'Login failed. Please check your email and password.';
                }
                
                displayAlert(errorMessage, 'danger');
                
                // Log the error for debugging
                console.error('Login error:', error);
            });
        });
    }
    
    // Toggle form elements (enable/disable)
    function toggleFormElements(form, enabled) {
        const formElements = form.elements;
        for (let i = 0; i < formElements.length; i++) {
            formElements[i].disabled = !enabled;
        }
    }
    
    // Display alert message
    function displayAlert(message, type) {
        const alertDiv = document.createElement('div');
        alertDiv.className = `alert alert-${type}`;
        alertDiv.textContent = message;
        
        const form = document.querySelector('.auth-form');
        form.parentNode.insertBefore(alertDiv, form);
    }
});