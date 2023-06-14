package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	DBSource = "postgresql://devuser:admin@localhost:5432/blog_go_db?sslmode=disable"
	DBDriver = "postgres"
)

var testStore Store
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(DBDriver, DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testStore = NewStore(testDB)
	testQueries = New(testDB)

	os.Exit(m.Run())
}
