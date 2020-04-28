package migrations_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	. "github.com/bbengfort/catena/migrations"
	"github.com/stretchr/testify/require"

	// use postgres for the test database
	_ "github.com/lib/pq"
)

// Test the migrations package and code -- does not run database commands.
func TestMigrations(t *testing.T) {
	// Must be same number of migrations as .sql files
	names, err := filepath.Glob("*.sql")
	require.NoError(t, err)
	require.Equal(t, len(names), Num(), "number of expected migrations doesn't match, have you run go generate?")

	// Migration 0 should be the migration schema migration (it's special)
	m, err := Revision(0, nil)
	require.NoError(t, err)
	require.Equal(t, "migrations schema", m.Name)

	// Hopefully we don't have this many migrations ...
	_, err = Revision(9999999, nil)
	require.Error(t, err)
}

// Test the migrations themselves -- runs database commands.
func TestDatabase(t *testing.T) {
	// postgres://localhost:5432/catena_test?sslmode=disable
	dburl := os.Getenv("CATENA_TEST_DATABASE")
	if dburl == "" {
		t.Skip("no test database available, set $CATENA_TEST_DATABASE")
	}

	conn, err := sql.Open("postgres", dburl)
	require.NoError(t, err, "could not connect to database")

	err = Refresh(conn)
	require.NoError(t, err)
}
