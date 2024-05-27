package config

type Config struct {
	DB Database
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
	}

	return &conf
}
