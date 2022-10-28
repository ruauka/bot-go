package storage

import (
	"context"

	"saver-bot/internal/domain/entities"
)

type Storage interface {
	Save(ctx context.Context, m *entities.Event) error
	GetAll(ctx context.Context) ([]entities.Event, error)
}
