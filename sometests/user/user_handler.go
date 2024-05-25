package user

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	User struct {
		Name  string `json:"name" form:"name"`
		Email string `json:"email" form:"email"`
	}
	Handler struct {
		db map[string]*User
	}
)

func (h *Handler) createUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
}

func (h *Handler) GetUser(c echo.Context) error {
	email := c.Param("email")
	user := h.db[email]
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, user)
}

func Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello world")
}
