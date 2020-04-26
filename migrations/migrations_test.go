package migrations_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/bbengfort/catena/migrations"
	"github.com/stretchr/testify/require"
)

// Test the migrations package and code -- does not run database commands.
func TestMigrations(t *testing.T) {
	names, err := filepath.Glob("*.sql")
	require.NoError(t, err)
	require.Equal(t, len(names), Num(), "number of expected migrations doesn't match, have you run go generate?")
}

// Test the migrations themselves -- runs database commands.
func TestDatabase(t *testing.T) {
	dburl := os.Getenv("CATENA_TEST_DATABASE")
	if dburl == "" {
		t.Skip("no test database available, set $CATENA_TEST_DATABASE")
	}
}
