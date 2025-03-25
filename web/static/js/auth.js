document.addEventListener('DOMContentLoaded', function() {
    // Handle login form submission
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            
            // Simple validation
            if (!email || !password) {
                showError('Please fill in all fields');
                return;
            }
            
            // Submit the form
            this.submit();
        });
    }
    
    // Handle register form submission
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            const username = document.getElementById('username').value;
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const confirmPassword = document.getElementById('confirm_password').value;
            
            // Simple validation
            if (!username || !email || !password || !confirmPassword) {
                showError('Please fill in all fields');
                return;
            }
            
            if (password !== confirmPassword) {
                showError('Passwords do not match');
                return;
            }
            
            if (password.length < 8) {
                showError('Password must be at least 8 characters long');
                return;
            }
            
            // Submit the form
            this.submit();
        });
    }
    
    // Helper function to show error messages
    function showError(message) {
        // Check if error div already exists
        let errorDiv = document.querySelector('.error-message');
        
        if (!errorDiv) {
            // Create error div if it doesn't exist
            errorDiv = document.createElement('div');
            errorDiv.className = 'error-message';
            
            // Insert before the form
            const form = document.querySelector('form');
            form.parentNode.insertBefore(errorDiv, form);
        }
        
        errorDiv.textContent = message;
        
        // Auto-dismiss after 5 seconds
        setTimeout(() => {
            errorDiv.remove();
        }, 5000);
    }
});