package api

import (
	"banana-back/domain/account"
	"banana-back/repositories"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type AccountHttpHandler struct {
	Logger     *slog.Logger
	Repository repositories.Repository[account.Account]
}

func NewAccountHttpHandler(logger *slog.Logger, repository repositories.Repository[account.Account]) AccountHttpHandler {
	return AccountHttpHandler{
		Logger:     logger,
		Repository: repository,
	}
}

func ActivateAccountRoutes(e *echo.Echo, handler AccountHttpHandler) {
	handler.Logger.Info("Activating account routes")
	g := e.Group("/accounts")
	g.GET("", handler.getAccounts)
	g.POST("", handler.createAccount)
}
