package application

import (
	"banana-back/domain/profile"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
)

// TODO to be tested
func CreateProfileUseCase(logger *slog.Logger, repository profile.ProfileRepository) func(c echo.Context, command *profile.CreateProfileCommand) error {
	return func(c echo.Context, command *profile.CreateProfileCommand) error {
		var existingProfile, err = repository.FindBy(c.Request().Context(), command.Username)

		var validatedCommand, validationError = profile.ValidateProfileUsername(
			command,
			existingProfile,
		)
		if validationError != nil {
			logger.Warn("Error while creating profile: ", err)
			return validationError
		}

		if validatedCommand != nil {
			if err := repository.Create(c.Request().Context(), validatedCommand); err != nil {
				logger.Warn("Error while creating profile: ", err)
				return fmt.Errorf("failed to migrate database: %w", err)
			}
		}

		logger.Info("Profile {} created", validatedCommand)
		return nil
	}
}
