package repositories

import (
	"context"
	"database/sql"
	"errors"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

type ContainerRepository interface {
	GetContainer(ctx context.Context, id int64) (*models.Container, error)
	GetContainerByName(ctx context.Context, name string) (*models.Container, error)
	ListContainers(ctx context.Context) ([]*models.Container, error)
	CreateContainer(ctx context.Context, input *models.CreateContainerInput) (int64, error)
	UpdateContainer(ctx context.Context, id int64, input *models.UpdateContainerInput) error
	DeleteContainer(ctx context.Context, id int64) error
}

type SQLCContainerRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCContainerRepository(db *sql.DB) *SQLCContainerRepository {
	return &SQLCContainerRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCContainerRepository) GetContainer(ctx context.Context, id int64) (*models.Container, error) {
	container, err := r.q.GetContainer(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("container", id)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbContainerToModel(container), nil
}

func (r *SQLCContainerRepository) GetContainerByName(ctx context.Context, name string) (*models.Container, error) {
	container, err := r.q.GetContainerByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("container", name)
		}
		return nil, apperrors.NewDatabaseError(err)
	}
	return mapDbContainerToModel(container), nil
}

func (r *SQLCContainerRepository) ListContainers(ctx context.Context) ([]*models.Container, error) {
	containers, err := r.q.ListContainers(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}
	result := make([]*models.Container, len(containers))
	for i, container := range containers {
		result[i] = mapDbContainerToModel(container)
	}
	return result, nil
}

func (r *SQLCContainerRepository) CreateContainer(ctx context.Context, input *models.CreateContainerInput) (int64, error) {
	result, err := r.q.CreateContainer(ctx, sqlcdb.CreateContainerParams{
		Name:         input.Name,
		MaxWeight:    int64(input.MaxWeight),
		AllowedItems: input.AllowedItems,
		Cost:         input.Cost,
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

func (r *SQLCContainerRepository) UpdateContainer(ctx context.Context, id int64, input *models.UpdateContainerInput) error {
	_, err := r.GetContainer(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.UpdateContainer(ctx, sqlcdb.UpdateContainerParams{
		Name:         input.Name,
		MaxWeight:    int64(input.MaxWeight),
		AllowedItems: input.AllowedItems,
		Cost:         input.Cost,
		ID:           id,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func (r *SQLCContainerRepository) DeleteContainer(ctx context.Context, id int64) error {
	_, err := r.GetContainer(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.q.DeleteContainer(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	return nil
}

func mapDbContainerToModel(container sqlcdb.Container) *models.Container {
	return &models.Container{
		ID:           container.ID,
		Name:         container.Name,
		MaxWeight:    int(container.MaxWeight),
		AllowedItems: container.AllowedItems,
		Cost:         container.Cost,
		CreatedAt:    container.CreatedAt,
		UpdatedAt:    container.UpdatedAt,
	}
}
