package testutil

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	postgresContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NewTestPg creates a new test database connection.
// A cleanup function is returned that should be called after the test is done.
func NewTestPg() (*sqlx.DB, func()) {

	// Start the postgres test container
	ctx := context.Background()
	container, err := postgresContainer.Run(ctx, "postgres:18-alpine",
		postgresContainer.WithDatabase("critiquefi_test"),
		postgresContainer.WithUsername("postgres"),
		postgresContainer.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(60*time.Second),
			),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgresContainer container: %v", err)
	}

	// Connect to the test database and run migrations
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}
	mig, err := migrate.New("file://../../../../migrations", connStr)
	if err != nil {
		log.Fatalf("error connecting to database for migrations: %v", err)
	}
	if err := mig.Up(); err != nil {
		log.Fatal(err)
	}

	// Create the test database connection
	testDB, err := postgres.NewDB(ctx, postgres.DBConfig{
		URL:             connStr,
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	})
	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}

	return testDB, func() {
		if err := testDB.Close(); err != nil {
			log.Fatalf("failed to close test db: %v", err)
		}

		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate postgresContainer container: %v", err)
		}
	}
}

// PrepareDB prepares the DB for testing by clearing and seeding tables.
func PrepareDB(t *testing.T, db *sqlx.DB) {
	t.Helper()

	// Reset and seed tables
	_, err := db.ExecContext(context.Background(), `
        TRUNCATE TABLE films, refresh_tokens, users RESTART IDENTITY CASCADE;

        INSERT INTO users (id, email, display_name, name, password_hash, is_admin, is_active) VALUES
        (1, 'user@critiquefi.com', 'test.user', 'Test User', 'password', false, true),
        (2, 'admin@critiquefi.com', 'admin.user', 'Test Admin', 'password', true, true),
        (3, 'deactivated@critiquefi.com', 'deactivated.user', 'Deactivated User', 'password', false, false);

		INSERT INTO refresh_tokens (token_hash, user_id, user_agent, expires_at) VALUES
		('test-token-hash', 1, 'test-user-agent', '2027-01-01 00:00:00');

		INSERT INTO films (id, film_type, title, description, release_date, runtime_minutes, external_references, created_by, updated_by) VALUES
		(1, 'FEATURE FILM', 'Fight Club', 'A movie about a fight between two characters', '2022-01-01', 95, '[{"name": "IMDB", "url": "https://www.imdb.com/title/tt0137523/"}]', 2, 2);
		
        SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT MAX(id) FROM users));
        SELECT setval(pg_get_serial_sequence('films', 'id'), (SELECT MAX(id) FROM films));
    `)
	if err != nil {
		log.Fatalf("failed to seed tables: %v", err)
	}
}
