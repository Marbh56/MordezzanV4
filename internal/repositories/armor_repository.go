package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type ArmorRepository interface {
	GetArmor(ctx context.Context, id int64) (*models.Armor, error)
	GetArmorByName(ctx context.Context, name string) (*models.Armor, error)
	ListArmors(ctx context.Context) ([]*models.Armor, error)
	CreateArmor(ctx context.Context, input *models.CreateArmorInput) (int64, error)
	UpdateArmor(ctx context.Context, id int64, input *models.UpdateArmorInput) error
	DeleteArmor(ctx context.Context, id int64) error
}

type SQLCArmorRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCArmorRepository(db *sql.DB) *SQLCArmorRepository {
	return &SQLCArmorRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCArmorRepository) GetArmor(ctx context.Context, id int64) (*models.Armor, error) {
	armor, err := r.q.GetArmor(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("armor", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbArmorToModel(armor), nil
}

func (r *SQLCArmorRepository) GetArmorByName(ctx context.Context, name string) (*models.Armor, error) {
	armor, err := r.q.GetArmorByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("armor", name)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbArmorToModel(armor), nil
}

func (r *SQLCArmorRepository) ListArmors(ctx context.Context) ([]*models.Armor, error) {
	armors, err := r.q.ListArmors(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.Armor, len(armors))
	for i, armor := range armors {
		result[i] = mapDbArmorToModel(armor)
	}
	return result, nil
}

func (r *SQLCArmorRepository) CreateArmor(ctx context.Context, input *models.CreateArmorInput) (int64, error) {
	result, err := r.q.CreateArmor(ctx, sqlcdb.CreateArmorParams{
		Name:            input.Name,
		ArmorType:       input.ArmorType,
		Ac:              int64(input.AC),
		Cost:            input.Cost,
		DamageReduction: int64(input.DamageReduction),
		Weight:          int64(input.Weight),
		WeightClass:     input.WeightClass,
		MovementRate:    int64(input.MovementRate),
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

func (r *SQLCArmorRepository) UpdateArmor(ctx context.Context, id int64, input *models.UpdateArmorInput) error {
	// First check if the armor exists
	_, err := r.GetArmor(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.UpdateArmor(ctx, sqlcdb.UpdateArmorParams{
		Name:            input.Name,
		ArmorType:       input.ArmorType,
		Ac:              int64(input.AC),
		Cost:            input.Cost,
		DamageReduction: int64(input.DamageReduction),
		Weight:          int64(input.Weight),
		WeightClass:     input.WeightClass,
		MovementRate:    int64(input.MovementRate),
		ID:              id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCArmorRepository) DeleteArmor(ctx context.Context, id int64) error {
	_, err := r.GetArmor(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeleteArmor(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func mapDbArmorToModel(armor sqlcdb.Armor) *models.Armor {
	return &models.Armor{
		ID:              armor.ID,
		Name:            armor.Name,
		ArmorType:       armor.ArmorType,
		AC:              int(armor.Ac),
		Cost:            armor.Cost,
		DamageReduction: int(armor.DamageReduction),
		Weight:          int(armor.Weight),
		WeightClass:     armor.WeightClass,
		MovementRate:    int(armor.MovementRate),
		CreatedAt:       armor.CreatedAt,
		UpdatedAt:       armor.UpdatedAt,
	}
}
