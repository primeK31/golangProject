package connections

import (
	"log"
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

	// Create users table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		balance INT NOT NULL DEFAULT 0
	)`)
	if err != nil {
		db.Close()
		log.Printf("Error creating table: %v", err)
		return nil, err
	}

	return db, nil
}


func CloseDB() {
	if db != nil {
		db.Close()
	}
}
