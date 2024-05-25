package config

import "fmt"

type Database struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

func (d Database) getConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.Username, d.Password, d.Host, d.Port, d.DatabaseName)
}
