package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *HttpHandler[AccountEntity]) getAccounts(c echo.Context) error {
	h.Logger.Info("Requesting all accounts slog")
	accounts, _ := h.Repository.List(c.Request().Context())
	if err := c.Bind(accounts); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, accounts)
}
