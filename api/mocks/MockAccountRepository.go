package mocks

import (
	"banana-back/domain/account"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) FindAll(ctx context.Context) ([]account.Account, error) {
	args := m.Called()
	return args[0].([]account.Account), args.Error(1)
}

func (m *MockAccountRepository) FindById(ctx context.Context, id int64) (*account.Account, error) {
	args := m.Called()
	return args[0].(*account.Account), args.Error(1)
}

func (m *MockAccountRepository) FindOneByField(ctx context.Context, field string, value string) (*account.Account, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args[0].(*account.Account), args.Error(1)
	}
}

func (m *MockAccountRepository) Create(ctx context.Context, input *account.Account) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAccountRepository) Update(ctx context.Context, input *account.Account) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAccountRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called()
	return args.Error(0)
}
