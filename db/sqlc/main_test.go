package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/util"
)

var testQueries *Queries

// Declare a global test object.
var testDB *sql.DB

// The main test used for the db connection.
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("cannot connect to the database", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database.", err)
	}

	// Pass in the new test object into the test query. This will be available
	// for use in the store_test.
	testQueries = New(testDB)

	os.Exit(m.Run())
}
