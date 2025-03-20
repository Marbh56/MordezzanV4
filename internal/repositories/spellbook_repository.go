package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"

	apperrors "mordezzanV4/internal/errors"
)

type SpellbookRepository interface {
	GetSpellbook(ctx context.Context, id int64) (*models.Spellbook, error)
	GetSpellbookByName(ctx context.Context, name string) (*models.Spellbook, error)
	ListSpellbooks(ctx context.Context) ([]*models.Spellbook, error)
	CreateSpellbook(ctx context.Context, input *models.CreateSpellbookInput) (int64, error)
	UpdateSpellbook(ctx context.Context, id int64, input *models.UpdateSpellbookInput) error
	DeleteSpellbook(ctx context.Context, id int64) error

	// Spell management
	AddSpellToSpellbook(ctx context.Context, spellbookID, spellID int64, characterClass string) error
	RemoveSpellFromSpellbook(ctx context.Context, spellbookID, spellID int64) error
	GetSpellsInSpellbook(ctx context.Context, spellbookID int64) ([]int64, error)
}

type SQLCSpellbookRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCSpellbookRepository(db *sql.DB) *SQLCSpellbookRepository {
	return &SQLCSpellbookRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCSpellbookRepository) GetSpellbook(ctx context.Context, id int64) (*models.Spellbook, error) {
	spellbook, err := r.q.GetSpellbook(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("spellbook", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	// Get spells in the spellbook
	spellIDs, err := r.GetSpellsInSpellbook(ctx, id)
	if err != nil {
		return nil, err
	}

	description := ""
	if spellbook.Description.Valid {
		description = spellbook.Description.String
	}

	return &models.Spellbook{
		ID:           spellbook.ID,
		Name:         spellbook.Name,
		Description:  description,
		TotalPages:   int(spellbook.TotalPages),
		UsedPages:    int(spellbook.UsedPages),
		Value:        int(spellbook.Value),
		Weight:       spellbook.Weight,
		SpellsStored: spellIDs,
		CreatedAt:    spellbook.CreatedAt,
		UpdatedAt:    spellbook.UpdatedAt,
	}, nil
}

func (r *SQLCSpellbookRepository) GetSpellbookByName(ctx context.Context, name string) (*models.Spellbook, error) {
	spellbook, err := r.q.GetSpellbookByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No error, just not found
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	// Get spells in the spellbook
	spellIDs, err := r.GetSpellsInSpellbook(ctx, spellbook.ID)
	if err != nil {
		return nil, err
	}

	description := ""
	if spellbook.Description.Valid {
		description = spellbook.Description.String
	}

	return &models.Spellbook{
		ID:           spellbook.ID,
		Name:         spellbook.Name,
		Description:  description,
		TotalPages:   int(spellbook.TotalPages),
		UsedPages:    int(spellbook.UsedPages),
		Value:        int(spellbook.Value),
		Weight:       spellbook.Weight,
		SpellsStored: spellIDs,
		CreatedAt:    spellbook.CreatedAt,
		UpdatedAt:    spellbook.UpdatedAt,
	}, nil
}

func (r *SQLCSpellbookRepository) ListSpellbooks(ctx context.Context) ([]*models.Spellbook, error) {
	spellbooks, err := r.q.ListSpellbooks(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]*models.Spellbook, len(spellbooks))
	for i, spellbook := range spellbooks {
		// Get spells in the spellbook
		spellIDs, err := r.GetSpellsInSpellbook(ctx, spellbook.ID)
		if err != nil {
			return nil, err
		}

		description := ""
		if spellbook.Description.Valid {
			description = spellbook.Description.String
		}

		result[i] = &models.Spellbook{
			ID:           spellbook.ID,
			Name:         spellbook.Name,
			Description:  description,
			TotalPages:   int(spellbook.TotalPages),
			UsedPages:    int(spellbook.UsedPages),
			Value:        int(spellbook.Value),
			Weight:       spellbook.Weight,
			SpellsStored: spellIDs,
			CreatedAt:    spellbook.CreatedAt,
			UpdatedAt:    spellbook.UpdatedAt,
		}
	}

	return result, nil
}

func (r *SQLCSpellbookRepository) CreateSpellbook(ctx context.Context, input *models.CreateSpellbookInput) (int64, error) {
	// Check if spellbook with same name already exists
	existingSpellbook, err := r.GetSpellbookByName(ctx, input.Name)
	if err != nil {
		return 0, err
	}
	if existingSpellbook != nil {
		return 0, apperrors.NewValidationError("name", "Spellbook with this name already exists")
	}

	result, err := r.q.CreateSpellbook(ctx, sqlcdb.CreateSpellbookParams{
		Name: input.Name,
		Description: sql.NullString{
			String: input.Description,
			Valid:  input.Description != "",
		},
		TotalPages: int64(input.TotalPages),
		Value:      int64(input.Value),
		Weight:     input.Weight,
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

func (r *SQLCSpellbookRepository) UpdateSpellbook(ctx context.Context, id int64, input *models.UpdateSpellbookInput) error {
	// Check if spellbook exists
	_, err := r.GetSpellbook(ctx, id)
	if err != nil {
		return err
	}

	// Check if updating to a name that's already taken
	existingSpellbook, err := r.GetSpellbookByName(ctx, input.Name)
	if err != nil {
		return err
	}
	if existingSpellbook != nil && existingSpellbook.ID != id {
		return apperrors.NewValidationError("name", "Spellbook with this name already exists")
	}

	err = r.q.UpdateSpellbook(ctx, sqlcdb.UpdateSpellbookParams{
		ID:   id,
		Name: input.Name,
		Description: sql.NullString{
			String: input.Description,
			Valid:  input.Description != "",
		},
		TotalPages: int64(input.TotalPages),
		UsedPages:  int64(input.UsedPages),
		Value:      int64(input.Value),
		Weight:     input.Weight,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

func (r *SQLCSpellbookRepository) DeleteSpellbook(ctx context.Context, id int64) error {
	// Check if spellbook exists
	_, err := r.GetSpellbook(ctx, id)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// Delete all spell associations
	err = qtx.DeleteAllSpellsFromSpellbook(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	// Delete the spellbook
	err = qtx.DeleteSpellbook(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	if err := tx.Commit(); err != nil {
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

func (r *SQLCSpellbookRepository) AddSpellToSpellbook(ctx context.Context, spellbookID, spellID int64, characterClass string) error {
	// Get the spellbook
	spellbook, err := r.GetSpellbook(ctx, spellbookID)
	if err != nil {
		return err
	}

	// Get the spell
	spellRepo := NewSQLCSpellRepository(r.db)
	spell, err := spellRepo.GetSpell(ctx, spellID)
	if err != nil {
		return err
	}

	// Determine pages required based on character class
	var pagesRequired int

	switch characterClass {
	case "Magician":
		pagesRequired = spell.MagLevel
	case "Cryo-mancer":
		pagesRequired = spell.CryLevel
	case "Illusionist":
		pagesRequired = spell.IllLevel
	case "Necromancer":
		pagesRequired = spell.NecLevel
	case "Pyromancer":
		pagesRequired = spell.PyrLevel
	case "Witch":
		pagesRequired = spell.WchLevel
	case "Cleric":
		pagesRequired = spell.ClrLevel
	case "Druid":
		pagesRequired = spell.DrdLevel
	default:
		return apperrors.NewBadRequest(fmt.Sprintf("Unsupported character class: %s", characterClass))
	}

	// Check if the spell is available for this class
	if pagesRequired == 0 {
		return apperrors.NewBadRequest(fmt.Sprintf("Spell '%s' is not available for %s class", spell.Name, characterClass))
	}

	// Check if there's enough space in the spellbook
	if spellbook.UsedPages+pagesRequired > spellbook.TotalPages {
		return apperrors.NewBadRequest(fmt.Sprintf(
			"Not enough pages in spellbook. Requires %d pages, only %d available",
			pagesRequired, spellbook.TotalPages-spellbook.UsedPages))
	}

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// Add the spell to the spellbook
	err = qtx.AddSpellToSpellbook(ctx, sqlcdb.AddSpellToSpellbookParams{
		SpellbookID:    spellbookID,
		SpellID:        spellID,
		CharacterClass: characterClass,
		PagesUsed:      int64(pagesRequired),
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	// Update the used pages count
	err = qtx.UpdateSpellbookUsedPages(ctx, sqlcdb.UpdateSpellbookUsedPagesParams{
		ID:        spellbookID,
		UsedPages: int64(spellbook.UsedPages + pagesRequired),
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	if err := tx.Commit(); err != nil {
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

func (r *SQLCSpellbookRepository) RemoveSpellFromSpellbook(ctx context.Context, spellbookID, spellID int64) error {
	// Get the spellbook
	spellbook, err := r.GetSpellbook(ctx, spellbookID)
	if err != nil {
		return err
	}

	// Get the spell entry from the spellbook to determine pages used
	spellEntry, err := r.q.GetSpellFromSpellbook(ctx, sqlcdb.GetSpellFromSpellbookParams{
		SpellbookID: spellbookID,
		SpellID:     spellID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.NewNotFound("spell in spellbook", spellID)
		}
		return apperrors.NewDatabaseError(err)
	}

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// Remove the spell from the spellbook
	err = qtx.RemoveSpellFromSpellbook(ctx, sqlcdb.RemoveSpellFromSpellbookParams{
		SpellbookID: spellbookID,
		SpellID:     spellID,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	// Update the used pages count
	newUsedPages := spellbook.UsedPages - int(spellEntry.PagesUsed)
	if newUsedPages < 0 {
		newUsedPages = 0 // Safeguard
	}

	err = qtx.UpdateSpellbookUsedPages(ctx, sqlcdb.UpdateSpellbookUsedPagesParams{
		ID:        spellbookID,
		UsedPages: int64(newUsedPages),
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	if err := tx.Commit(); err != nil {
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

func (r *SQLCSpellbookRepository) GetSpellsInSpellbook(ctx context.Context, spellbookID int64) ([]int64, error) {
	spellIDs, err := r.q.GetSpellsInSpellbook(ctx, spellbookID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	return spellIDs, nil
}
