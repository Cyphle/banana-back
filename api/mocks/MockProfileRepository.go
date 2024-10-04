package mocks

import (
	"banana-back/domain/profile"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockProfileRepository struct {
	mock.Mock
}

func (r *MockProfileRepository) FindBy(ctx context.Context, username string) (*profile.Profile, error) {
	args := r.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args[0].(*profile.Profile), args.Error(1)
	}
}

func (r *MockProfileRepository) Create(ctx context.Context, command *profile.CreateProfileCommand) error {
	args := r.Called()
	return args.Error(0)
}
