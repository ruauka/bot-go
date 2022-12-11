package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"bot/internal/config"
)

func NewRabbitMQConnect(cfg *config.Config) (*amqp.Connection, error) {
	var host string

	if cfg.Level == "dev" {
		host = "localhost"
	} else {
		host = "queue"
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:5672/", host))
	if err != nil {
		return nil, err
	}

	log.Println("Connect to queue: ok")

	return conn, nil
}
