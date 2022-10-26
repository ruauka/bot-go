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

func (e *eventStorage) SaveManic(ctx context.Context, m *entities.Manic) error {
	q := `INSERT INTO manic (event_date) VALUES (?)`

	if _, err := e.db.ExecContext(ctx, q, m.Date); err != nil {
		return fmt.Errorf("can't save manic: %w", err)
	}

	return nil
}

func (e *eventStorage) SaveMassage(ctx context.Context, m *entities.Massage) error {
	q := `INSERT INTO massage (event_date) VALUES (?)`

	if _, err := e.db.ExecContext(ctx, q, m.Date); err != nil {
		return fmt.Errorf("can't save massage: %w", err)
	}

	return nil
}

func (e *eventStorage) GetAllEvents(ctx context.Context) ([]entities.Manic, []entities.Massage, error) {
	var manics []entities.Manic

	query := fmt.Sprintf("SELECT id, event_date FROM manic")
	if err := e.db.SelectContext(ctx, &manics, query); err != nil {
		return nil, nil, err
	}

	var massages []entities.Massage

	query = fmt.Sprintf("SELECT id, event_date FROM massage")
	if err := e.db.SelectContext(ctx, &massages, query); err != nil {
		return nil, nil, err
	}

	return manics, massages, nil
}
