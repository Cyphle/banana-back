package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabase_ShouldGetConnectionString(t *testing.T) {
	db := Database{
		Host:         "localhost",
		Port:         5432,
		Username:     "postgres",
		Password:     "postgres",
		DatabaseName: "postgres",
	}

	connectionString := db.GetConnectionString()

	assert.Equal(
		t,
		"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		connectionString,
	)
}
