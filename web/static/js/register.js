document.addEventListener('DOMContentLoaded', function() {
    console.log("Register.js loaded");
    const form = document.querySelector('.auth-form');
    
    if (!form) {
        console.error("Form element not found!");
        return;
    }
    
    console.log("Form element found:", form);
    
    // Add form submission validation
    form.addEventListener('submit', function(event) {
        // Prevent default form behavior
        event.preventDefault();
        console.log("Form submit event captured");
        
        // Get form elements
        const username = document.getElementById('username');
        const email = document.getElementById('email');
        const password = document.getElementById('password');
        
        if (!username || !email || !password) {
            console.error("Required form elements not found!", { 
                username: username ? "Found" : "Not found",
                email: email ? "Found" : "Not found",
                password: password ? "Found" : "Not found"
            });
            alert("Form is missing required elements. Check console for details.");
            return;
        }
        
        console.log("Form data collected:", {
            username: username.value,
            email: email.value,
            password: "********"
        });
        
        // Very simple validation
        if (!username.value || !email.value || !password.value) {
            console.error("Missing required fields");
            alert("Please fill in all required fields");
            return;
        }
        
        // Show that we're attempting to submit
        alert("Attempting to submit form to server...");
        
        // Collect form data
        const formData = {
            username: username.value,
            email: email.value,
            password: password.value
        };
        
        // Simple request with detailed logging
        console.log("Sending fetch request to /api/users");
        
        fetch('/api/users', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        })
        .then(response => {
            console.log("Received server response:", response.status, response.statusText);
            return response.json().catch(e => {
                console.log("Error parsing JSON response:", e);
                return { error: "Could not parse server response" };
            });
        })
        .then(data => {
            console.log("Response data:", data);
            alert("Server response received! Check console for details.");
            
            if (data.error) {
                alert("Error: " + data.error);
            } else {
                alert("Registration successful! Redirecting to login page...");
                window.location.href = '/auth/login-page';
            }
        })
        .catch(error => {
            console.error("Fetch error:", error);
            alert("Error during fetch request: " + error.message);
        });
    });
});