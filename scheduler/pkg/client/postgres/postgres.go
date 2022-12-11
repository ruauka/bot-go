package postgres

import (
	"fmt"
	"log"

	// import db migrations engine.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"scheduler/internal/config"
)

// NewPostgresConnect create connect with DB.
func NewPostgresConnect(cfg *config.Config) (*sqlx.DB, error) {
	var host string

	if cfg.Level == "dev" {
		host = "localhost"
	} else {
		host = "database"
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PgUser,
		cfg.PgPassword,
		host,
		"5432",
		cfg.PgDbName,
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
