package main

import (
	"testing"
	"os"

	"database/sql"
)

func Test(t *testing.T) {
	createDB(t)
	database, _ := sql.Open("sqlite3", "database_test.db")
	createUsersTable(database)
}

func createDB(t *testing.T) {
	file, err := os.Create("database_test.db")
	if err != nil {
		t.Error(err.Error() + "\nfunction: createDB")
	}
	file.Close()
}