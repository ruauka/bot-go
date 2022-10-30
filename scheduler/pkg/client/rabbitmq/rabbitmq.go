package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnect() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@queue:5672/")
	if err != nil {
		return nil, err
	}

	log.Println("Connect to queue: ok")

	return conn, nil
}
