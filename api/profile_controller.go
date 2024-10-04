package api

import (
	"banana-back/application"
	"banana-back/domain/profile"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

func CreateProfileHandler(logger *slog.Logger, repository profile.ProfileRepository) func(c echo.Context) error {
	return func(c echo.Context) error {
		u := new(profile.CreateProfileCommand)
		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		var err = application.CreateProfileUseCase(logger, repository)(c, u)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to create profile: %w", err))
		} else {
			return c.JSON(http.StatusCreated, u)
		}
	}
}

func ProfileRoutes(e *echo.Echo, logger *slog.Logger, repository profile.ProfileRepository) {
	logger.Info("Activating profiles routes")
	g := e.Group("/profiles")
	g.POST("", CreateProfileHandler(logger, repository))
}
