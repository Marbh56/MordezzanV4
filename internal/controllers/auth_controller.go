package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"mordezzanV4/internal/auth"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
}

type AuthController struct {
	userRepo  repositories.UserRepository
	tmpl      *template.Template
	jwtConfig auth.JWTConfig
}

func NewAuthController(userRepo repositories.UserRepository, tmpl *template.Template, jwtSecret string) *AuthController {
	return &AuthController{
		userRepo: userRepo,
		tmpl:     tmpl,
		jwtConfig: auth.JWTConfig{
			Secret:     jwtSecret,
			Issuer:     "mordezzanV4",
			Expiration: 24 * time.Hour, // 24 hour token expiration
		},
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request format"))
		return
	}

	// Check if either email or password is empty
	if loginReq.Email == "" || loginReq.Password == "" {
		apperrors.HandleValidationErrors(w, map[string]string{
			"credentials": "Email and password are required",
		})
		return
	}

	// Find user by email
	user, err := c.userRepo.GetUserByEmail(r.Context(), loginReq.Email)
	if err != nil {
		// Don't reveal whether the email exists or not for security
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid credentials"))
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password)); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid credentials"))
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, "user", c.jwtConfig)
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Create and send response
	resp := LoginResponse{
		Token:     token,
		ExpiresIn: int64(c.jwtConfig.Expiration.Seconds()),
		UserID:    user.ID,
		Username:  user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *AuthController) RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := c.tmpl.ExecuteTemplate(w, "login", nil); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *AuthController) RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := c.tmpl.ExecuteTemplate(w, "register", nil); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}
