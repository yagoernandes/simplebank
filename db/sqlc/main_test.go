package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

//15:00

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simplebank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	// Connect to the test database
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// Initialize the testQueries object
	testQueries = New(testDB)

	// Run the tests
	os.Exit(m.Run())
}
