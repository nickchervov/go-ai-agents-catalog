//go:build integration

package postgres

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type AgentsRepositoryTestSuite struct {
	suite.Suite
	db        *sqlx.DB
	repo      *PostgresRepository
	container testcontainers.Container
	ctx       context.Context
}

func (s *AgentsRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()

	container, err := postgres.Run(
		s.ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("agent_catalog_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	s.Require().NoError(err, "starting postgres container")

	dsn, err := container.ConnectionString(s.ctx, "sslmode=disable")
	s.Require().NoError(err)

	_, file, _, _ := runtime.Caller(0)
	migrationPath = filepath.Join(filepath.Dir(file), "..", "..", "..", "pkg", "migrations")

	db, err := sqlx.Connect("pgx", dsn)
	s.Require().NoError(err)

	repo, err := NewPostgresRepository(s.ctx, dsn)
	s.Require().NoError(err)

	s.container = container
	s.db = db
	s.repo = repo
}
func (s *AgentsRepositoryTestSuite) TearDownSuite() {
	if s.repo != nil {
		s.repo.Close()
	}
	if s.container != nil {
		err := s.container.Terminate(s.ctx)
		s.Require().NoError(err, "unexpected error: %v", err)
	}
}

func (s *AgentsRepositoryTestSuite) SetupTest() {
	_, err := s.db.Exec("TRUNCATE assistant_run, assistants, categories, users CASCADE")
	s.Require().NoError(err)
}

func TestRunAgentsRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AgentsRepositoryTestSuite))
}
