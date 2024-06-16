package main

import (
	"banana-back/api"
	"banana-back/config"
	"banana-back/initializers"
	"banana-back/repositories"
	"banana-back/sometests/user"
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

func MySum(xi ...int) int {
	sum := 0
	for _, y := range xi {
		sum += y
	}
	return sum
}

func main() {
	ctx := context.Background()
	log := slog.Default()

	conf := config.Get()

	// Setup database
	log.Info("Initializing database")
	dbClient, err := initializers.InitDatabase(ctx, conf.DB)
	if err != nil {
		log.Error("failed to init database", "err", err)
		return
	}

	// ECHO
	handler := api.NewAccountHttpHandler(
		log,
		repositories.NewAccountRepository(dbClient),
	)
	api.ActivateAccountRoutes(conf.WebServer, handler)

	conf.WebServer.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	conf.WebServer.GET("/users", user.Test)

	conf.WebServer.Logger.Fatal(conf.WebServer.Start(":1323"))
	// END ECHO
}
