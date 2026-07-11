package authpg_test

import (
	"os"
	"testing"

	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var testDB *sqlx.DB

// TestMain runs before all tests and sets up a test database
func TestMain(m *testing.M) {
	var cleanup func()
	testDB, cleanup = testutil.NewTestPg()
	exitCode := m.Run()
	cleanup()
	os.Exit(exitCode)
}
