package models

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// The getTestEnvVariables() helper reads the .env file and returns the
// requested key value
func getEnvVariables(t *testing.T, key string) string {
	// Find the basepath of the current root of project.
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Unable to load .env directory!")
	}

	basepath := filepath.Dir(file)

	// Use the godotenv.Load() method to load the env files in base
	// directory.
	err := godotenv.Load(filepath.Join(basepath, "../..", "/.env"))
	if err != nil {
		t.Fatal("No .env file found!")
	}

	return os.Getenv(key)
}

func newTestDB(t *testing.T) *sql.DB {
	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contain multiple SQL statements, we need to
	// use the "multiStatement=true" parameter in our DSN. This instructs our
	// PostgreSQL database driver to support executing multiple SQL statements
	// in our db.Exec() call.
	// dsn := flag.String("dsn", getTestEnvVariables(t, "TEST_DATABASE_URL"), "PostgresSQL data source name")
	// db, err := sql.Open("postgres", dsn)
	db, err := sql.Open("postgres", getEnvVariables(t, "TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script from file and execute the statements.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Use the t.Cleanup() to register a func *which will automatically be
	// called by Go when the current test/ sub-test, which calls newTestDB(),
	// has finished. In this func we read and execute the teardown script, and
	// close the database connection pool.
	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	// Return the database connection pool.
	return db
}
