package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type RingRepository interface {
	GetRing(ctx context.Context, id int64) (*models.Ring, error)
	GetRingByName(ctx context.Context, name string) (*models.Ring, error)
	ListRings(ctx context.Context) ([]*models.Ring, error)
	CreateRing(ctx context.Context, input *models.CreateRingInput) (int64, error)
	UpdateRing(ctx context.Context, id int64, input *models.UpdateRingInput) error
	DeleteRing(ctx context.Context, id int64) error
}

type SQLCRingRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCRingRepository(db *sql.DB) *SQLCRingRepository {
	return &SQLCRingRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCRingRepository) GetRing(ctx context.Context, id int64) (*models.Ring, error) {
	ring, err := r.q.GetRing(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("ring", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbRingToModel(ring), nil
}

func (r *SQLCRingRepository) GetRingByName(ctx context.Context, name string) (*models.Ring, error) {
	ring, err := r.q.GetRingByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("ring", name)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbRingToModel(ring), nil
}

func (r *SQLCRingRepository) ListRings(ctx context.Context) ([]*models.Ring, error) {
	rings, err := r.q.ListRings(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.Ring, len(rings))
	for i, ring := range rings {
		result[i] = mapDbRingToModel(ring)
	}
	return result, nil
}

func (r *SQLCRingRepository) CreateRing(ctx context.Context, input *models.CreateRingInput) (int64, error) {
	result, err := r.q.CreateRing(ctx, sqlcdb.CreateRingParams{
		Name:        input.Name,
		Description: input.Description,
		Cost:        input.Cost,
		Weight:      int64(input.Weight),
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

func (r *SQLCRingRepository) UpdateRing(ctx context.Context, id int64, input *models.UpdateRingInput) error {
	_, err := r.GetRing(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.UpdateRing(ctx, sqlcdb.UpdateRingParams{
		Name:        input.Name,
		Description: input.Description,
		Cost:        input.Cost,
		Weight:      int64(input.Weight),
		ID:          id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCRingRepository) DeleteRing(ctx context.Context, id int64) error {
	_, err := r.GetRing(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeleteRing(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func mapDbRingToModel(ring sqlcdb.Ring) *models.Ring {
	return &models.Ring{
		ID:          ring.ID,
		Name:        ring.Name,
		Description: ring.Description,
		Cost:        ring.Cost,
		Weight:      int(ring.Weight),
		CreatedAt:   ring.CreatedAt,
		UpdatedAt:   ring.UpdatedAt,
	}
}
