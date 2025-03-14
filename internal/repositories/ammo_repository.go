package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type AmmoRepository interface {
	GetAmmo(ctx context.Context, id int64) (*models.Ammo, error)
	GetAmmoByName(ctx context.Context, name string) (*models.Ammo, error)
	ListAmmo(ctx context.Context) ([]*models.Ammo, error)
	CreateAmmo(ctx context.Context, input *models.CreateAmmoInput) (int64, error)
	UpdateAmmo(ctx context.Context, id int64, input *models.UpdateAmmoInput) error
	DeleteAmmo(ctx context.Context, id int64) error
}

type SQLCAmmoRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCAmmoRepository(db *sql.DB) *SQLCAmmoRepository {
	return &SQLCAmmoRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCAmmoRepository) GetAmmo(ctx context.Context, id int64) (*models.Ammo, error) {
	ammo, err := r.q.GetAmmo(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("ammo", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbAmmoToModel(ammo), nil
}

func (r *SQLCAmmoRepository) GetAmmoByName(ctx context.Context, name string) (*models.Ammo, error) {
	ammo, err := r.q.GetAmmoByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("ammo", name)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbAmmoToModel(ammo), nil
}

func (r *SQLCAmmoRepository) ListAmmo(ctx context.Context) ([]*models.Ammo, error) {
	ammoList, err := r.q.ListAmmo(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.Ammo, len(ammoList))
	for i, ammo := range ammoList {
		result[i] = mapDbAmmoToModel(ammo)
	}
	return result, nil
}

func (r *SQLCAmmoRepository) CreateAmmo(ctx context.Context, input *models.CreateAmmoInput) (int64, error) {
	result, err := r.q.CreateAmmo(ctx, sqlcdb.CreateAmmoParams{
		Name:   input.Name,
		Cost:   input.Cost,
		Weight: int64(input.Weight),
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

func (r *SQLCAmmoRepository) UpdateAmmo(ctx context.Context, id int64, input *models.UpdateAmmoInput) error {
	_, err := r.GetAmmo(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.UpdateAmmo(ctx, sqlcdb.UpdateAmmoParams{
		Name:   input.Name,
		Cost:   input.Cost,
		Weight: int64(input.Weight),
		ID:     id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCAmmoRepository) DeleteAmmo(ctx context.Context, id int64) error {
	_, err := r.GetAmmo(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeleteAmmo(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func mapDbAmmoToModel(ammo sqlcdb.Ammo) *models.Ammo {
	return &models.Ammo{
		ID:        ammo.ID,
		Name:      ammo.Name,
		Cost:      ammo.Cost,
		Weight:    int(ammo.Weight),
		CreatedAt: ammo.CreatedAt,
		UpdatedAt: ammo.UpdatedAt,
	}
}
