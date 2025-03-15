package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
	"time"
)

type SQLCUserRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

type UserRepository interface {
	GetUser(ctx context.Context, id int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers(ctx context.Context) ([]*models.User, error)
	CreateUser(ctx context.Context, username, email, passwordHash string) (int64, error)
	UpdateUser(ctx context.Context, id int64, username, email string) error
	DeleteUser(ctx context.Context, id int64) error
}

func NewSQLCUserRepository(db *sql.DB) *SQLCUserRepository {
	return &SQLCUserRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCUserRepository) GetUser(ctx context.Context, id int64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	user, err := r.q.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("user", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return &models.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *SQLCUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	user, err := r.q.GetFullUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("user", email)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return &models.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

func (r *SQLCUserRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	users, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.User, len(users))
	for i, user := range users {
		result[i] = &models.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}
	return result, nil
}

func (r *SQLCUserRepository) CreateUser(ctx context.Context, username, email, passwordHash string) (int64, error) {
	result, err := r.q.CreateUser(ctx, sqlcdb.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" {
			return 0, apperrors.NewValidationError("username", "Username already taken")
		}
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return 0, apperrors.NewValidationError("email", "Email already registered")
		}
		return 0, apperrors.NewDatabaseError(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}
	return id, nil
}

func (r *SQLCUserRepository) UpdateUser(ctx context.Context, id int64, username, email string) error {
	_, err := r.GetUser(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.UpdateUser(ctx, sqlcdb.UpdateUserParams{
		ID:       id,
		Username: username,
		Email:    email,
	})
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" {
			return apperrors.NewValidationError("username", "Username already taken")
		}
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return apperrors.NewValidationError("email", "Email already registered")
		}
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCUserRepository) DeleteUser(ctx context.Context, id int64) error {
	_, err := r.GetUser(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeleteUser(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}
