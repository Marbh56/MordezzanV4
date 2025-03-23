package repositories

import "context"

// Repository is a generic interface for CRUD operations
type Repository[T any, CreateInput any, UpdateInput any] interface {
	Get(ctx context.Context, id int64) (*T, error)
	List(ctx context.Context) ([]*T, error)
	Create(ctx context.Context, input *CreateInput) (int64, error)
	Update(ctx context.Context, id int64, input *UpdateInput) error
	Delete(ctx context.Context, id int64) error
}
