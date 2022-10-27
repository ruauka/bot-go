package main

import (
	"fmt"
	"log"

	q "scheduler/internal/adapters/queue"
	s "scheduler/internal/adapters/storage"
	"scheduler/pkg/client/rabbitmq"
	"scheduler/pkg/client/sqlite"
)

func main() {
	db, err := sqlite.NewSqliteConnect()
	if err != nil {
		fmt.Println(err)
	}

	storage := s.NewEventStorage(db)

	queueConn, err := rabbitmq.NewRabbitMQConnect()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect to RabbitMQ: %s", err.Error()))
	}
	defer func() { _ = queueConn.Close() }()

	queue := q.NewEventQueue(queueConn, storage)
	queue.Start()
}
