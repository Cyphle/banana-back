package api

import (
	"banana-back/api/mocks"
	"banana-back/domain/account"
	account2 "banana-back/repositories/account"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

func TestFindAccounts(t *testing.T) {
	t.Run("should get accounts", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRep := mocks.MockAccountRepository{}
		mockRep.On("FindAll").Return([]account.Account{
			account.Account{
				ID:   1,
				Name: "Coucou",
			},
		}, nil)
		handler := NewAccountHttpHandler(logger, &mockRep)

		// Assertions
		if assert.NoError(t, handler.getAccounts(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "{\"data\":[{\"id\":1,\"name\":\"Coucou\"}]}\n", rec.Body.String())
		}
	})

	t.Run("should call repository when getting accounts", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRep := mocks.MockAccountRepository{}
		mockRep.On("FindAll").Return([]account.Account{
			account.Account{
				ID:   1,
				Name: "Coucou",
			},
		}, nil)
		handler := NewAccountHttpHandler(logger, &mockRep)

		handler.getAccounts(c)

		mockRep.AssertExpectations(t)
	})
}

func TestFindAccountById(t *testing.T) {
	t.Run("should find one account for given id", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockRep := mocks.MockAccountRepository{}
		mockRep.On("FindById").Return(&account.Account{
			ID:   1,
			Name: "Coucou",
		}, nil)
		handler := NewAccountHttpHandler(logger, &mockRep)

		// Assertions
		if assert.NoError(t, handler.findAccountHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "{\"id\":1,\"name\":\"Coucou\"}\n", rec.Body.String())
		}
		mockRep.AssertExpectations(t)
	})
}

func TestCreateAccount(t *testing.T) {
	t.Run("should create an account", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader("{ \"name\": \"John Smith\" }"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRep := mocks.MockAccountRepository{}
		mockRep.On("FindOneByField").Return(
			nil,
			account2.ErrAccountNotFound,
		).Once()
		mockRep.On("FindOneByField").Return(&account.Account{
			ID:   1,
			Name: "John Smith",
		}, nil).Once()
		mockRep.On("Create").Return(nil)
		h := NewAccountHttpHandler(logger, &mockRep)

		// Assertions
		if assert.NoError(t, h.createAccountHandler(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "{\"id\":1,\"name\":\"John Smith\"}\n", rec.Body.String())
		}
		mockRep.AssertExpectations(t)
	})
}

func TestUpdateAccount(t *testing.T) {
	t.Run("should update an account", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/accounts", strings.NewReader("{ \"name\": \"John Smith\" }"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockRep := mocks.MockAccountRepository{}
		mockRep.On("Update").Return(nil)
		h := NewAccountHttpHandler(logger, &mockRep)

		// Assertions
		if assert.NoError(t, h.updateAccount(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "{\"id\":1,\"name\":\"John Smith\"}\n", rec.Body.String())
		}
		mockRep.AssertExpectations(t)
	})
}

func TestDeleteAccount(t *testing.T) {
	t.Run("should delete an account for given id", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/accounts", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockRep := mocks.MockAccountRepository{}
		mockRep.On("Delete").Return(nil)
		handler := NewAccountHttpHandler(logger, &mockRep)

		// Assertions
		if assert.NoError(t, handler.deleteAccount(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
		mockRep.AssertExpectations(t)
	})
}
