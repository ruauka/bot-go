package sqlite

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	// import sqlite driver.
	_ "github.com/mattn/go-sqlite3"
)

// NewSqliteConnect creates new SQLite storage.
func NewSqliteConnect() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "storage.db")
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	log.Println("Connect to db: ok")

	if err := Init(context.Background(), db); err != nil {
		return nil, err
	}

	return db, nil
}

func Init(ctx context.Context, db *sqlx.DB) error {
	q := `
		CREATE TABLE IF NOT EXISTS event (
		    id INTEGER PRIMARY KEY,
		    date varchar(255),
		    type varchar(255),
		    telega_id varchar(255)
		);`

	_, err := db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	log.Println("Create db schemas: ok")

	return nil
}
