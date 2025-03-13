package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type EquipmentRepository interface {
	GetEquipment(ctx context.Context, id int64) (*models.Equipment, error)
	GetEquipmentByName(ctx context.Context, name string) (*models.Equipment, error)
	ListEquipment(ctx context.Context) ([]*models.Equipment, error)
	CreateEquipment(ctx context.Context, input *models.CreateEquipmentInput) (int64, error)
	UpdateEquipment(ctx context.Context, id int64, input *models.UpdateEquipmentInput) error
	DeleteEquipment(ctx context.Context, id int64) error
}

type SQLCEquipmentRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCEquipmentRepository(db *sql.DB) *SQLCEquipmentRepository {
	return &SQLCEquipmentRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCEquipmentRepository) GetEquipment(ctx context.Context, id int64) (*models.Equipment, error) {
	equipment, err := r.q.GetEquipment(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("equipment", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbEquipmentToModel(equipment), nil
}

func (r *SQLCEquipmentRepository) GetEquipmentByName(ctx context.Context, name string) (*models.Equipment, error) {
	equipment, err := r.q.GetEquipmentByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("equipment", name)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbEquipmentToModel(equipment), nil
}

func (r *SQLCEquipmentRepository) ListEquipment(ctx context.Context) ([]*models.Equipment, error) {
	equipmentList, err := r.q.ListEquipment(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.Equipment, len(equipmentList))
	for i, equipment := range equipmentList {
		result[i] = mapDbEquipmentToModel(equipment)
	}
	return result, nil
}

func (r *SQLCEquipmentRepository) CreateEquipment(ctx context.Context, input *models.CreateEquipmentInput) (int64, error) {
	result, err := r.q.CreateEquipment(ctx, sqlcdb.CreateEquipmentParams{
		Name:        input.Name,
		Description: input.Description,
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

func (r *SQLCEquipmentRepository) UpdateEquipment(ctx context.Context, id int64, input *models.UpdateEquipmentInput) error {
	_, err := r.GetEquipment(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.UpdateEquipment(ctx, sqlcdb.UpdateEquipmentParams{
		Name:        input.Name,
		Description: input.Description,
		Cost:        input.Cost,
		Weight:      int64(input.Weight),
		ID:          id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCEquipmentRepository) DeleteEquipment(ctx context.Context, id int64) error {
	_, err := r.GetEquipment(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeleteEquipment(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func mapDbEquipmentToModel(equipment sqlcdb.Equipment) *models.Equipment {
	return &models.Equipment{
		ID:          equipment.ID,
		Name:        equipment.Name,
		Description: equipment.Description,
		Cost:        equipment.Cost,
		Weight:      int(equipment.Weight),
		CreatedAt:   equipment.CreatedAt,
		UpdatedAt:   equipment.UpdatedAt,
	}
}
