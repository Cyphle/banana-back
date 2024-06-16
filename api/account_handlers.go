package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *HttpHandler[Account]) getAccounts(c echo.Context) error {
	h.Logger.Info("Requesting all accounts")
	accounts, _ := h.Repository.List(c.Request().Context())
	if err := c.Bind(accounts); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, accounts)
}

// TODO to be tested
func (h *HttpHandler[Account]) createAccount(c echo.Context) error {
	h.Logger.Info("Creating an account")

	u := new(CreateAccountCommandView)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	//h.Repository.Create()

	fmt.Println(u)

	return c.NoContent(http.StatusNoContent)
}
