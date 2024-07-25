package config

import "github.com/labstack/echo/v4"

type Config struct {
	DB            Database
	WebServer     *echo.Echo
	WebServerPort int
}

func Get() *Config {
	conf := Config{
		DB: Database{
			Host:         "localhost",
			Port:         5433,
			Username:     "postgres",
			Password:     "postgres",
			DatabaseName: "banana",
		},
		WebServer:     echo.New(),
		WebServerPort: 8080,
	}

	return &conf
}
