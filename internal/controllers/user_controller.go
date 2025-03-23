package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"mordezzanV4/internal/contextkeys"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository interface for the controller
type UserRepository interface {
	GetUser(ctx context.Context, id int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers(ctx context.Context) ([]*models.User, error)
	CreateUser(ctx context.Context, username, email, passwordHash string) (int64, error)
	UpdateUser(ctx context.Context, id int64, username, email string) error
	DeleteUser(ctx context.Context, id int64) error
}

type UserController struct {
	BaseController[
		models.User,
		models.CreateUserInput,
		models.UpdateUserInput,
		interface{ Validate() error },
	]
	userRepo UserRepository
	tmpl     *template.Template
}

func NewUserController(userRepo UserRepository, tmpl *template.Template) *UserController {
	return &UserController{
		BaseController: BaseController[
			models.User,
			models.CreateUserInput,
			models.UpdateUserInput,
			interface{ Validate() error },
		]{
			Repository:   nil, // We'll handle operations manually
			TemplateName: "user.html",
			Tmpl:         tmpl,
		},
		userRepo: userRepo,
		tmpl:     tmpl,
	}
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid user ID format"))
		return
	}

	user, err := c.userRepo.GetUser(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	if err := c.tmpl.ExecuteTemplate(w, "user.html", user); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.userRepo.ListUsers(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

	validationErrors := make(map[string]string)

	if len(strings.TrimSpace(input.Username)) < 3 {
		validationErrors["username"] = "Username must be at least 3 characters long"
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		validationErrors["email"] = "Invalid email address"
	}

	if len(input.Password) < 8 {
		validationErrors["password"] = "Password must be at least 8 characters long"
	}

	if len(validationErrors) > 0 {
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	// Check if email already exists
	existingUser, _ := c.userRepo.GetUserByEmail(r.Context(), input.Email)
	if existingUser != nil {
		validationErrors["email"] = "Email address already in use"
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	id, err := c.userRepo.CreateUser(r.Context(), input.Username, input.Email, string(hashedPassword))
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	user, err := c.userRepo.GetUser(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid user ID format"))
		return
	}

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

	validationErrors := make(map[string]string)

	if len(strings.TrimSpace(input.Username)) < 3 {
		validationErrors["username"] = "Username must be at least 3 characters long"
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		validationErrors["email"] = "Invalid email address"
	}

	if len(validationErrors) > 0 {
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	// Check if email already exists for another user
	existingUser, _ := c.userRepo.GetUserByEmail(r.Context(), input.Email)
	if existingUser != nil && existingUser.ID != id {
		validationErrors["email"] = "Email address already in use"
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	if err := c.userRepo.UpdateUser(r.Context(), id, input.Username, input.Email); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	updatedUser, err := c.userRepo.GetUser(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid user ID format"))
		return
	}

	if err := c.userRepo.DeleteUser(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *UserController) RenderSettingsPage(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userIDValue := r.Context().Value(contextkeys.UserIDKey)
	if userIDValue == nil {
		apperrors.HandleError(w, apperrors.NewUnauthorized("User not authenticated"))
		return
	}

	// Convert userID to int64
	userID, ok := userIDValue.(int64)
	if !ok {
		apperrors.HandleError(w, apperrors.NewInternalError(errors.New("invalid user ID format in context")))
		return
	}

	// Get user data
	user, err := c.userRepo.GetUser(r.Context(), userID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	data := map[string]interface{}{
		"User":            user,
		"IsAuthenticated": true,
	}

	// Add meta tag for user ID (used by JavaScript)
	data["MetaTags"] = []map[string]string{
		{
			"name":    "user-id",
			"content": strconv.FormatInt(user.ID, 10),
		},
	}

	if err := c.tmpl.ExecuteTemplate(w, "settings", data); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *UserController) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userIDValue := r.Context().Value(contextkeys.UserIDKey)
	if userIDValue == nil {
		apperrors.HandleError(w, apperrors.NewUnauthorized("User not authenticated"))
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		apperrors.HandleError(w, apperrors.NewInternalError(errors.New("invalid user ID format in context")))
		return
	}

	// Parse request body
	var input struct {
		Username        string `json:"username"`
		Email           string `json:"email"`
		CurrentPassword string `json:"current_password,omitempty"`
		NewPassword     string `json:"new_password,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

	// Validate input
	validationErrors := make(map[string]string)

	if len(strings.TrimSpace(input.Username)) < 3 {
		validationErrors["username"] = "Username must be at least 3 characters long"
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		validationErrors["email"] = "Invalid email address"
	}

	// Check if changing password
	isChangingPassword := input.CurrentPassword != "" && input.NewPassword != ""

	if isChangingPassword {
		// Get full user with password hash for verification
		fullUser, err := c.userRepo.GetUserByEmail(r.Context(), input.Email)
		if err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
			return
		}

		// Verify current password
		err = bcrypt.CompareHashAndPassword([]byte(fullUser.PasswordHash), []byte(input.CurrentPassword))
		if err != nil {
			validationErrors["current_password"] = "Current password is incorrect"
		}

		if len(input.NewPassword) < 8 {
			validationErrors["new_password"] = "New password must be at least 8 characters long"
		}
	}

	if len(validationErrors) > 0 {
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	// Check if email is already in use by another user
	existingUser, _ := c.userRepo.GetUserByEmail(r.Context(), input.Email)
	if existingUser != nil && existingUser.ID != userID {
		validationErrors["email"] = "Email address already in use by another account"
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	// Update user information
	if err := c.userRepo.UpdateUser(r.Context(), userID, input.Username, input.Email); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Update password if requested
	if isChangingPassword {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
			return
		}

		if err := c.updateUserPassword(r.Context(), userID, string(hashedPassword)); err != nil {
			apperrors.HandleError(w, err)
			return
		}
	}

	// Get updated user data to return in response
	updatedUser, err := c.userRepo.GetUser(r.Context(), userID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Return updated user information
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func (c *UserController) updateUserPassword(ctx context.Context, id int64, passwordHash string) error {
	// Check if the user repository implements the password update method
	if passwordUpdater, ok := c.userRepo.(interface {
		UpdateUserPassword(ctx context.Context, id int64, passwordHash string) error
	}); ok {
		return passwordUpdater.UpdateUserPassword(ctx, id, passwordHash)
	}

	return apperrors.NewInternalError(errors.New("password update not implemented"))
}
