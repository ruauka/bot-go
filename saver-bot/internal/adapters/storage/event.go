package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"saver-bot/internal/domain/entities"
)

// userAuthStorage auth storage struct.
type eventStorage struct {
	db *sqlx.DB
}

// NewEventStorage auth storage func builder.
func NewEventStorage(db *sqlx.DB) EventStorage {
	return &eventStorage{db: db}
}

func (e *eventStorage) Save(ctx context.Context, event *entities.Event) error {
	query := `INSERT INTO event (date, type, username, telega_id) VALUES (?, ?, ?, ?)`

	if _, err := e.db.ExecContext(ctx, query, event.Date, event.Type, event.Username, event.TelegaID); err != nil {
		return fmt.Errorf("can't save event: %w", err)
	}

	return nil
}

func (e *eventStorage) GetAll(ctx context.Context) ([]entities.Event, error) {
	var events []entities.Event

	query := fmt.Sprintf("SELECT id, date, type, username, telega_id FROM event")
	if err := e.db.SelectContext(ctx, &events, query); err != nil {
		return events, err
	}

	return events, nil
}
