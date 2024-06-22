package mocks

import (
	"banana-back/domain"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) FindAll(ctx context.Context) ([]domain.Account, error) {
	args := m.Called()
	return args[0].([]domain.Account), args.Error(1)
}

func (m *MockAccountRepository) FindById(ctx context.Context, id int64) (*domain.Account, error) {
	args := m.Called()
	return args[0].(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) Create(ctx context.Context, input *domain.Account) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAccountRepository) Update(ctx context.Context, input *domain.Account) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAccountRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called()
	return args.Error(0)
}
