package api

import (
	"banana-back/domain/account"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type AccountHttpHandler struct {
	Logger     *slog.Logger
	Repository account.AccountRepository
}

func NewAccountHttpHandler(logger *slog.Logger, repository account.AccountRepository) *AccountHttpHandler {
	return &AccountHttpHandler{
		Logger:     logger,
		Repository: repository,
	}
}

func ActivateAccountRoutes(e *echo.Echo, handler *AccountHttpHandler) {
	handler.Logger.Info("Activating account routes")
	g := e.Group("/accounts")
	g.GET("", handler.getAccounts)
	g.GET("/:id", handler.findAccountHandler)
	g.POST("", handler.createAccountHandler)
	g.PUT("/:id", handler.updateAccount)
	g.DELETE("/:id", handler.deleteAccount)
}
