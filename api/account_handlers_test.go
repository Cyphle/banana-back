package api

import (
	"banana-back/repositories"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAccountRepository struct {
}

func (m MockAccountRepository) List(ctx context.Context) ([]repositories.AccountEntity, error) {
	res := make([]repositories.AccountEntity, 1)
	res[0] = repositories.AccountEntity{
		ID:   1,
		Name: "Coucou",
	}
	return res, nil
}

func NewHttpHandlerWithMock() HttpHandler[repositories.AccountEntity] {
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
