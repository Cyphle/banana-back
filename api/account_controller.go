package api

import (
	"banana-back/repositories"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

func ActivateAccountRoutes(e *echo.Echo, handler HttpHandler) {
	e.GET("/", handler.getAccounts)
}

type HttpHandler struct {
	Logger            *slog.Logger
	AccountRepository *repositories.AccountRepository
}

func (h *HttpHandler) getAccounts(c echo.Context) error {
	h.Logger.Info("Requesting all accounts")
	accounts, _ := h.AccountRepository.List(c.Request().Context())
	if err := c.Bind(accounts); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, accounts)
}
