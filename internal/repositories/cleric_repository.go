package repositories

import (
	"context"
	"database/sql"
	"errors"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

// ClericRepository defines methods for accessing cleric class data
type ClericRepository interface {
	GetClericClassData(ctx context.Context, level int) (*models.ClericClassData, error)
	ListClericClassData(ctx context.Context) ([]*models.ClericClassData, error)
	GetNextClericLevel(ctx context.Context, currentLevel int) (*models.ClericClassData, error)
}

// SQLCClericRepository implements ClericRepository using SQLC
type SQLCClericRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

// NewSQLCClericRepository creates a new cleric data repository
func NewSQLCClericRepository(db *sql.DB) *SQLCClericRepository {
	return &SQLCClericRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

// GetClericClassData retrieves cleric class data for a specific level
func (r *SQLCClericRepository) GetClericClassData(ctx context.Context, level int) (*models.ClericClassData, error) {
	clericData, err := r.q.GetClericClassData(ctx, int64(level))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return mapDbClericDataToModel(clericData), nil
}

// ListClericClassData retrieves all cleric class data levels
func (r *SQLCClericRepository) ListClericClassData(ctx context.Context) ([]*models.ClericClassData, error) {
	clericDataList, err := r.q.ListClericClassData(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*models.ClericClassData, len(clericDataList))
	for i, data := range clericDataList {
		result[i] = mapDbClericDataToModel(data)
	}

	return result, nil
}

// GetNextClericLevel gets the next cleric level data after currentLevel
func (r *SQLCClericRepository) GetNextClericLevel(ctx context.Context, currentLevel int) (*models.ClericClassData, error) {
	nextLevel, err := r.q.GetNextClericLevel(ctx, int64(currentLevel))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return mapDbClericDataToModel(nextLevel), nil
}

// mapDbClericDataToModel converts DB cleric data to model
func mapDbClericDataToModel(data sqlcdb.ClericClassDatum) *models.ClericClassData {
	return &models.ClericClassData{
		Level:            int(data.Level),
		ExperiencePoints: int(data.ExperiencePoints),
		HitDice:          data.HitDice,
		SavingThrow:      int(data.SavingThrow),
		FightingAbility:  int(data.FightingAbility),
		TurningAbility:   int(data.TurningAbility),
		// Use explicit values for spell slots until sqlc generation is updated
		SpellSlotsLevel1: 0,
		SpellSlotsLevel2: 0,
		SpellSlotsLevel3: 0,
		SpellSlotsLevel4: 0,
		SpellSlotsLevel5: 0,
	}
}
