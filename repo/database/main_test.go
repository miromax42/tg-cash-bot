package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
)

type PostgresTestSuite struct {
	suite.Suite
	c                *ent.Client
	connectionString string
	generatorID      func() int64
}

func (s *PostgresTestSuite) SetupSuite() {
	s.c, s.connectionString = getDB(s.T())
	s.generatorID = generator()
}

func TestPersonalSettingsSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PostgresTestSuite))
}

func (s *PostgresTestSuite) applyFixture(filePath string, values map[string]interface{}) {
	const dialect = "postgres"

	sqlDB, err := sql.Open(dialect, s.connectionString)
	require.NoError(s.T(), err)
	defer sqlDB.Close()

	fixtures, err := testfixtures.New(
		testfixtures.Template(),
		testfixtures.TemplateData(values),
		testfixtures.Database(sqlDB.DB()),
		testfixtures.Dialect(dialect),
		testfixtures.FilesMultiTables(filePath),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	require.NoError(s.T(), err)
	require.NoError(s.T(), fixtures.Load())
}

func getDB(t *testing.T) (*ent.Client, string) {
	var (
		container        = NewTestDatabase(t)
		connectionString = container.ConnectionString(t)
	)
	db, err := ent.Open("postgres", connectionString)
	require.NoError(t, err)

	err = db.Schema.Create(context.Background())
	require.NoError(t, err)

	return db, connectionString
}

func generator() func() int64 {
	var inc int64
	return func() int64 {
		inc++
		return inc
	}
}

type TestDatabase struct {
	instance testcontainers.Container
}

func NewTestDatabase(t *testing.T) *TestDatabase {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{"5432/tcp"},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	return &TestDatabase{
		instance: postgres,
	}
}

func (db *TestDatabase) Port(t *testing.T) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	p, err := db.instance.MappedPort(ctx, "5432")
	require.NoError(t, err)
	return p.Int()
}

func (db *TestDatabase) ConnectionString(t *testing.T) string {
	return fmt.Sprintf("postgresql://postgres:postgres@127.0.0.1:%d/postgres?sslmode=disable", db.Port(t))
}

func (db *TestDatabase) Close(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	require.NoError(t, db.instance.Terminate(ctx))
}
