package api

import (
	"banana-back/domain"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *AccountHttpHandler) getAccounts(c echo.Context) error {
	h.Logger.Info("Requesting all accounts")
	accounts, _ := h.Repository.FindAll(c.Request().Context())
	response := ArrayResponse[domain.Account]{
		Data: accounts,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AccountHttpHandler) findAccount(c echo.Context) error {
	h.Logger.Info("Findind account {}")

	var accountId AccountIdPathParam
	err := c.Bind(&accountId)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	account, _ := h.Repository.FindById(c.Request().Context(), accountId.ID)

	return c.JSON(http.StatusOK, account)
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

	return c.JSON(http.StatusOK, account)
}
