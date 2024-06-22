package main

import (
	"banana-back/api"
	"banana-back/config"
	"banana-back/initializers"
	"banana-back/repositories"
	"context"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

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

	conf.WebServer.Logger.Fatal(conf.WebServer.Start(fmt.Sprintf("%s%d", ":", conf.WebServerPort)))
	// END ECHO
}
