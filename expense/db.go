package expense

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

var db *sql.DB

func InitDB() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	var err error
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}
	createExpenseTable :=
		`CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			title TEXT,
			amount FLOAT,
			note TEXT,
			tags TEXT[]
		)`
	_, err = db.Exec(createExpenseTable)

	if err != nil {
		log.Fatal("Can't create table", err)
	}
}

func setMockDB(mockDB *sql.DB) {
	db = mockDB
}
