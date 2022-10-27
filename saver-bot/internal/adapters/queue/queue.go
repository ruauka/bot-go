package queue

import (
	"encoding/json"
	"fmt"
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	amqp "github.com/rabbitmq/amqp091-go"

	"saver-bot/internal/domain/entities"
	"saver-bot/internal/domain/usecases"
)

type Queue interface {
	QueueChanListen()
}

type queue struct {
	bot  *tg.BotAPI
	conn *amqp.Connection
}

func NewQueue(bot *tg.BotAPI, mq *amqp.Connection) Queue {
	return &queue{
		bot:  bot,
		conn: mq,
	}
}

func (mq *queue) QueueChanListen() {
	ch, err := mq.conn.Channel()
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

	messages, err := ch.Consume(
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

	go func() {
		for d := range messages {
			var event entities.Event

			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Fatalf("queue unmarshal message err %s", err)
			}

			if event.Type == usecases.Manic {
				msg := tg.NewMessage(event.TelegaID,
					fmt.Sprintf("Напоминаю. У тебя сегодня %s в %s 💅", "маникюр", event.Date[11:]),
				)
				mq.bot.Send(msg)
			}

			if event.Type == usecases.Massage {
				msg := tg.NewMessage(event.TelegaID,
					fmt.Sprintf("Напоминаю. У тебя сегодня %s в %s 💆‍♀", "массаж", event.Date[11:]),
				)
				mq.bot.Send(msg)
			}

			log.Printf("Received a message: %s", d.Body)
		}
	}()

	<-forever
}
