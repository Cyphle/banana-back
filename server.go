package main

import (
	"banana-back/config"
	"banana-back/initializers"
	"context"
	"log/slog"
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
	_, err := initializers.InitDatabase(ctx, conf.DB)
	if err != nil {
		log.Error("failed to init database", "err", err)
		return
	}

	// Setup Echo

	// =============================
	// ===> OLD TESTS TO SEtuP TOOLS

	// ECHO
	//e := echo.New()
	//e.GET("/", func(c echo.Context) error {
	//	return c.String(http.StatusOK, hello.Hello()+", World!")
	//})
	//
	//e.GET("/users", user.Test)
	//e.Logger.Fatal(e.Start(":1323"))
	// END ECHO
}
