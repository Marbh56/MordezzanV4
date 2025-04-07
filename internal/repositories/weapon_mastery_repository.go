package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
	"strings"
)

type WeaponMasteryRepository interface {
	GetWeaponMasteriesByCharacter(ctx context.Context, characterID int64) ([]*models.WeaponMastery, error)
	AddWeaponMastery(ctx context.Context, input *models.AddWeaponMasteryInput) (int64, error)
	UpdateWeaponMastery(ctx context.Context, characterID int64, weaponBaseName string, input *models.UpdateWeaponMasteryInput) error
	DeleteWeaponMastery(ctx context.Context, characterID int64, weaponBaseName string) error
	CountWeaponMasteries(ctx context.Context, characterID int64, masteryLevel string) (int, error)
	GetWeaponMasteryByID(ctx context.Context, id int64) (*models.WeaponMastery, error)
	GetWeaponMasteryByBaseName(ctx context.Context, characterID int64, weaponBaseName string) (*models.WeaponMastery, error)
	GetWeaponBaseNameFromID(ctx context.Context, weaponID int64) (string, error)
}

type SQLCWeaponMasteryRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCWeaponMasteryRepository(db *sql.DB) *SQLCWeaponMasteryRepository {
	return &SQLCWeaponMasteryRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCWeaponMasteryRepository) GetWeaponMasteriesByCharacter(ctx context.Context, characterID int64) ([]*models.WeaponMastery, error) {
	masteriesRows, err := r.q.GetWeaponMasteriesByCharacter(ctx, characterID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]*models.WeaponMastery, len(masteriesRows))
	for i, row := range masteriesRows {
		result[i] = &models.WeaponMastery{
			ID:             row.ID,
			CharacterID:    row.CharacterID,
			WeaponBaseName: row.WeaponBaseName,
			MasteryLevel:   row.MasteryLevel,
			CreatedAt:      row.CreatedAt,
			UpdatedAt:      row.UpdatedAt,
		}
	}
	return result, nil
}

func (r *SQLCWeaponMasteryRepository) AddWeaponMastery(ctx context.Context, input *models.AddWeaponMasteryInput) (int64, error) {
	// Check if the character exists
	_, err := r.q.GetCharacter(ctx, input.CharacterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperrors.NewNotFound("character", input.CharacterID)
		}
		return 0, apperrors.NewDatabaseError(err)
	}

	// Check if there's already a mastery for this weapon base name
	mastery, err := r.GetWeaponMasteryByBaseName(ctx, input.CharacterID, input.WeaponBaseName)
	if err == nil && mastery != nil {
		return 0, apperrors.NewBadRequest("Character already has mastery for this weapon type")
	}

	// If it's grand mastery, check if they already have one
	if input.MasteryLevel == "grand_mastery" {
		count, err := r.CountWeaponMasteries(ctx, input.CharacterID, "grand_mastery")
		if err != nil {
			return 0, err
		}
		if count >= 1 {
			return 0, apperrors.NewBadRequest("Character already has a grand mastery weapon")
		}
	}

	err = r.q.AddWeaponMastery(ctx, sqlcdb.AddWeaponMasteryParams{
		CharacterID:    input.CharacterID,
		WeaponBaseName: input.WeaponBaseName,
		MasteryLevel:   input.MasteryLevel,
	})
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	// Get the last inserted ID
	var id int64
	err = r.db.QueryRowContext(ctx, "SELECT last_insert_rowid()").Scan(&id)
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	return id, nil
}

