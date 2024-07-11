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

	accountViews := make([]AccountView, len(accounts))
	for _, account := range accounts {
		accountViews = append(accountViews, AccountView{
			ID:   account.ID,
			Name: account.Name,
		})
	}

	response := ArrayResponse[domain.Account]{
		Data: accountViews,
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

func (h *AccountHttpHandler) createAccountHandler(c echo.Context) error {
	h.Logger.Info("Creating an account")

	u := new(CreateAccountRequest)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	// TODO ici ça devrait être une couche métier qui reçoit une request, qui transforme en commande si ok et qui save dans le repo
	/*
		genre une fonction qu'on peut passer en paramètre du createAccountHandler pour faire de l'injection et testable

		func createAccount(repository *AccountRepository, request Request)
	*/
	// TODO à injecte comme ça func (h *AccountHttpHandler) createAccountHandler(c echo.Context, f func() domain.Account) error
	domain.CreateAccount(h.Repository, u)

	// TODO to be deleted from here
	account := &domain.Account{
		ID:   -1,
		Name: u.Name,
	}
	h.Repository.Create(c.Request().Context(), account)
	// TODO to here

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
