package main

import (
	"banana-back/hello"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, hello.Hello()+", World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
