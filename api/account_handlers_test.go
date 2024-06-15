package api

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccounts(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/accounts")
	h := &HttpHandler{}

	// Assertions
	if assert.NoError(t, h.getAccounts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		//assert.Equal(t, userJSON, rec.Body.String())
	}
}
