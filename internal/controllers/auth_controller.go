package controllers

import (
	"encoding/json"
	"html/template"
	"mordezzanV4/internal/errors"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents the request body for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response body for login
type LoginResponse struct {
	Success  bool   `json:"success"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Message  string `json:"message,omitempty"`
}

// AuthController handles authentication-related requests
type AuthController struct {
	userRepo       repositories.UserRepository
	tmpl           *template.Template
	sessionManager *scs.SessionManager
}

// NewAuthController creates a new AuthController instance
func NewAuthController(userRepo repositories.UserRepository, tmpl *template.Template, sessionManager *scs.SessionManager) *AuthController {
	return &AuthController{
		userRepo:       userRepo,
		tmpl:           tmpl,
		sessionManager: sessionManager,
	}
}

// Login handles user login and session creation
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest

	// Check if this is a form submission or JSON API request
	contentType := r.Header.Get("Content-Type")
	acceptHeader := r.Header.Get("Accept")
	isAPIRequest := contentType == "application/json" || acceptHeader == "application/json"

	if isAPIRequest {
		// Handle JSON API request
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			errors.HandleError(w, errors.NewBadRequest("Invalid request format"))
			return
		}
	} else {
		// Handle form submission
		if err := r.ParseForm(); err != nil {
			errors.HandleError(w, errors.NewBadRequest("Invalid form submission"))
			return
		}
		loginReq.Email = r.FormValue("email")
		loginReq.Password = r.FormValue("password")
	}

	// Validate login request
	if loginReq.Email == "" || loginReq.Password == "" {
		if isAPIRequest {
			errors.HandleValidationErrors(w, map[string]string{
				"credentials": "Email and password are required",
			})
		} else {
			data := map[string]interface{}{
				"Error": "Email and password are required",
			}
			c.tmpl.ExecuteTemplate(w, "login", data)
		}
		return
	}

	// Get user by email
	user, err := c.userRepo.GetUserByEmail(r.Context(), loginReq.Email)
	if err != nil {
		if isAPIRequest {
			errors.HandleError(w, errors.NewBadRequest("Invalid credentials"))
		} else {
			data := map[string]interface{}{
				"Error": "Invalid email or password",
			}
			c.tmpl.ExecuteTemplate(w, "login", data)
		}
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password)); err != nil {
		if isAPIRequest {
			errors.HandleError(w, errors.NewBadRequest("Invalid credentials"))
		} else {
			data := map[string]interface{}{
				"Error": "Invalid email or password",
			}
			c.tmpl.ExecuteTemplate(w, "login", data)
		}
		return
	}

	// Login successful - create session
	c.sessionManager.Put(r.Context(), "userID", user.ID)
	c.sessionManager.Put(r.Context(), "username", user.Username)
	c.sessionManager.Put(r.Context(), "isAuthenticated", true)

	logger.Info("User logged in: %s (ID: %d)", user.Username, user.ID)

	// Return response based on request type
	if isAPIRequest {
		resp := LoginResponse{
			Success:  true,
			UserID:   user.ID,
			Username: user.Username,
			Message:  "Login successful",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			errors.HandleError(w, errors.NewInternalError(err))
		}
	} else {
		// Redirect to dashboard
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Logout handles user logout by destroying the session
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user info for logging before destroying session
	if userID := c.sessionManager.GetInt64(r.Context(), "userID"); userID > 0 {
		logger.Info("User logged out: ID %d", userID)
	}

	// Destroy the session
	c.sessionManager.Destroy(r.Context())

	// Check if this is an API request
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Logged out successfully",
		})
		return
	}

	// Redirect to home page for browser requests
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// RenderLoginPage renders the login page template
func (c *AuthController) RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	// Check if user is already logged in
	if c.sessionManager.GetBool(r.Context(), "isAuthenticated") {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Render login template
	if err := c.tmpl.ExecuteTemplate(w, "login", nil); err != nil {
		errors.HandleError(w, errors.NewInternalError(err))
	}
}

// RenderRegisterPage renders the registration page template
func (c *AuthController) RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	// Check if user is already logged in
	if c.sessionManager.GetBool(r.Context(), "isAuthenticated") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Render register template
	if err := c.tmpl.ExecuteTemplate(w, "register", nil); err != nil {
		errors.HandleError(w, errors.NewInternalError(err))
	}
}

// Register handles new user registration
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUserInput

	// Check if this is a form submission or JSON API request
	contentType := r.Header.Get("Content-Type")
	acceptHeader := r.Header.Get("Accept")
	isAPIRequest := contentType == "application/json" || acceptHeader == "application/json"

	if isAPIRequest {
		// Handle JSON API request
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			errors.HandleError(w, errors.NewBadRequest("Invalid request format"))
			return
		}
	} else {
		// Handle form submission
		if err := r.ParseForm(); err != nil {
			errors.HandleError(w, errors.NewBadRequest("Invalid form submission"))
			return
		}
		input.Username = r.FormValue("username")
		input.Email = r.FormValue("email")
		input.Password = r.FormValue("password")
	}

	// Validate input
	if err := input.Validate(); err != nil {
		if isAPIRequest {
			errors.HandleError(w, err)
		} else {
			data := map[string]interface{}{
				"Error": err.Error(),
			}
			c.tmpl.ExecuteTemplate(w, "register", data)
		}
		return
	}

	// Check if email already exists
	existingUser, _ := c.userRepo.GetUserByEmail(r.Context(), input.Email)
	if existingUser != nil {
		if isAPIRequest {
			errors.HandleValidationErrors(w, map[string]string{
				"email": "Email address already in use",
			})
		} else {
			data := map[string]interface{}{
				"Error": "Email address already in use",
			}
			c.tmpl.ExecuteTemplate(w, "register", data)
		}
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		errors.HandleError(w, errors.NewInternalError(err))
		return
	}

	// Create user
	userID, err := c.userRepo.CreateUser(r.Context(), input.Username, input.Email, string(hashedPassword))
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	// Log the creation
	logger.Info("New user registered: %s (ID: %d)", input.Username, userID)

	// Automatically log in the new user
	c.sessionManager.Put(r.Context(), "userID", userID)
	c.sessionManager.Put(r.Context(), "username", input.Username)
	c.sessionManager.Put(r.Context(), "isAuthenticated", true)

	// Return response based on request type
	if isAPIRequest {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"user_id": userID,
			"message": "Registration successful",
		})
	} else {
		// Redirect to dashboard
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// GetCurrentUser gets the currently authenticated user
func (c *AuthController) GetCurrentUser(r *http.Request) (*models.User, error) {
	userID := c.sessionManager.GetInt64(r.Context(), "userID")
	if userID == 0 {
		return nil, errors.NewUnauthorized("Not authenticated")
	}

	return c.userRepo.GetUser(r.Context(), userID)
}
