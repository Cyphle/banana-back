package api

import (
	"banana-back/domain"
	"banana-back/repositories"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type AccountHttpHandler struct {
	Logger     *slog.Logger
	Repository repositories.Repository[domain.Account]
}

func NewAccountHttpHandler(logger *slog.Logger, repository repositories.Repository[domain.Account]) *AccountHttpHandler {
	return &AccountHttpHandler{
		Logger:     logger,
		Repository: repository,
	}
}

func ActivateAccountRoutes(e *echo.Echo, handler *AccountHttpHandler) {
	handler.Logger.Info("Activating account routes")
	g := e.Group("/accounts")
	g.GET("", handler.getAccounts)
	g.GET("/:id", handler.findAccount)
	g.POST("", handler.createAccount)
	g.PUT("/:id", handler.updateAccount)
	g.DELETE("/:id", handler.deleteAccount)
}
