package repositories

import "context"

type Repository[T any] interface {
	List(ctx context.Context) ([]T, error)
	//Create(ctx context.Context, input *T) error
}
