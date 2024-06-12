package config

import "github.com/labstack/echo/v4"

type Config struct {
	DB        Database
	WebServer *echo.Echo
}

func Get() *Config {
	conf := Config{
		DB: Database{
			Host:         "localhost",
			Port:         5432,
			Username:     "postgres",
			Password:     "postgres",
			DatabaseName: "banana",
		},
		WebServer: echo.New(),
	}

	return &conf
}
