package api

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllBooks(t *testing.T) {
	t.Run("should return 200 status ok", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepository := MockRepository{}
		mockRepository.On("FindAll").Return([]BookEntity{}, nil)

		controller := Controller{&mockRepository}
		controller.GetAllBooks(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should return books", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepository := MockRepository{}
		mockRepository.On("FindAll").Return([]BookEntity{
			{
				Model:  gorm.Model{ID: 1},
				Isbn:   "999",
				Title:  "Learn Something",
				Author: "Jay",
			},
		}, nil)

		controller := Controller{&mockRepository}
		controller.GetAllBooks(c)

		var books []Book
		json.Unmarshal(rec.Body.Bytes(), &books)
		assert.Equal(t, 1, len(books))
		assert.Equal(t, "123", books[0].Isbn)
		assert.Equal(t, "Learn Something", books[0].Title)
		assert.Equal(t, "Jay", books[0].Author)
	})

	t.Run("should call repository to fetch books", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepository := MockRepository{}
		mockRepository.On("FindAll").Return([]BookEntity{}, nil)

		controller := Controller{&mockRepository}
		controller.GetAllBooks(c)

		mockRepository.AssertExpectations(t)
	})
}

// To mock Gorm repository interface to access book
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindAll() ([]BookEntity, error) {
	args := m.Called()
	return args[0].([]BookEntity), args.Error(1)
}
