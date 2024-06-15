package api

import (
	"context"
)

type RepositoryInterface[T any] interface {
	List(ctx context.Context) ([]T, error)
}

type MockAccountRepository struct {
}

func (m MockAccountRepository) List(ctx context.Context) ([]string, error) {
	res := make([]string, 0)
	return res, nil
}

type SomeReceiver[T any] struct {
	Repository RepositoryInterface[T]
}

func Other() {
	mock := MockAccountRepository{}
	h := SomeReceiver[string]{
		Repository: mock,
	}
}
