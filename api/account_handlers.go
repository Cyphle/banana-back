package api

import (
	"banana-back/domain"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *AccountHttpHandler) getAccounts(c echo.Context) error {
	h.Logger.Info("Requesting all accounts")
	accounts, _ := h.Repository.List(c.Request().Context())
	if err := c.Bind(accounts); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, accounts)
}

func (h *AccountHttpHandler) createAccount(c echo.Context) error {
	h.Logger.Info("Creating an account")

	u := new(CreateAccountCommandView)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	account := &domain.Account{
		ID:   -1,
		Name: u.Name,
	}
	h.Repository.Create(c.Request().Context(), account)

	fmt.Println(u)

	return c.NoContent(http.StatusNoContent)
}
