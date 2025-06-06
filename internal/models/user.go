package models

import (
	apperrors "mordezzanV4/internal/errors"
	"net/mail"
	"strings"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose in JSON responses
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if len(strings.TrimSpace(u.Username)) < 3 {
		return apperrors.NewValidationError("username", "username must be at least 3 characters long")
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return apperrors.NewValidationError("email", "invalid email address")
	}
	return nil
}

type CreateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (i *CreateUserInput) Validate() error {
	validationErrors := make(map[string]string)
	if len(strings.TrimSpace(i.Username)) < 3 {
		validationErrors["username"] = "Username must be at least 3 characters long"
	}
	if _, err := mail.ParseAddress(i.Email); err != nil {
		validationErrors["email"] = "Invalid email address"
	}
	if len(i.Password) < 8 {
		validationErrors["password"] = "Password must be at least 8 characters long"
	}
	if len(validationErrors) > 0 {
		for field, message := range validationErrors {
			return apperrors.NewValidationError(field, message)
		}
	}
	return nil
}

type UpdateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (i *UpdateUserInput) Validate() error {
	validationErrors := make(map[string]string)
	if len(strings.TrimSpace(i.Username)) < 3 {
		validationErrors["username"] = "Username must be at least 3 characters long"
	}
	if _, err := mail.ParseAddress(i.Email); err != nil {
		validationErrors["email"] = "Invalid email address"
	}
	if len(validationErrors) > 0 {
		for field, message := range validationErrors {
			return apperrors.NewValidationError(field, message)
		}
	}
	return nil
}
