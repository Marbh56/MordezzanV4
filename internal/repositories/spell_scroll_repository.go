package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type SpellScrollRepository interface {
	GetSpellScroll(ctx context.Context, id int64) (*models.SpellScroll, error)
	ListSpellScrolls(ctx context.Context) ([]*models.SpellScroll, error)
	GetSpellScrollsBySpell(ctx context.Context, spellID int64) ([]*models.SpellScroll, error)
	CreateSpellScroll(ctx context.Context, input *models.CreateSpellScrollInput) (int64, error)
	UpdateSpellScroll(ctx context.Context, id int64, input *models.UpdateSpellScrollInput) error
	DeleteSpellScroll(ctx context.Context, id int64) error
}

type SQLCSpellScrollRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCSpellScrollRepository(db *sql.DB) *SQLCSpellScrollRepository {
	return &SQLCSpellScrollRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCSpellScrollRepository) GetSpellScroll(ctx context.Context, id int64) (*models.SpellScroll, error) {
	spellScroll, err := r.q.GetSpellScroll(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("spell scroll", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbSpellScrollToModel(spellScroll), nil
}

func (r *SQLCSpellScrollRepository) ListSpellScrolls(ctx context.Context) ([]*models.SpellScroll, error) {
	spellScrolls, err := r.q.ListSpellScrolls(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.SpellScroll, len(spellScrolls))
	for i, spellScroll := range spellScrolls {
		result[i] = mapDbListSpellScrollToModel(spellScroll)
	}
	return result, nil
}

func (r *SQLCSpellScrollRepository) GetSpellScrollsBySpell(ctx context.Context, spellID int64) ([]*models.SpellScroll, error) {
	spellScrolls, err := r.q.GetSpellScrollsBySpell(ctx, spellID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.SpellScroll, len(spellScrolls))
	for i, spellScroll := range spellScrolls {
		result[i] = mapDbSpellScrollsBySpellToModel(spellScroll)
	}
	return result, nil
}

func (r *SQLCSpellScrollRepository) CreateSpellScroll(ctx context.Context, input *models.CreateSpellScrollInput) (int64, error) {
	result, err := r.q.CreateSpellScroll(ctx, sqlcdb.CreateSpellScrollParams{
		SpellID:      input.SpellID,
		CastingLevel: int64(input.CastingLevel),
		Cost:         input.Cost,
		Weight:       int64(input.Weight),
		Description:  sql.NullString{String: input.Description, Valid: input.Description != ""},
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

func (r *SQLCSpellScrollRepository) UpdateSpellScroll(ctx context.Context, id int64, input *models.UpdateSpellScrollInput) error {
	_, err := r.GetSpellScroll(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.UpdateSpellScroll(ctx, sqlcdb.UpdateSpellScrollParams{
		ID:           id,
		SpellID:      input.SpellID,
		CastingLevel: int64(input.CastingLevel),
		Cost:         input.Cost,
		Weight:       int64(input.Weight),
		Description:  sql.NullString{String: input.Description, Valid: input.Description != ""},
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCSpellScrollRepository) DeleteSpellScroll(ctx context.Context, id int64) error {
	_, err := r.GetSpellScroll(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeleteSpellScroll(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func mapDbSpellScrollToModel(spellScroll sqlcdb.GetSpellScrollRow) *models.SpellScroll {
	return &models.SpellScroll{
		ID:           spellScroll.ID,
		SpellID:      spellScroll.SpellID,
		SpellName:    spellScroll.SpellName,
		CastingLevel: int(spellScroll.CastingLevel),
		Cost:         spellScroll.Cost,
		Weight:       int(spellScroll.Weight),
		Description:  spellScroll.Description.String,
		CreatedAt:    spellScroll.CreatedAt,
		UpdatedAt:    spellScroll.UpdatedAt,
	}
}

func mapDbListSpellScrollToModel(spellScroll sqlcdb.ListSpellScrollsRow) *models.SpellScroll {
	return &models.SpellScroll{
		ID:           spellScroll.ID,
		SpellID:      spellScroll.SpellID,
		SpellName:    spellScroll.SpellName,
		CastingLevel: int(spellScroll.CastingLevel),
		Cost:         spellScroll.Cost,
		Weight:       int(spellScroll.Weight),
		Description:  spellScroll.Description.String,
		CreatedAt:    spellScroll.CreatedAt,
		UpdatedAt:    spellScroll.UpdatedAt,
	}
}

func mapDbSpellScrollsBySpellToModel(spellScroll sqlcdb.GetSpellScrollsBySpellRow) *models.SpellScroll {
	return &models.SpellScroll{
		ID:           spellScroll.ID,
		SpellID:      spellScroll.SpellID,
		SpellName:    spellScroll.SpellName,
		CastingLevel: int(spellScroll.CastingLevel),
		Cost:         spellScroll.Cost,
		Weight:       int(spellScroll.Weight),
		Description:  spellScroll.Description.String,
		CreatedAt:    spellScroll.CreatedAt,
		UpdatedAt:    spellScroll.UpdatedAt,
	}
}
