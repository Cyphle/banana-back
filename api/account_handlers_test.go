package api

import (
	"banana-back/domain/account"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

func NewHttpHandlerWithMock() HttpHandler[account.Account] {
	mock := &MockAccountRepository{}
	return NewHttpHandler(mock)
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

//func TestCreateUser(t *testing.T) {
//	// Setup
//	e := echo.New()
//	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//	h := &handler{mockDB}
//
//	// Assertions
//	if assert.NoError(t, h.createUser(c)) {
//		assert.Equal(t, http.StatusCreated, rec.Code)
//		assert.Equal(t, userJSON, rec.Body.String())
//	}
//}
