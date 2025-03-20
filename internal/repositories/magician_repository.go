package repositories

import (
	"context"
	"database/sql"
	"errors"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

// MagicianRepository defines the interface for magician data operations
type MagicianRepository interface {
	// Get magician data by level
	GetMagicianClassData(ctx context.Context, level int) (*models.MagicianClassData, error)
	// Get the next level's data
	GetNextMagicianLevel(ctx context.Context, currentLevel int) (*models.MagicianClassData, error)
	// List all magician level data
	ListMagicianClassData(ctx context.Context) ([]*models.MagicianClassData, error)
}

// SQLCMagicianRepository implements MagicianRepository with SQLC
type SQLCMagicianRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

// NewSQLCMagicianRepository creates a new magician repository
func NewSQLCMagicianRepository(db *sql.DB) *SQLCMagicianRepository {
	return &SQLCMagicianRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

// GetMagicianClassData retrieves class data for the specified magician level
func (r *SQLCMagicianRepository) GetMagicianClassData(ctx context.Context, level int) (*models.MagicianClassData, error) {
	data, err := r.q.GetMagicianClassData(ctx, int64(level))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("magician class data", level)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	return &models.MagicianClassData{
		Level:            int(data.Level),
		ExperiencePoints: int(data.ExperiencePoints),
		HitDice:          data.HitDice,
		SavingThrow:      int(data.SavingThrow),
		FightingAbility:  int(data.FightingAbility),
		CastingAbility:   int(data.CastingAbility),
		SpellSlotsLevel1: int(data.SpellSlotsLevel1),
		SpellSlotsLevel2: int(data.SpellSlotsLevel2),
		SpellSlotsLevel3: int(data.SpellSlotsLevel3),
		SpellSlotsLevel4: int(data.SpellSlotsLevel4),
		SpellSlotsLevel5: int(data.SpellSlotsLevel5),
		SpellSlotsLevel6: int(data.SpellSlotsLevel6),
	}, nil
}

// GetNextMagicianLevel gets the next level data for Magicians
func (r *SQLCMagicianRepository) GetNextMagicianLevel(ctx context.Context, currentLevel int) (*models.MagicianClassData, error) {
	data, err := r.q.GetNextMagicianLevel(ctx, int64(currentLevel))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("next magician level", currentLevel+1)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	return &models.MagicianClassData{
		Level:            int(data.Level),
		ExperiencePoints: int(data.ExperiencePoints),
		HitDice:          data.HitDice,
		SavingThrow:      int(data.SavingThrow),
		FightingAbility:  int(data.FightingAbility),
		CastingAbility:   int(data.CastingAbility),
		SpellSlotsLevel1: int(data.SpellSlotsLevel1),
		SpellSlotsLevel2: int(data.SpellSlotsLevel2),
		SpellSlotsLevel3: int(data.SpellSlotsLevel3),
		SpellSlotsLevel4: int(data.SpellSlotsLevel4),
		SpellSlotsLevel5: int(data.SpellSlotsLevel5),
		SpellSlotsLevel6: int(data.SpellSlotsLevel6),
	}, nil
}

// ListMagicianClassData retrieves all magician level data
func (r *SQLCMagicianRepository) ListMagicianClassData(ctx context.Context) ([]*models.MagicianClassData, error) {
	dataList, err := r.q.ListMagicianClassData(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]*models.MagicianClassData, len(dataList))
	for i, data := range dataList {
		result[i] = &models.MagicianClassData{
			Level:            int(data.Level),
			ExperiencePoints: int(data.ExperiencePoints),
			HitDice:          data.HitDice,
			SavingThrow:      int(data.SavingThrow),
			FightingAbility:  int(data.FightingAbility),
			CastingAbility:   int(data.CastingAbility),
			SpellSlotsLevel1: int(data.SpellSlotsLevel1),
			SpellSlotsLevel2: int(data.SpellSlotsLevel2),
			SpellSlotsLevel3: int(data.SpellSlotsLevel3),
			SpellSlotsLevel4: int(data.SpellSlotsLevel4),
			SpellSlotsLevel5: int(data.SpellSlotsLevel5),
			SpellSlotsLevel6: int(data.SpellSlotsLevel6),
		}
	}

	return result, nil
}
