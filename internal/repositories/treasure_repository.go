package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type TreasureRepository interface {
	GetTreasure(ctx context.Context, id int64) (*models.Treasure, error)
	GetTreasureByCharacter(ctx context.Context, characterID int64) (*models.Treasure, error)
	ListTreasures(ctx context.Context) ([]*models.Treasure, error)
	CreateTreasure(ctx context.Context, input *models.CreateTreasureInput) (int64, error)
	UpdateTreasure(ctx context.Context, id int64, input *models.UpdateTreasureInput) error
	DeleteTreasure(ctx context.Context, id int64) error
}

type SQLCTreasureRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCTreasureRepository(db *sql.DB) *SQLCTreasureRepository {
	return &SQLCTreasureRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCTreasureRepository) GetTreasure(ctx context.Context, id int64) (*models.Treasure, error) {
	treasure, err := r.q.GetTreasure(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("treasure", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbTreasureToModel(treasure), nil
}

func (r *SQLCTreasureRepository) GetTreasureByCharacter(ctx context.Context, characterID int64) (*models.Treasure, error) {
	nullableCharacterID := sql.NullInt64{
		Int64: characterID,
		Valid: true,
	}

	treasure, err := r.q.GetTreasureByCharacter(ctx, nullableCharacterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("treasure for character", characterID)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbTreasureToModel(treasure), nil
}

func (r *SQLCTreasureRepository) ListTreasures(ctx context.Context) ([]*models.Treasure, error) {
	treasures, err := r.q.ListTreasures(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.Treasure, len(treasures))
	for i, treasure := range treasures {
		result[i] = mapDbTreasureToModel(treasure)
	}
	return result, nil
}

func (r *SQLCTreasureRepository) CreateTreasure(ctx context.Context, input *models.CreateTreasureInput) (int64, error) {
	var characterID sql.NullInt64
	if input.CharacterID != nil {
		characterID.Int64 = *input.CharacterID
		characterID.Valid = true
	}

	result, err := r.q.CreateTreasure(ctx, sqlcdb.CreateTreasureParams{
		CharacterID:    characterID,
		PlatinumCoins:  int64(input.PlatinumCoins),
		GoldCoins:      int64(input.GoldCoins),
		ElectrumCoins:  int64(input.ElectrumCoins),
		SilverCoins:    int64(input.SilverCoins),
		CopperCoins:    int64(input.CopperCoins),
		Gems:           sql.NullString{String: input.Gems, Valid: input.Gems != ""},
		ArtObjects:     sql.NullString{String: input.ArtObjects, Valid: input.ArtObjects != ""},
		OtherValuables: sql.NullString{String: input.OtherValuables, Valid: input.OtherValuables != ""},
		TotalValueGold: input.TotalValueGold,
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

func (r *SQLCTreasureRepository) UpdateTreasure(ctx context.Context, id int64, input *models.UpdateTreasureInput) error {
	_, err := r.GetTreasure(ctx, id)
	if err != nil {
		return err
	}

	_, err = r.q.UpdateTreasure(ctx, sqlcdb.UpdateTreasureParams{
		PlatinumCoins:  int64(input.PlatinumCoins),
		GoldCoins:      int64(input.GoldCoins),
		ElectrumCoins:  int64(input.ElectrumCoins),
		SilverCoins:    int64(input.SilverCoins),
		CopperCoins:    int64(input.CopperCoins),
		Gems:           sql.NullString{String: input.Gems, Valid: input.Gems != ""},
		ArtObjects:     sql.NullString{String: input.ArtObjects, Valid: input.ArtObjects != ""},
		OtherValuables: sql.NullString{String: input.OtherValuables, Valid: input.OtherValuables != ""},
		TotalValueGold: input.TotalValueGold,
		ID:             id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCTreasureRepository) DeleteTreasure(ctx context.Context, id int64) error {
	_, err := r.GetTreasure(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeleteTreasure(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func mapDbTreasureToModel(treasure sqlcdb.Treasure) *models.Treasure {
	var characterID *int64
	if treasure.CharacterID.Valid {
		id := treasure.CharacterID.Int64
		characterID = &id
	}

	return &models.Treasure{
		ID:             treasure.ID,
		CharacterID:    characterID,
		PlatinumCoins:  int(treasure.PlatinumCoins),
		GoldCoins:      int(treasure.GoldCoins),
		ElectrumCoins:  int(treasure.ElectrumCoins),
		SilverCoins:    int(treasure.SilverCoins),
		CopperCoins:    int(treasure.CopperCoins),
		Gems:           treasure.Gems.String,
		ArtObjects:     treasure.ArtObjects.String,
		OtherValuables: treasure.OtherValuables.String,
		TotalValueGold: treasure.TotalValueGold,
		CreatedAt:      treasure.CreatedAt,
		UpdatedAt:      treasure.UpdatedAt,
	}
}
