package main

import (
	"fmt"
	"log"

	q "bot/internal/adapters/queue"
	s "bot/internal/adapters/storage"
	tg "bot/internal/adapters/telegram"
	"bot/internal/config"
	"bot/internal/domain/usecases"
	"bot/pkg/client/postgres"
	"bot/pkg/client/rabbitmq"
	"bot/pkg/client/telegram"
)

func main() {
	// config create
	cfg := config.GetConfig()

	bot, updates, err := telegram.NewTelegramBot(cfg)
	if err != nil {
		log.Fatalf("failed to connect telegram: %s", err.Error())
	}

	mq, err := rabbitmq.NewRabbitMQConnect(cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect to RabbitMQ: %s", err.Error()))
	}

	defer func() { _ = mq.Close() }()

	queue := q.NewQueue(bot, mq)

	//db, err := sqlite.NewSqliteConnect()
	db, err := postgres.NewPostgresConnect(cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to init db: %s", err.Error()))
	}

	storage := s.NewStorage(db)

	usecase := usecases.NewUsecases(storage, queue)

	app := tg.NewApp(usecase, bot)
	app.Start(updates)
}
