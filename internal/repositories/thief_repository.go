package repositories

import (
	"context"
	"database/sql"
	"errors"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

// ThiefDataRepository defines methods for accessing thief class data
type ThiefDataRepository interface {
	GetThiefClassData(ctx context.Context, level int) (*models.ThiefClassData, error)
	ListThiefClassData(ctx context.Context) ([]*models.ThiefClassData, error)
	GetNextThiefLevel(ctx context.Context, currentLevel int) (*models.ThiefClassData, error)
}

// SQLCThiefDataRepository implements ThiefDataRepository using SQLC
type SQLCThiefDataRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

// NewSQLCThiefDataRepository creates a new thief data repository
func NewSQLCThiefDataRepository(db *sql.DB) *SQLCThiefDataRepository {
	return &SQLCThiefDataRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

// GetThiefClassData retrieves thief class data for a specific level
func (r *SQLCThiefDataRepository) GetThiefClassData(ctx context.Context, level int) (*models.ThiefClassData, error) {
	thiefData, err := r.q.GetThiefClassData(ctx, int64(level))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return mapDbThiefDataToModel(thiefData), nil
}

// ListThiefClassData retrieves all thief class data levels
func (r *SQLCThiefDataRepository) ListThiefClassData(ctx context.Context) ([]*models.ThiefClassData, error) {
	thiefDataList, err := r.q.ListThiefClassData(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*models.ThiefClassData, len(thiefDataList))
	for i, data := range thiefDataList {
		result[i] = mapDbThiefDataToModel(data)
	}

	return result, nil
}

// GetNextThiefLevel gets the next thief level data after currentLevel
func (r *SQLCThiefDataRepository) GetNextThiefLevel(ctx context.Context, currentLevel int) (*models.ThiefClassData, error) {
	nextLevel, err := r.q.GetNextThiefLevel(ctx, int64(currentLevel))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return mapDbThiefDataToModel(nextLevel), nil
}

// mapDbThiefDataToModel converts DB thief data to model
func mapDbThiefDataToModel(data sqlcdb.ThiefClassDatum) *models.ThiefClassData {
	return &models.ThiefClassData{
		Level:            int(data.Level),
		ExperiencePoints: int(data.ExperiencePoints),
		HitDice:          data.HitDice,
		SavingThrow:      int(data.SavingThrow),
		FightingAbility:  int(data.FightingAbility),
	}
}
