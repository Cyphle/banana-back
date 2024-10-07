package api

import (
	"banana-back/api/mocks"
	"banana-back/repositories"
	"banana-back/testutils"
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
		h := CreateProfileHandler(testutils.Logger, &mockRep)

		// Assertions
		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, "{\"username\":\"johnsmith\",\"email\":\"johnsmith@banana.com\",\"first_name\":\"\",\"last_name\":\"\"}\n", rec.Body.String())
		}
		mockRep.AssertExpectations(t)
	})
}
