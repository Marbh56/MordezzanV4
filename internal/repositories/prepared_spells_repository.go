package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

// PreparedSpellRepository defines operations for prepared spells
type PreparedSpellRepository interface {
	GetPreparedSpell(ctx context.Context, id int64) (*models.PreparedSpell, error)
	GetPreparedSpellsByCharacter(ctx context.Context, characterID int64) ([]*models.PreparedSpell, error)
	IsSpellPrepared(ctx context.Context, characterID int64, spellID int64) (bool, error)
	CountPreparedSpellsByLevel(ctx context.Context, characterID int64, slotLevel int) (int, error)
	PrepareSpell(ctx context.Context, characterID int64, spellID int64, slotLevel int) error
	UnprepareSpell(ctx context.Context, characterID int64, spellID int64) error
	ClearPreparedSpells(ctx context.Context, characterID int64) error
}

// SQLCPreparedSpellRepository implements PreparedSpellRepository using SQLC
type SQLCPreparedSpellRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

// NewSQLCPreparedSpellRepository creates a new prepared spell repository
func NewSQLCPreparedSpellRepository(db *sql.DB) *SQLCPreparedSpellRepository {
	return &SQLCPreparedSpellRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

// GetPreparedSpell gets a prepared spell by ID
func (r *SQLCPreparedSpellRepository) GetPreparedSpell(ctx context.Context, id int64) (*models.PreparedSpell, error) {
	preparedSpell, err := r.q.GetPreparedSpell(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("prepared spell", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	return &models.PreparedSpell{
		ID:          preparedSpell.ID,
		CharacterID: preparedSpell.CharacterID,
		SpellID:     preparedSpell.SpellID,
		SlotLevel:   int(preparedSpell.SlotLevel),
		PreparedAt:  preparedSpell.PreparedAt,
	}, nil
}

// GetPreparedSpellsByCharacter gets all prepared spells for a character
func (r *SQLCPreparedSpellRepository) GetPreparedSpellsByCharacter(ctx context.Context, characterID int64) ([]*models.PreparedSpell, error) {
	preparedSpells, err := r.q.GetPreparedSpellsByCharacter(ctx, characterID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]*models.PreparedSpell, len(preparedSpells))
	for i, ps := range preparedSpells {
		result[i] = &models.PreparedSpell{
			ID:          ps.ID,
			CharacterID: ps.CharacterID,
			SpellID:     ps.SpellID,
			SlotLevel:   int(ps.SlotLevel),
			PreparedAt:  ps.PreparedAt,
		}
	}

	return result, nil
}

// IsSpellPrepared checks if a spell is already prepared by the character
func (r *SQLCPreparedSpellRepository) IsSpellPrepared(ctx context.Context, characterID int64, spellID int64) (bool, error) {
	count, err := r.q.CountPreparedSpell(ctx, sqlcdb.CountPreparedSpellParams{
		CharacterID: characterID,
		SpellID:     spellID,
	})
	if err != nil {
		return false, apperrors.NewDatabaseError(err)
	}
	return count > 0, nil
}

// CountPreparedSpellsByLevel counts how many spells a character has prepared at a specific level
func (r *SQLCPreparedSpellRepository) CountPreparedSpellsByLevel(ctx context.Context, characterID int64, slotLevel int) (int, error) {
	count, err := r.q.CountPreparedSpellsByLevel(ctx, sqlcdb.CountPreparedSpellsByLevelParams{
		CharacterID: characterID,
		SlotLevel:   int64(slotLevel),
	})
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}
	return int(count), nil
}

// PrepareSpell prepares a spell for a character
func (r *SQLCPreparedSpellRepository) PrepareSpell(ctx context.Context, characterID int64, spellID int64, slotLevel int) error {
	_, err := r.q.PrepareSpell(ctx, sqlcdb.PrepareSpellParams{
		CharacterID: characterID,
		SpellID:     spellID,
		SlotLevel:   int64(slotLevel),
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

// UnprepareSpell removes a prepared spell
func (r *SQLCPreparedSpellRepository) UnprepareSpell(ctx context.Context, characterID int64, spellID int64) error {
	err := r.q.UnprepareSpell(ctx, sqlcdb.UnprepareSpellParams{
		CharacterID: characterID,
		SpellID:     spellID,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

// ClearPreparedSpells removes all prepared spells for a character
func (r *SQLCPreparedSpellRepository) ClearPreparedSpells(ctx context.Context, characterID int64) error {
	err := r.q.ClearPreparedSpells(ctx, characterID)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}
