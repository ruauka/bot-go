package usecases

import (
	"context"

	"bot/internal/adapters/storage"
	"bot/internal/domain/entities"
)

type eventUsecase struct {
	storage storage.Storage
}

func NewEventUsecase(storage storage.Storage) EventUsecase {
	return &eventUsecase{
		storage: storage,
	}
}

func (e *eventUsecase) Save(event *entities.Event) error {
	return e.storage.Save(context.Background(), event)
}

func (e *eventUsecase) Remove(date, button, username string) error {
	return e.storage.Remove(context.Background(), date, button, username)
}

func (e *eventUsecase) GetAll() ([]entities.Event, error) {
	return e.storage.GetAll(context.Background())
}
