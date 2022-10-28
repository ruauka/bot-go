package main

import (
	"fmt"
	"log"

	"saver-bot/internal/adapters/queue"
	s "saver-bot/internal/adapters/storage"
	tg "saver-bot/internal/adapters/telegram"
	"saver-bot/internal/domain/usecases"
	"saver-bot/pkg/client/rabbitmq"
	"saver-bot/pkg/client/sqlite"
	"saver-bot/pkg/client/telegram"
)

func main() {
	bot, updates, err := telegram.NewTelegramBot()
	if err != nil {
		log.Fatalf("failed to connect telegram: %s", err.Error())
	}

	conn, err := rabbitmq.NewRabbitMQConnect()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect to RabbitMQ: %s", err.Error()))
	}

	defer func() { _ = conn.Close() }()

	mq := queue.NewQueue(bot, conn)

	db, err := sqlite.NewSqliteConnect()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to init db: %s", err.Error()))
	}

	storage := s.NewStorage(db)

	usecase := usecases.NewUsecases(storage, bot, mq)

	app := tg.NewApp(usecase)
	app.Start(updates)
}
