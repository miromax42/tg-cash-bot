package database

import (
	"context"
	"fmt"
	"os"
	"testing"

	"entgo.io/ent/dialect/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo/database/db_container_test"
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

func getDB(t *testing.T) (db *ent.Client, url string) {
	var connectionString string
	if _, ok := os.LookupEnv("CI_PROJECT_ID"); ok {
		user, ok := os.LookupEnv("POSTGRES_USER")
		require.True(t, ok)

		password, ok := os.LookupEnv("POSTGRES_PASSWORD")
		require.True(t, ok)

		host, ok := os.LookupEnv("POSTGRES_HOSTNAME")
		require.True(t, ok)

		connectionString = fmt.Sprintf("postgresql://%s:%s@%s:5432/postgres?sslmode=disable", user, password, host)
	} else {
		container, err := db_container_test.NewTestDatabase()
		require.NoError(t, err)
		connectionString = container.ConnectionString()
	}

	db, err := ent.Open("postgres", connectionString)
	require.NoError(t, err)

	err = db.Schema.Create(context.Background())
	require.NoErrorf(t, err, "connectionString=%q", connectionString)

	return db, connectionString
}

func generator() func() int64 {
	var inc int64

	return func() int64 {
		inc++

		return inc
	}
}
