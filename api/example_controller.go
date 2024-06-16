package api

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

// From https://jskim1991.medium.com/go-building-an-application-using-echo-framework-with-tests-controller-e4ca1187478c

// API
type Controller struct {
	Repository
}

func (m *Controller) GetAllBooks(c echo.Context) error {
	fetchedBooks, _ := m.Repository.FindAll()
	var books []Book
	for _, fetchedBook := range fetchedBooks {
		book := Book{
			Isbn:   fetchedBook.Isbn,
			Title:  fetchedBook.Title,
			Author: fetchedBook.Author,
		}
		books = append(books, book)
	}

	return c.JSON(http.StatusOK, books)
}

// DTO
type Book struct {
	Isbn   string
	Title  string
	Author string
}

// ORM
type BookEntity struct {
	gorm.Model
	Isbn   string
	Title  string
	Author string
}

type Repository interface {
	FindAll() ([]BookEntity, error)
}

type DefaultRepository struct {
}

func (m *DefaultRepository) FindAll() ([]BookEntity, error) {
	return []BookEntity{{
		Model:  gorm.Model{ID: 1},
		Isbn:   "9780321278654",
		Title:  "Extreme Programming Explained: Embrace Change",
		Author: "Kent Beck, Cynthia Andres",
	}}, nil
}
