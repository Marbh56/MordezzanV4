package repositories

import (
	"context"
	"database/sql"
	"errors"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

// FighterDataRepository defines methods for accessing fighter class data
type FighterDataRepository interface {
	GetFighterClassData(ctx context.Context, level int) (*models.FighterClassData, error)
	ListFighterClassData(ctx context.Context) ([]*models.FighterClassData, error)
	GetNextFighterLevel(ctx context.Context, currentLevel int) (*models.FighterClassData, error)
}

// SQLCFighterDataRepository implements FighterDataRepository using SQLC
type SQLCFighterDataRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

// NewSQLCFighterDataRepository creates a new fighter data repository
func NewSQLCFighterDataRepository(db *sql.DB) *SQLCFighterDataRepository {
	return &SQLCFighterDataRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

// GetFighterClassData retrieves fighter class data for a specific level
func (r *SQLCFighterDataRepository) GetFighterClassData(ctx context.Context, level int) (*models.FighterClassData, error) {
	fighterData, err := r.q.GetFighterClassData(ctx, int64(level))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return mapDbFighterDataToModel(fighterData), nil
}

// ListFighterClassData retrieves all fighter class data levels
func (r *SQLCFighterDataRepository) ListFighterClassData(ctx context.Context) ([]*models.FighterClassData, error) {
	fighterDataList, err := r.q.ListFighterClassData(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*models.FighterClassData, len(fighterDataList))
	for i, data := range fighterDataList {
		result[i] = mapDbFighterDataToModel(data)
	}

	return result, nil
}

// GetNextFighterLevel gets the next fighter level data after currentLevel
func (r *SQLCFighterDataRepository) GetNextFighterLevel(ctx context.Context, currentLevel int) (*models.FighterClassData, error) {
	nextLevel, err := r.q.GetNextFighterLevel(ctx, int64(currentLevel))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return mapDbFighterDataToModel(nextLevel), nil
}

// mapDbFighterDataToModel converts DB fighter data to model
func mapDbFighterDataToModel(data sqlcdb.FighterClassDatum) *models.FighterClassData {
	return &models.FighterClassData{
		Level:            int(data.Level),
		ExperiencePoints: int(data.ExperiencePoints),
		HitDice:          data.HitDice,
		SavingThrow:      int(data.SavingThrow),
		FightingAbility:  int(data.FightingAbility),
	}
}
