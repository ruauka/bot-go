package queue

import (
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	QueueChanListen()
}

type queue struct {
	bot    *tg.BotAPI
	mqConn *amqp.Connection
}

func NewQueue(bot *tg.BotAPI, mq *amqp.Connection) Queue {
	return &queue{
		bot:    bot,
		mqConn: mq,
	}
}

func (mq *queue) QueueChanListen() {
	ch, err := mq.mqConn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %s", err.Error())
	}

	log.Println("Connect to queue channel: ok")

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"my_queue_1", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err.Error())
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %s", err.Error())
	}

	var forever chan struct{}

	// 202460681 м
	// 394622071 с

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			msg := tg.NewMessage(394622071, string(d.Body))
			mq.bot.Send(msg)
		}
	}()

	<-forever
}
