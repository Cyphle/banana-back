package initializers

import (
	"banana-back/config"
	"banana-back/migrations"
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"log/slog"
	"time"
)

const (
	dbMaxIdleConns    = 1
	dbMaxOpenConns    = 2
	dbConnMaxLifetime = time.Hour
)

func InitDatabase(ctx context.Context, conf config.Database) (*bun.DB, error) {
	pool, err := pgxpool.New(ctx, conf.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	pool.Config().MinConns = dbMaxIdleConns
	pool.Config().MaxConns = dbMaxOpenConns
	pool.Config().MaxConnLifetime = dbConnMaxLifetime

	dbClient := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New())

	if err := doMigrate(conf); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return dbClient, nil
}

// To add migration file run in CLI: migrate create -ext sql -dir migrations -seq <migration name>
func doMigrate(conf config.Database) error {
	migrationsFS, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migration, err := migrate.NewWithSourceInstance("iofs", migrationsFS, conf.GetConnectionString())
	if err != nil {
		return fmt.Errorf("failed to initiate migrations: %w", err)
	}

	log := slog.Default()
	if err != nil {
		log.Error("failed to init database", "err", err)
	}

	if err := migration.Up(); err != nil && err.Error() != "no change" {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
