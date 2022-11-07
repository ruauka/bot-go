package usecases

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"bot/internal/adapters/queue"
	"bot/internal/adapters/storage"
)

type StorageUsecase interface {
	CommandHandle(update *tg.Update)
	ButtonsHandle(update *tg.Update, button string)
	ChatStateHandle(update *tg.Update, state *State)
	MakeResponse(update *tg.Update, text string)
	IsChatState(userID int64) *State
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
