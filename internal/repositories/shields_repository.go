package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type ShieldRepository interface {
	GetShield(ctx context.Context, id int64) (*models.Shield, error)
	GetShieldByName(ctx context.Context, name string) (*models.Shield, error)
	ListShields(ctx context.Context) ([]*models.Shield, error)
	CreateShield(ctx context.Context, input *models.CreateShieldInput) (int64, error)
	UpdateShield(ctx context.Context, id int64, input *models.UpdateShieldInput) error
	DeleteShield(ctx context.Context, id int64) error
}

type SQLCShieldRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCShieldRepository(db *sql.DB) *SQLCShieldRepository {
	return &SQLCShieldRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCShieldRepository) GetShield(ctx context.Context, id int64) (*models.Shield, error) {
	shield, err := r.q.GetShield(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("shield", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbShieldToModel(shield), nil
}

func (r *SQLCShieldRepository) GetShieldByName(ctx context.Context, name string) (*models.Shield, error) {
	shield, err := r.q.GetShieldByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("shield", name)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbShieldToModel(shield), nil
}

func (r *SQLCShieldRepository) ListShields(ctx context.Context) ([]*models.Shield, error) {
	shields, err := r.q.ListShields(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.Shield, len(shields))
	for i, shield := range shields {
		result[i] = mapDbShieldToModel(shield)
	}
	return result, nil
}

func (r *SQLCShieldRepository) CreateShield(ctx context.Context, input *models.CreateShieldInput) (int64, error) {
	result, err := r.q.CreateShield(ctx, sqlcdb.CreateShieldParams{
		Name:            input.Name,
		Cost:            input.Cost,
		Weight:          int64(input.Weight),
		DefenseModifier: int64(input.DefenseModifier),
	})
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}
	return id, nil
}

func (r *SQLCShieldRepository) UpdateShield(ctx context.Context, id int64, input *models.UpdateShieldInput) error {
	_, err := r.GetShield(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.UpdateShield(ctx, sqlcdb.UpdateShieldParams{
		Name:            input.Name,
		Cost:            input.Cost,
		Weight:          int64(input.Weight),
		DefenseModifier: int64(input.DefenseModifier),
		ID:              id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCShieldRepository) DeleteShield(ctx context.Context, id int64) error {
	_, err := r.GetShield(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeleteShield(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func mapDbShieldToModel(shield sqlcdb.Shield) *models.Shield {
	return &models.Shield{
		ID:              shield.ID,
		Name:            shield.Name,
		Cost:            shield.Cost,
		Weight:          int(shield.Weight),
		DefenseModifier: int(shield.DefenseModifier),
		CreatedAt:       shield.CreatedAt,
		UpdatedAt:       shield.UpdatedAt,
	}
}
