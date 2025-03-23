document.addEventListener('DOMContentLoaded', function() {
    // Mobile navigation toggle
    const navbarToggle = document.getElementById('navbar-toggle');
    const navbarMenu = document.getElementById('navbar-menu');
    
    if (navbarToggle && navbarMenu) {
        navbarToggle.addEventListener('click', function() {
            navbarMenu.classList.toggle('active');
            
            // Toggle aria-expanded for accessibility
            const expanded = navbarToggle.getAttribute('aria-expanded') === 'true' || false;
            navbarToggle.setAttribute('aria-expanded', !expanded);
        });
    }
    
    // Mobile dropdown toggles
    const dropdownToggles = document.querySelectorAll('.dropdown-toggle');
    
    dropdownToggles.forEach(toggle => {
        toggle.addEventListener('click', function(e) {
            // Only apply this behavior on mobile
            if (window.innerWidth <= 768) {
                e.preventDefault();
                
                const parent = this.parentNode;
                parent.classList.toggle('active');
                
                // Toggle aria-expanded for accessibility
                const expanded = this.getAttribute('aria-expanded') === 'true' || false;
                this.setAttribute('aria-expanded', !expanded);
            }
        });
    });
    
    // Handle logout
    const logoutLink = document.querySelector('.logout-link');
    if (logoutLink) {
        logoutLink.addEventListener('click', function(e) {
            e.preventDefault();
            
            // Clear the JWT token
            localStorage.removeItem('jwt_token');
            
            // Redirect to login page
            window.location.href = '/auth/login-page?success=You have been logged out successfully';
        });
    }
    
    // Close the mobile menu when clicking outside
    document.addEventListener('click', function(e) {
        if (navbarMenu && navbarMenu.classList.contains('active')) {
            if (!navbarMenu.contains(e.target) && e.target !== navbarToggle) {
                navbarMenu.classList.remove('active');
                navbarToggle.setAttribute('aria-expanded', 'false');
            }
        }
    });
    
    // Close dropdown menus when clicking outside on mobile
    document.addEventListener('click', function(e) {
        if (window.innerWidth <= 768) {
            const activeDropdowns = document.querySelectorAll('.dropdown.active');
            activeDropdowns.forEach(dropdown => {
                if (!dropdown.contains(e.target)) {
                    dropdown.classList.remove('active');
                    const toggle = dropdown.querySelector('.dropdown-toggle');
                    if (toggle) {
                        toggle.setAttribute('aria-expanded', 'false');
                    }
                }
            });
        }
    });
});