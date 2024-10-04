package application

import (
	"banana-back/domain/profile"
	"context"
	"fmt"
	"log/slog"
)

func CreateProfileUseCase(logger *slog.Logger, repository profile.ProfileRepository) func(c context.Context, command *profile.CreateProfileCommand) error {
	return func(c context.Context, command *profile.CreateProfileCommand) error {
		var existingProfile, err = repository.FindBy(c, command.Username)

		var validatedCommand, validationError = profile.ValidateProfileUsername(
			command,
			existingProfile,
		)
		if validationError != nil {
			logger.Warn("Error while creating profile: ", err)
			return validationError
		}

		if validatedCommand != nil {
			if err := repository.Create(c, validatedCommand); err != nil {
				logger.Warn("Error while creating profile: ", err)
				return fmt.Errorf("failed to migrate database: %w", err)
			}
		}

		logger.Info("Profile created")
		return nil
	}
}
