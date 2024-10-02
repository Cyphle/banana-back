package api

import (
	"banana-back/domain/account"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// TODO peut être que les handler faudrait les mettre dans la couche application et en faire des use cases. si on garde les handlers ça fait une indirection en plus
func (h *AccountHttpHandler) getAccounts(c echo.Context) error {
	h.Logger.Info("Requesting all accounts")
	accounts, _ := h.Repository.FindAll(c.Request().Context())

	var accountViews []AccountView
	for _, account := range accounts {
		accountViews = append(accountViews, AccountView{
			ID:   account.ID,
			Name: account.Name,
		})
	}

	response := ArrayResponse[account.Account]{
		Data: accountViews,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AccountHttpHandler) findAccountHandler(c echo.Context) error {
	var accountId AccountIdPathParam
	err := c.Bind(&accountId)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	h.Logger.Info(fmt.Sprintf("Findind account %v", accountId))

	account, _ := h.Repository.FindById(c.Request().Context(), accountId.ID)

	return c.JSON(http.StatusOK, AccountView{
		ID:   account.ID,
		Name: account.Name,
	})
}

func (h *AccountHttpHandler) createAccountHandler(c echo.Context) error {
	h.Logger.Info("Creating an account")

	u := new(CreateAccountRequest)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	existingAccount, _ := h.Repository.FindOneByField(c.Request().Context(), "name", u.Name)

	// As domain is pure functions only, no need to inject as it does not make any side effect
	if accountToCreate, err := account.CreateAccount(
		&account.CreateAccountCommand{Name: u.Name},
		existingAccount,
	); err != nil {
		h.Logger.Error("failed to create an account: %w", err)
		return c.String(http.StatusBadRequest, err.Error())
	} else {
		if err := h.Repository.Create(c.Request().Context(), accountToCreate); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to create account: %w", err))
		} else {
			createdAccount, _ := h.Repository.FindOneByField(c.Request().Context(), "name", accountToCreate.Name)
			return c.JSON(http.StatusCreated, AccountView{
				ID:   createdAccount.ID,
				Name: createdAccount.Name,
			})
		}
	}
}

func (h *AccountHttpHandler) updateAccount(c echo.Context) error {
	h.Logger.Info("Update an account")

	u := new(UpdateAccountCommandView)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	h.Repository.Update(c.Request().Context(), &account.Account{
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
