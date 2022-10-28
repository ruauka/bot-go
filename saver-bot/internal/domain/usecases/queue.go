package usecases

import (
	"saver-bot/internal/adapters/queue"
)

type queueUsecase struct {
	queue queue.Queue
}

func NewQueueUsecase(queue queue.Queue) QueueUsecase {
	return &queueUsecase{
		queue: queue,
	}
}

func (q *queueUsecase) QueueChanListen() {
	q.queue.QueueChanListen()
}
