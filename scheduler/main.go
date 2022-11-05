package main

import (
	"fmt"
	"log"

	q "scheduler/internal/adapters/queue"
	s "scheduler/internal/adapters/storage"
	"scheduler/pkg/client/postgres"
	"scheduler/pkg/client/rabbitmq"
)

func main() {
	db, err := postgres.NewPostgresConnect()
	if err != nil {
		fmt.Println(err)
	}

	storage := s.NewEventStorage(db)

	mq, err := rabbitmq.NewRabbitMQConnect()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect to RabbitMQ: %s", err.Error()))
	}
	defer func() { _ = mq.Close() }()

	app := q.NewApp(mq, storage)
	app.Start()
}
