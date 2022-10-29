package queue

import (
	"encoding/json"
	"fmt"
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	amqp "github.com/rabbitmq/amqp091-go"

	"saver-bot/internal/domain/entities"
)

var (
	ManicSmall   = "маникюр"
	MassageSmall = "массаж"
	SportSmall   = "Спорт"
	MeetingSmall = "Встреча"
	Massage      = "Массаж"
	Manic        = "Маникюр"
	Sport        = "Спорт"
	Meeting      = "Встреча"
)

type queue struct {
	bot      *tg.BotAPI
	ch       *amqp.Channel
	queue    amqp.Queue
	messages <-chan amqp.Delivery
}

func NewQueue(bot *tg.BotAPI, conn *amqp.Connection) Queue {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %s", err.Error())
	}

	log.Println("Connect to queue channel: ok")

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

	return &queue{
		bot:      bot,
		ch:       ch,
		queue:    q,
		messages: messages,
	}
}

func (mq *queue) QueueChanListen() {
	defer func() { _ = mq.ch.Close() }()

	var forever chan struct{}

	go func() {
		for msg := range mq.messages {
			var event entities.Event

			err := json.Unmarshal(msg.Body, &event)
			if err != nil {
				log.Fatalf("queue unmarshal message err %s", err)
			}

			mq.bot.Send(eventCheck(event))

			log.Printf("Received a message: %s", msg.Body)
		}
	}()

	<-forever
}

func eventCheck(event entities.Event) tg.MessageConfig {
	if event.Type == Manic {
		return tg.NewMessage(event.TelegaID,
			fmt.Sprintf("Напоминаю. У тебя сегодня %s в %s 💅", ManicSmall, event.Date[11:]),
		)
	}

	if event.Type == Massage {
		return tg.NewMessage(event.TelegaID,
			fmt.Sprintf("Напоминаю. У тебя сегодня %s в %s 💆‍♀", MassageSmall, event.Date[11:]),
		)
	}

	if event.Type == Sport {
		return tg.NewMessage(event.TelegaID,
			fmt.Sprintf("Напоминаю. У тебя сегодня %s в %s 🏃‍♀", SportSmall, event.Date[11:]),
		)
	}

	if event.Type == Meeting {
		return tg.NewMessage(event.TelegaID,
			fmt.Sprintf("Напоминаю. У тебя сегодня %s в %s 🗓", MeetingSmall, event.Date[11:]),
		)
	}

	return tg.MessageConfig{}
}
