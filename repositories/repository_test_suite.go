package repositories

import (
	"banana-back/migrations"
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"testing"
	"time"
)

func TestRepository(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepositorySuite))
}

// RepositorySuite is a suite for the repositories.
type RepositorySuite struct {
	suite.Suite
	client      *bun.DB
	pgContainer *postgres.PostgresContainer
}

// SetupSuite is called once before the execution of the suite.
func (s *RepositorySuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// init container
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase("test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	s.Require().NoError(err)
	s.pgContainer = pgContainer
	connectionString, err := s.pgContainer.ConnectionString(ctx)
	connectionString += "sslmode=disable"
	s.Require().NoError(err)

	// init pg client
	pool, err := pgxpool.New(ctx, connectionString)
	s.Require().NoError(err)
	s.client = bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New())

	// migrate
	migrationsFS, err := iofs.New(migrations.FS, ".")
	s.Require().NoError(err)
	migration, err := migrate.NewWithSourceInstance("iofs", migrationsFS, connectionString)
	s.Require().NoError(err)
	s.Require().NoError(migration.Up())
	sourceErr, databaseErr := migration.Close()
	s.Require().NoError(sourceErr)
	s.Require().NoError(databaseErr)

	// snapshot
	s.Require().NoError(s.pgContainer.Snapshot(ctx, postgres.WithSnapshotName("empty")))
}

// TearDownSuite is called after all tests in this suite have run.
func (s *RepositorySuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Require().NoError(s.client.Close())
	s.Require().NoError(s.pgContainer.Terminate(ctx))
}

// SetupSubTest is called before each subtest.
func (s *RepositorySuite) SetupSubTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Require().NoError(s.pgContainer.Restore(ctx))
}
