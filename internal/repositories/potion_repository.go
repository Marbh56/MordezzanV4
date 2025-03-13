package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type PotionRepository interface {
	GetPotion(ctx context.Context, id int64) (*models.Potion, error)
	GetPotionByName(ctx context.Context, name string) (*models.Potion, error)
	ListPotions(ctx context.Context) ([]*models.Potion, error)
	CreatePotion(ctx context.Context, input *models.CreatePotionInput) (int64, error)
	UpdatePotion(ctx context.Context, id int64, input *models.UpdatePotionInput) error
	DeletePotion(ctx context.Context, id int64) error
}

type SQLCPotionRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCPotionRepository(db *sql.DB) *SQLCPotionRepository {
	return &SQLCPotionRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCPotionRepository) GetPotion(ctx context.Context, id int64) (*models.Potion, error) {
	potion, err := r.q.GetPotion(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("potion", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbPotionToModel(potion), nil
}

func (r *SQLCPotionRepository) GetPotionByName(ctx context.Context, name string) (*models.Potion, error) {
	potion, err := r.q.GetPotionByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("potion", name)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbPotionToModel(potion), nil
}

func (r *SQLCPotionRepository) ListPotions(ctx context.Context) ([]*models.Potion, error) {
	potions, err := r.q.ListPotions(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.Potion, len(potions))
	for i, potion := range potions {
		result[i] = mapDbPotionToModel(potion)
	}
	return result, nil
}

func (r *SQLCPotionRepository) CreatePotion(ctx context.Context, input *models.CreatePotionInput) (int64, error) {
	result, err := r.q.CreatePotion(ctx, sqlcdb.CreatePotionParams{
		Name:        input.Name,
		Description: input.Description,
		Uses:        int64(input.Uses),
		Weight:      int64(input.Weight),
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

func (r *SQLCPotionRepository) UpdatePotion(ctx context.Context, id int64, input *models.UpdatePotionInput) error {
	_, err := r.GetPotion(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.UpdatePotion(ctx, sqlcdb.UpdatePotionParams{
		Name:        input.Name,
		Description: input.Description,
		Uses:        int64(input.Uses),
		Weight:      int64(input.Weight),
		ID:          id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCPotionRepository) DeletePotion(ctx context.Context, id int64) error {
	_, err := r.GetPotion(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeletePotion(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func mapDbPotionToModel(potion sqlcdb.Potion) *models.Potion {
	return &models.Potion{
		ID:          potion.ID,
		Name:        potion.Name,
		Description: potion.Description,
		Uses:        int(potion.Uses),
		Weight:      int(potion.Weight),
		CreatedAt:   potion.CreatedAt,
		UpdatedAt:   potion.UpdatedAt,
	}
}
