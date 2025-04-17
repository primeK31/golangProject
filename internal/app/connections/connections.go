package connections

import (
	// "log"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)


var db *sql.DB


func ConnectDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}


func CloseDB() {
	if db != nil {
		db.Close()
	}
}
