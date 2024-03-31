package main

import (
	"banana-back/db"
	"banana-back/hello"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	fmt.Println("TEST")
	fmt.Println("Accounts " + db.ToString())

	db.Add(db.BankAccount{
		Name: "My account",
	})

	fmt.Println("Accounts now: " + db.ToString())
	fmt.Println("END TEST")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, hello.Hello()+", World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
