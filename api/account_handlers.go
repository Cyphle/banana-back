package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *HttpHandler) getAccounts(c echo.Context) error {
	h.Logger.Info("Requesting all accounts")
	accounts, _ := h.AccountRepository.List(c.Request().Context())
	if err := c.Bind(accounts); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, accounts)
}
