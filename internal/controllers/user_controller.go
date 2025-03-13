package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userRepo repositories.UserRepository
	tmpl     *template.Template
}

func NewUserController(userRepo repositories.UserRepository, tmpl *template.Template) *UserController {
	return &UserController{
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
