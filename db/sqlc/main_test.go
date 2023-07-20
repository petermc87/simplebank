package db

import(
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

// Declare a global test object.
var testDB *sql.DB

// The main test used for the db connection.
func TestMain(m *testing.M){
	
	var err error 
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil{
		log.Fatal("Cannot connect to the database.", err)
	}
	
	// Pass in the new test object into the test query. This will be available
	// for use in the store_test.
	testQueries = New(testDB)

	os.Exit(m.Run())
}