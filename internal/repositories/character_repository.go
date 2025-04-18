package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type CharacterRepository interface {
	GetCharacter(ctx context.Context, id int64) (*models.Character, error)
	GetCharactersByUser(ctx context.Context, userID int64) ([]*models.Character, error)
	ListCharacters(ctx context.Context) ([]*models.Character, error)
	CreateCharacter(ctx context.Context, input *models.CreateCharacterInput) (int64, error)
	UpdateCharacter(ctx context.Context, id int64, input *models.UpdateCharacterInput) error
	DeleteCharacter(ctx context.Context, id int64) error
}

type SQLCCharacterRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCCharacterRepository(db *sql.DB) *SQLCCharacterRepository {
	return &SQLCCharacterRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCCharacterRepository) GetCharacter(ctx context.Context, id int64) (*models.Character, error) {
	dbCharacter, err := r.q.GetCharacter(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("character", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	character := mapDbCharacterRowToModel(dbCharacter)
	character.CalculateDerivedStats()

	return character, nil
}

func (r *SQLCCharacterRepository) GetCharactersByUser(ctx context.Context, userID int64) ([]*models.Character, error) {
	characters, err := r.q.GetCharactersByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Character, len(characters))
	for i, c := range characters {
		result[i] = mapDbCharactersByUserRowToModel(c)
	}

	return result, nil
}

func (r *SQLCCharacterRepository) ListCharacters(ctx context.Context) ([]*models.Character, error) {
	characters, err := r.q.ListCharacters(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Character, len(characters))
	for i, c := range characters {
		result[i] = mapDbListCharactersRowToModel(c)
	}

	return result, nil
}

func (r *SQLCCharacterRepository) CreateCharacter(ctx context.Context, input *models.CreateCharacterInput) (int64, error) {
	result, err := r.q.CreateCharacter(ctx, sqlcdb.CreateCharacterParams{
		UserID:             input.UserID,
		Name:               input.Name,
		Class:              input.Class,
		Level:              int64(input.Level),
		Strength:           int64(input.Strength),
		Dexterity:          int64(input.Dexterity),
		Constitution:       int64(input.Constitution),
		Wisdom:             int64(input.Wisdom),
		Intelligence:       int64(input.Intelligence),
		Charisma:           int64(input.Charisma),
		MaxHitPoints:       int64(input.MaxHitPoints),
		CurrentHitPoints:   int64(input.CurrentHitPoints),
		TemporaryHitPoints: int64(input.TemporaryHitPoints),
		ExperiencePoints:   int64(input.ExperiencePoints),
	})
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SQLCCharacterRepository) UpdateCharacter(ctx context.Context, id int64, input *models.UpdateCharacterInput) error {
	_, err := r.q.UpdateCharacter(ctx, sqlcdb.UpdateCharacterParams{
		Name:               input.Name,
		Class:              input.Class,
		Level:              int64(input.Level),
		Strength:           int64(input.Strength),
		Dexterity:          int64(input.Dexterity),
		Constitution:       int64(input.Constitution),
		Wisdom:             int64(input.Wisdom),
		Intelligence:       int64(input.Intelligence),
		Charisma:           int64(input.Charisma),
		MaxHitPoints:       int64(input.MaxHitPoints),
		CurrentHitPoints:   int64(input.CurrentHitPoints),
		TemporaryHitPoints: int64(input.TemporaryHitPoints),
		ExperiencePoints:   int64(input.ExperiencePoints),
		ID:                 id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.NewNotFound("character", id)
		}
		return err
	}

	return nil
}

func (r *SQLCCharacterRepository) DeleteCharacter(ctx context.Context, id int64) error {
	_, err := r.q.DeleteCharacter(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.NewNotFound("character", id)
		}
		return err
	}

	return nil
}

// This function is kept for reference but isn't used directly anymore
func mapDbCharacterToModel(character sqlcdb.Character) *models.Character {
	return &models.Character{
		ID:                 character.ID,
		UserID:             character.UserID,
		Name:               character.Name,
		Class:              character.Class,
		Level:              int(character.Level),
		Strength:           int(character.Strength),
		Dexterity:          int(character.Dexterity),
		Constitution:       int(character.Constitution),
		Wisdom:             int(character.Wisdom),
		Intelligence:       int(character.Intelligence),
		Charisma:           int(character.Charisma),
		MaxHitPoints:       int(character.MaxHitPoints),
		CurrentHitPoints:   int(character.CurrentHitPoints),
		TemporaryHitPoints: int(character.TemporaryHitPoints),
		ExperiencePoints:   int(character.ExperiencePoints),
		CreatedAt:          character.CreatedAt,
		UpdatedAt:          character.UpdatedAt,
	}
}

// Maps GetCharacterRow to models.Character
func mapDbCharacterRowToModel(character sqlcdb.GetCharacterRow) *models.Character {
	return &models.Character{
		ID:                 character.ID,
		UserID:             character.UserID,
		Name:               character.Name,
		Class:              character.Class,
		Level:              int(character.Level),
		Strength:           int(character.Strength),
		Dexterity:          int(character.Dexterity),
		Constitution:       int(character.Constitution),
		Wisdom:             int(character.Wisdom),
		Intelligence:       int(character.Intelligence),
		Charisma:           int(character.Charisma),
		MaxHitPoints:       int(character.MaxHitPoints),
		CurrentHitPoints:   int(character.CurrentHitPoints),
		TemporaryHitPoints: int(character.TemporaryHitPoints),
		ExperiencePoints:   int(character.ExperiencePoints),
		CreatedAt:          character.CreatedAt,
		UpdatedAt:          character.UpdatedAt,
	}
}

// Maps GetCharactersByUserRow to models.Character
func mapDbCharactersByUserRowToModel(character sqlcdb.GetCharactersByUserRow) *models.Character {
	return &models.Character{
		ID:                 character.ID,
		UserID:             character.UserID,
		Name:               character.Name,
		Class:              character.Class,
		Level:              int(character.Level),
		Strength:           int(character.Strength),
		Dexterity:          int(character.Dexterity),
		Constitution:       int(character.Constitution),
		Wisdom:             int(character.Wisdom),
		Intelligence:       int(character.Intelligence),
		Charisma:           int(character.Charisma),
		MaxHitPoints:       int(character.MaxHitPoints),
		CurrentHitPoints:   int(character.CurrentHitPoints),
		TemporaryHitPoints: int(character.TemporaryHitPoints),
		ExperiencePoints:   int(character.ExperiencePoints),
		CreatedAt:          character.CreatedAt,
		UpdatedAt:          character.UpdatedAt,
	}
}

// Maps ListCharactersRow to models.Character
func mapDbListCharactersRowToModel(character sqlcdb.ListCharactersRow) *models.Character {
	return &models.Character{
		ID:                 character.ID,
		UserID:             character.UserID,
		Name:               character.Name,
		Class:              character.Class,
		Level:              int(character.Level),
		Strength:           int(character.Strength),
		Dexterity:          int(character.Dexterity),
		Constitution:       int(character.Constitution),
		Wisdom:             int(character.Wisdom),
		Intelligence:       int(character.Intelligence),
		Charisma:           int(character.Charisma),
		MaxHitPoints:       int(character.MaxHitPoints),
		CurrentHitPoints:   int(character.CurrentHitPoints),
		TemporaryHitPoints: int(character.TemporaryHitPoints),
		ExperiencePoints:   int(character.ExperiencePoints),
		CreatedAt:          character.CreatedAt,
		UpdatedAt:          character.UpdatedAt,
	}
}
