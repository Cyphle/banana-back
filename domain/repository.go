package domain

import "context"

type Repository[T any] interface {
	FindAll(ctx context.Context) ([]T, error)
	FindById(ctx context.Context, id int64) (*T, error)
	Create(ctx context.Context, input *T) error
	Update(ctx context.Context, input *T) error
	Delete(ctx context.Context, id int64) error
}
