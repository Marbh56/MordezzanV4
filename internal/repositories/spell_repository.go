package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

// SpellRepository defines the interface for spell data operations
type SpellRepository interface {
	GetSpell(ctx context.Context, id int64) (*models.Spell, error)
	GetSpellsByCharacter(ctx context.Context, characterID int64) ([]*models.Spell, error)
	ListSpells(ctx context.Context) ([]*models.Spell, error)
	CreateSpell(ctx context.Context, input *models.CreateSpellInput) (int64, error)
	UpdateSpell(ctx context.Context, id int64, input *models.UpdateSpellInput) error
	DeleteSpell(ctx context.Context, id int64) error
}

// SQLCSpellRepository implements SpellRepository using SQLC
type SQLCSpellRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

// NewSQLCSpellRepository creates a new SQLCSpellRepository
func NewSQLCSpellRepository(db *sql.DB) *SQLCSpellRepository {
	return &SQLCSpellRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

// GetSpell retrieves a spell by its ID
func (r *SQLCSpellRepository) GetSpell(ctx context.Context, id int64) (*models.Spell, error) {
	spell, err := r.q.GetSpell(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("spell", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbSpellToModel(spell), nil
}

// GetSpellsByCharacter retrieves all spells for a specific character
func (r *SQLCSpellRepository) GetSpellsByCharacter(ctx context.Context, characterID int64) ([]*models.Spell, error) {
	spells, err := r.q.GetSpellsByCharacter(ctx, characterID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.Spell, len(spells))
	for i, spell := range spells {
		result[i] = mapDbSpellToModel(spell)
	}
	return result, nil
}

// ListSpells retrieves all spells
func (r *SQLCSpellRepository) ListSpells(ctx context.Context) ([]*models.Spell, error) {
	spells, err := r.q.ListSpells(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.Spell, len(spells))
	for i, spell := range spells {
		result[i] = mapDbSpellToModel(spell)
	}
	return result, nil
}

// CreateSpell creates a new spell
func (r *SQLCSpellRepository) CreateSpell(ctx context.Context, input *models.CreateSpellInput) (int64, error) {
	result, err := r.q.CreateSpell(ctx, sqlcdb.CreateSpellParams{
		CharacterID:  input.CharacterID,
		Name:         input.Name,
		MagLevel:     int64(input.MagLevel),
		CryLevel:     int64(input.CryLevel),
		IllLevel:     int64(input.IllLevel),
		NecLevel:     int64(input.NecLevel),
		PyrLevel:     int64(input.PyrLevel),
		WchLevel:     int64(input.WchLevel),
		ClrLevel:     int64(input.ClrLevel),
		DrdLevel:     int64(input.DrdLevel),
		Range:        input.Range,
		Duration:     input.Duration,
		AreaOfEffect: sql.NullString{String: input.AreaOfEffect, Valid: input.AreaOfEffect != ""},
		Components:   sql.NullString{String: input.Components, Valid: input.Components != ""},
		Description:  input.Description,
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

// UpdateSpell updates an existing spell
func (r *SQLCSpellRepository) UpdateSpell(ctx context.Context, id int64, input *models.UpdateSpellInput) error {
	// First check if the spell exists
	_, err := r.GetSpell(ctx, id)
	if err != nil {
		return err
	}

	_, err = r.q.UpdateSpell(ctx, sqlcdb.UpdateSpellParams{
		Name:         input.Name,
		MagLevel:     int64(input.MagLevel),
		CryLevel:     int64(input.CryLevel),
		IllLevel:     int64(input.IllLevel),
		NecLevel:     int64(input.NecLevel),
		PyrLevel:     int64(input.PyrLevel),
		WchLevel:     int64(input.WchLevel),
		ClrLevel:     int64(input.ClrLevel),
		DrdLevel:     int64(input.DrdLevel),
		Range:        input.Range,
		Duration:     input.Duration,
		AreaOfEffect: sql.NullString{String: input.AreaOfEffect, Valid: input.AreaOfEffect != ""},
		Components:   sql.NullString{String: input.Components, Valid: input.Components != ""},
		Description:  input.Description,
		ID:           id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

// DeleteSpell deletes a spell by its ID
func (r *SQLCSpellRepository) DeleteSpell(ctx context.Context, id int64) error {
	// First check if the spell exists
	_, err := r.GetSpell(ctx, id)
	if err != nil {
		return err
	}

	_, err = r.q.DeleteSpell(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

// mapDbSpellToModel converts a database spell to a model spell
func mapDbSpellToModel(spell sqlcdb.Spell) *models.Spell {
	return &models.Spell{
		ID:           spell.ID,
		CharacterID:  spell.CharacterID,
		Name:         spell.Name,
		MagLevel:     int(spell.MagLevel),
		CryLevel:     int(spell.CryLevel),
		IllLevel:     int(spell.IllLevel),
		NecLevel:     int(spell.NecLevel),
		PyrLevel:     int(spell.PyrLevel),
		WchLevel:     int(spell.WchLevel),
		ClrLevel:     int(spell.ClrLevel),
		DrdLevel:     int(spell.DrdLevel),
		Range:        spell.Range,
		Duration:     spell.Duration,
		AreaOfEffect: spell.AreaOfEffect.String,
		Components:   spell.Components.String,
		Description:  spell.Description,
		CreatedAt:    spell.CreatedAt,
		UpdatedAt:    spell.UpdatedAt,
	}
}
