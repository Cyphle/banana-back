// Package migrations provides the migrations for the database.
package migrations

import "embed"

// FS embeds the migrations directory.
//
//go:embed *.sql
var FS embed.FS

// To add migration file run in CLI: migrate create -ext sql -dir migrations -seq <migration name>
