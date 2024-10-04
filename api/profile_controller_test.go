package api

import (
	"banana-back/api/mocks"
	"banana-back/repositories"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateProfile(t *testing.T) {
	t.Run("should create a profile", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/profiles", strings.NewReader("{\n  \"username\": \"johnsmith\",\n  \"email\": \"johnsmith@banana.com\",\n  \"firstName\": \"John\",\n  \"lastName\": \"Smith\"}"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRep := mocks.MockProfileRepository{}
		mockRep.On("FindBy").Return(
			nil,
			repositories.ErrProfileNotFound,
		).Once()
		mockRep.On("Create").Return(nil)
		h := CreateProfileHandler(logger, &mockRep)

		// Assertions
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "{\"id\":1,\"name\":\"John Smith\"}\n", rec.Body.String())
		}
		mockRep.AssertExpectations(t)
	})
}
