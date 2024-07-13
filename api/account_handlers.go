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

func (h *AccountHttpHandler) createAccountHandler(createAccount func(command *domain.CreateAccountCommand, existingAccount *domain.Account) (*domain.Account, error)) func(echo.Context) error {
	return func(c echo.Context) error {
		h.Logger.Info("Creating an account")

		u := new(CreateAccountRequest)
		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		existingAccount, _ := h.Repository.FindOneByField(c.Request().Context(), "name", u.Name)

		if accountToCreate, err := createAccount(
			&domain.CreateAccountCommand{Name: u.Name},
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