func (r *SQLCWeaponMasteryRepository) UpdateWeaponMastery(ctx context.Context, characterID int64, weaponBaseName string, input *models.UpdateWeaponMasteryInput) error {
	// Check if the mastery exists
	mastery, err := r.GetWeaponMasteryByBaseName(ctx, characterID, weaponBaseName)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.NewNotFound("weapon mastery", 0)
		}
		return err
	}

	// If it's grand mastery, check if they already have one
	if input.MasteryLevel == "grand_mastery" && mastery.MasteryLevel != "grand_mastery" {
		// Count existing grand masteries
		count, err := r.CountWeaponMasteries(ctx, characterID, "grand_mastery")
		if err != nil {
			return err
		}
		if count >= 1 {
			return apperrors.NewBadRequest("Character already has a grand mastery weapon")
		}
	}

	err = r.q.UpdateWeaponMasteryLevel(ctx, sqlcdb.UpdateWeaponMasteryLevelParams{
		MasteryLevel:   input.MasteryLevel,
		CharacterID:    characterID,
		WeaponBaseName: weaponBaseName,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.NewNotFound("weapon mastery", 0)
		}
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

func (r *SQLCWeaponMasteryRepository) DeleteWeaponMastery(ctx context.Context, characterID int64, weaponBaseName string) error {
	err := r.q.DeleteWeaponMastery(ctx, sqlcdb.DeleteWeaponMasteryParams{
		CharacterID:    characterID,
		WeaponBaseName: weaponBaseName,
	})

	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCWeaponMasteryRepository) CountWeaponMasteries(ctx context.Context, characterID int64, masteryLevel string) (int, error) {
	count, err := r.q.CountWeaponMasteries(ctx, sqlcdb.CountWeaponMasteriesParams{
		CharacterID:  characterID,
		MasteryLevel: masteryLevel,
	})
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	return int(count), nil
}

func (r *SQLCWeaponMasteryRepository) GetWeaponMasteryByID(ctx context.Context, id int64) (*models.WeaponMastery, error) {
	mastery, err := r.q.GetWeaponMasteryByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("weapon mastery", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	return &models.WeaponMastery{
		ID:             mastery.ID,
		CharacterID:    mastery.CharacterID,
		WeaponBaseName: mastery.WeaponBaseName,
		MasteryLevel:   mastery.MasteryLevel,
		CreatedAt:      mastery.CreatedAt,
		UpdatedAt:      mastery.UpdatedAt,
	}, nil
}

func (r *SQLCWeaponMasteryRepository) GetWeaponMasteryByBaseName(ctx context.Context, characterID int64, weaponBaseName string) (*models.WeaponMastery, error) {
	mastery, err := r.q.GetWeaponMasteryByBaseName(ctx, sqlcdb.GetWeaponMasteryByBaseNameParams{
		CharacterID:    characterID,
		WeaponBaseName: weaponBaseName,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("weapon mastery", 0)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	return &models.WeaponMastery{
		ID:             mastery.ID,
		CharacterID:    mastery.CharacterID,
		WeaponBaseName: mastery.WeaponBaseName,
		MasteryLevel:   mastery.MasteryLevel,
		CreatedAt:      mastery.CreatedAt,
		UpdatedAt:      mastery.UpdatedAt,
	}, nil
}

func (r *SQLCWeaponMasteryRepository) GetWeaponBaseNameFromID(ctx context.Context, weaponID int64) (string, error) {
	weapon, err := r.q.GetWeapon(ctx, weaponID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", apperrors.NewNotFound("weapon", weaponID)
		}
		return "", apperrors.NewDatabaseError(err)
	}

	// Extract base name by removing any "+" and numbers after it
	baseName := weapon.Name
	if idx := strings.Index(baseName, " +"); idx != -1 {
		baseName = baseName[:idx]
	}

	// Also handle other magical affixes
	suffixes := []string{
		" of Slaying",
		" of Fire",
		" of Frost",
		" of Lightning",
		" of Venom",
		" of Speed",
		" of Accuracy",
		" of Power",
	}

	for _, suffix := range suffixes {
		if idx := strings.Index(baseName, suffix); idx != -1 {
			baseName = baseName[:idx]
			break
		}
	}

	return baseName, nil
}

// Additional SQLC query definitions needed:

// name: GetWeaponMasteryByBaseName :one
// SELECT id, character_id, weapon_base_name, mastery_level, created_at, updated_at
// FROM weapon_masteries
// WHERE character_id = ? AND weapon_base_name = ?;
