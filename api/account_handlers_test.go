package api

import (
	"banana-back/domain/account"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockAccountRepository struct {
}

func (m MockAccountRepository) List(ctx context.Context) ([]account.Account, error) {
	res := make([]account.Account, 1)
	res[0] = account.Account{
		ID:   1,
		Name: "Coucou",
	}
	return res, nil
}

func (m MockAccountRepository) Create(ctx context.Context, input *account.Account) error {
	return nil
}

func NewHttpHandlerWithMock() AccountHttpHandler {
	logger := slog.Default()
	mock := &MockAccountRepository{}
	return NewAccountHttpHandler(logger, mock)
}

func TestGetAccounts(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/accounts")
	h := NewHttpHandlerWithMock()

	// Assertions
	if assert.NoError(t, h.getAccounts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "{ data: [{ ID: 1, Name: \"coucou\"}]}", rec.Body.String())
	}
}

func TestCreateUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader("{ \"name\": \"John Smith\" }"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHttpHandlerWithMock()

	// Assertions
	if assert.NoError(t, h.createAccount(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		// TODO assert to have been called (cf https://jskim1991.medium.com/go-building-an-application-using-echo-framework-with-tests-controller-e4ca1187478c)
	}
}
