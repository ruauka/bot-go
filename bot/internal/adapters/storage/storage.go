package storage

import (
	"context"

	"bot/internal/domain/entities"
)

type Storage interface {
	Save(ctx context.Context, m *entities.Event) error
	GetAll(ctx context.Context) ([]entities.Event, error)
	Remove(ctx context.Context, date, button, username string) error
}
