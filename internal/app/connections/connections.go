package connections

import (
	// "log"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)


var db *sql.DB


func ConnectDB(databaseURL string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for {
		db, err = sql.Open("mysql", databaseURL)
		if err == nil {
			break
		}
		log.Println("Waiting for database...")
		time.Sleep(2 * time.Second)
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
