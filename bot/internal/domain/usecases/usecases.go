package usecases

import (
	"time"

	"github.com/matperez/go-cbr-client"

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

func NewUsecases(storage storage.Storage, queue queue.Queue, CBRFClient cbr.Client) *Usecases {
	return &Usecases{
		Event:    NewEventUsecase(storage),
		Currency: NewCurrencyUsecase(CBRFClient),
		Queue:    NewQueueUsecase(queue),
	}
}
