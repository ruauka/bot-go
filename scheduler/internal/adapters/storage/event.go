package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"scheduler/internal/entities"
)

type eventStorage struct {
	db *sqlx.DB
}

// NewEventStorage auth storage func builder.
func NewEventStorage(db *sqlx.DB) Storage {
	return &eventStorage{db: db}
}

func (e *eventStorage) GetAll(ctx context.Context) []entities.Event {
	var events []entities.Event

	query := fmt.Sprintf("SELECT id, date, type, username, telega_id FROM event")
	if err := e.db.SelectContext(ctx, &events, query); err != nil {
		return nil
	}

	return events
}

func (e *eventStorage) Delete(ctx context.Context, eventID string) {
	query := fmt.Sprintf("DELETE FROM event WHERE id=$1")
	if _, err := e.db.ExecContext(ctx, query, eventID); err != nil {
		log.Fatalf("problem delete event")
	}
}
