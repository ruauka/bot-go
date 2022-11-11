package postgres

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
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

	log.Println("Connect to database: ok")

	makeMigrations(dsn)

	return db, nil
}

// makeMigrations - make db migrations.
func makeMigrations(dsn string) {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("No changes in database...\n")
			return
		}
		log.Fatal(err)
	}

	log.Printf("Database migration is done...\n")
}
