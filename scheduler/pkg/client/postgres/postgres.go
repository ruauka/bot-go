package postgres

import (
	"fmt"
	"log"

	// import db migrations engine.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// NewPostgresConnect create connect with DB.
func NewPostgresConnect() (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		"pg",
		"pass",
		"database",
		"5432",
		"crud",
		"disable",
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successful connection to the database...")

	return db, nil
}
