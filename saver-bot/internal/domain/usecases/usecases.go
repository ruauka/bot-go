package usecases

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"saver-bot/internal/adapters/queue"
	"saver-bot/internal/adapters/storage"
)

type StorageUsecase interface {
	CommandHandle(update *tg.Update)
	MenuButtonsHandle(update *tg.Update, button string)
	ChatStateHandle(update *tg.Update, state *State)
	MakeResponse(update *tg.Update, text string)
}

type QueueUsecase interface {
	QueueChanListen()
}

type Usecases struct {
	Storage StorageUsecase
	Queue   QueueUsecase
}

func NewUsecases(storage storage.Storage, bot *tg.BotAPI, queue queue.Queue) *Usecases {
	return &Usecases{
		Storage: NewStorageUsecase(storage, bot),
		Queue:   NewQueueUsecase(queue),
	}
}
