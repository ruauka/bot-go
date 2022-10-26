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
func NewEventStorage(db *sqlx.DB) Storage {
	return &eventStorage{db: db}
}

func (e *eventStorage) Save(ctx context.Context, m *entities.Event) error {
	q := `INSERT INTO event (date, type, telega_id) VALUES (?, ?, ?)`

	if _, err := e.db.ExecContext(ctx, q, m.Date, m.Type, m.TelegaID); err != nil {
		return fmt.Errorf("can't save event: %w", err)
	}

	return nil
}

func (e *eventStorage) GetAll(ctx context.Context) ([]entities.Event, error) {
	var events []entities.Event

	query := fmt.Sprintf("SELECT id, date, type, telega_id FROM event")
	if err := e.db.SelectContext(ctx, &events, query); err != nil {
		return events, err
	}

	return events, nil
}
