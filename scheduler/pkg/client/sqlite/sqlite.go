package sqlite

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	// import sqlite driver.
	_ "github.com/mattn/go-sqlite3"
)

// NewSqliteConnect creates new SQLite storage.
func NewSqliteConnect() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "../saver-bot/storage.db")
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	log.Println("Connect to db: ok")

	return db, nil
}
