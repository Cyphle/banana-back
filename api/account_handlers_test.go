package api

import (
	"banana-back/domain"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) List(ctx context.Context) ([]domain.Account, error) {
	args := m.Called()
	return args[0].([]domain.Account), args.Error(1)
}

func (m *MockAccountRepository) Create(ctx context.Context, input *domain.Account) error {
	return nil
}

func TestGetAccounts(t *testing.T) {
	logger := slog.Default()

	t.Run("should get accounts", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRep := MockAccountRepository{}
		mockRep.On("List").Return([]domain.Account{
			domain.Account{
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

		mockRep := MockAccountRepository{}
		mockRep.On("List").Return([]domain.Account{
			domain.Account{
				ID:   1,
				Name: "Coucou",
			},
		}, nil)
		handler := NewAccountHttpHandler(logger, &mockRep)

		handler.getAccounts(c)

		mockRep.AssertExpectations(t)
	})
}

func TestCreateUser(t *testing.T) {
	logger := slog.Default()

	t.Run("should create an account", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader("{ \"name\": \"John Smith\" }"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := NewAccountHttpHandler(logger, &MockAccountRepository{})

		// Assertions
		if assert.NoError(t, h.createAccount(c)) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
			// TODO assert to have been called (cf https://jskim1991.medium.com/go-building-an-application-using-echo-framework-with-tests-controller-e4ca1187478c)
		}
	})
}
