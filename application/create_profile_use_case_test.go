package application

import (
	"banana-back/api/mocks"
	"banana-back/domain/profile"
	"banana-back/repositories"
	"banana-back/testutils"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateProfileUseCase(t *testing.T) {
	t.Run("should process create profile use case", func(t *testing.T) {
		mockRep := mocks.MockProfileRepository{}
		useCase := CreateProfileUseCase(testutils.Logger, &mockRep)
		mockRep.On("FindBy").Return(
			nil,
			repositories.ErrProfileNotFound,
		).Once()
		mockRep.On("Create").Return(nil)

		err := useCase(context.Background(), &profile.CreateProfileCommand{})

		assert.NoError(t, err)
		mockRep.AssertExpectations(t)
	})
}
