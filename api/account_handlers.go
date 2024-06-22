package api

import (
	"banana-back/domain"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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
	var accountId AccountIdPathParam
	err := c.Bind(&accountId)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	h.Logger.Info(fmt.Sprintf("Findind account %v", accountId))

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

func (h *AccountHttpHandler) updateAccount(c echo.Context) error {
	h.Logger.Info("Update an account")

	u := new(UpdateAccountCommandView)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	h.Repository.Update(c.Request().Context(), &domain.Account{
		ID:   u.ID,
		Name: u.Name,
	})

	return c.JSON(http.StatusOK, u)
}

func (h *AccountHttpHandler) deleteAccount(c echo.Context) error {
	accountId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	h.Logger.Info(fmt.Sprintf("Deleting account %v", accountId))

	if err := h.Repository.Delete(c.Request().Context(), accountId); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	return c.NoContent(http.StatusOK)
}
