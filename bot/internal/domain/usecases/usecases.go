package usecases

import (
	"time"

	"bot/internal/adapters/queue"
	"bot/internal/adapters/storage"
	"bot/internal/domain/entities"
)

type EventUsecase interface {
	Save(event *entities.Event) error
	Remove(date, button, username string) error
	GetAll() ([]entities.Event, error)
}

type CurrencyUsecase interface {
	Get(currency string, date time.Time) (float64, error)
}

type QueueUsecase interface {
	QueueChanListen()
}

type Usecases struct {
	Event    EventUsecase
	Currency CurrencyUsecase
	Queue    QueueUsecase
}

func NewUsecases(storage storage.Storage, queue queue.Queue) *Usecases {
	return &Usecases{
		Event:    NewEventUsecase(storage),
		Currency: NewCurrencyUsecase(),
		Queue:    NewQueueUsecase(queue),
	}
}
