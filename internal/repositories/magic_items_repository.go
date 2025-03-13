package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type MagicItemRepository interface {
	GetMagicItem(ctx context.Context, id int64) (*models.MagicItem, error)
	GetMagicItemByName(ctx context.Context, name string) (*models.MagicItem, error)
	ListMagicItems(ctx context.Context) ([]*models.MagicItem, error)
	ListMagicItemsByType(ctx context.Context, itemType string) ([]*models.MagicItem, error)
	CreateMagicItem(ctx context.Context, input *models.CreateMagicItemInput) (int64, error)
	UpdateMagicItem(ctx context.Context, id int64, input *models.UpdateMagicItemInput) error
	DeleteMagicItem(ctx context.Context, id int64) error
}

type SQLCMagicItemRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCMagicItemRepository(db *sql.DB) *SQLCMagicItemRepository {
	return &SQLCMagicItemRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCMagicItemRepository) GetMagicItem(ctx context.Context, id int64) (*models.MagicItem, error) {
	magicItem, err := r.q.GetMagicItem(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("magic item", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbMagicItemToModel(magicItem), nil
}

func (r *SQLCMagicItemRepository) GetMagicItemByName(ctx context.Context, name string) (*models.MagicItem, error) {
	magicItem, err := r.q.GetMagicItemByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("magic item", name)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbMagicItemToModel(magicItem), nil
}

func (r *SQLCMagicItemRepository) ListMagicItems(ctx context.Context) ([]*models.MagicItem, error) {
	magicItems, err := r.q.ListMagicItems(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.MagicItem, len(magicItems))
	for i, item := range magicItems {
		result[i] = mapDbMagicItemToModel(item)
	}
	return result, nil
}

func (r *SQLCMagicItemRepository) ListMagicItemsByType(ctx context.Context, itemType string) ([]*models.MagicItem, error) {
	magicItems, err := r.q.ListMagicItemsByType(ctx, itemType)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.MagicItem, len(magicItems))
	for i, item := range magicItems {
		result[i] = mapDbMagicItemToModel(item)
	}
	return result, nil
}

func (r *SQLCMagicItemRepository) CreateMagicItem(ctx context.Context, input *models.CreateMagicItemInput) (int64, error) {
	var charges sql.NullInt64
	if input.Charges != nil {
		charges.Int64 = int64(*input.Charges)
		charges.Valid = true
	}

	result, err := r.q.CreateMagicItem(ctx, sqlcdb.CreateMagicItemParams{
		Name:        input.Name,
		ItemType:    input.ItemType,
		Description: input.Description,
		Charges:     charges,
		Cost:        input.Cost,
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

func (r *SQLCMagicItemRepository) UpdateMagicItem(ctx context.Context, id int64, input *models.UpdateMagicItemInput) error {
	_, err := r.GetMagicItem(ctx, id)
	if err != nil {
		return err
	}

	var charges sql.NullInt64
	if input.Charges != nil {
		charges.Int64 = int64(*input.Charges)
		charges.Valid = true
	}

	_, err = r.q.UpdateMagicItem(ctx, sqlcdb.UpdateMagicItemParams{
		Name:        input.Name,
		ItemType:    input.ItemType,
		Description: input.Description,
		Charges:     charges,
		Cost:        input.Cost,
		Weight:      int64(input.Weight),
		ID:          id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCMagicItemRepository) DeleteMagicItem(ctx context.Context, id int64) error {
	_, err := r.GetMagicItem(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeleteMagicItem(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func mapDbMagicItemToModel(item sqlcdb.MagicItem) *models.MagicItem {
	var charges *int
	if item.Charges.Valid {
		chargesVal := int(item.Charges.Int64)
		charges = &chargesVal
	}

	return &models.MagicItem{
		ID:          item.ID,
		Name:        item.Name,
		ItemType:    item.ItemType,
		Description: item.Description,
		Charges:     charges,
		Cost:        item.Cost,
		Weight:      int(item.Weight),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}