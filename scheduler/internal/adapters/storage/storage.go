package storage

import (
	"context"

	"scheduler/internal/entities"
)

type Storage interface {
	GetAll(ctx context.Context) []entities.Event
	Delete(ctx context.Context, eventID string)
}
