package queue

import (
	"encoding/json"
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	amqp "github.com/rabbitmq/amqp091-go"

	"bot/internal/domain/entities"
)

type Queue interface {
	QueueChanListen()
}

type queue struct {
	bot             *tg.BotAPI
	ch              *amqp.Channel
	event, forecast <-chan amqp.Delivery
}

func NewQueue(bot *tg.BotAPI, conn *amqp.Connection) Queue {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %s", err.Error())
	}

	defer log.Println("Connect to 'event' queue channel: ok")
	defer log.Println("Connect to 'forecast' queue channel: ok")

	return &queue{
		bot:      bot,
		ch:       ch,
		event:    NewChannelConnect(ch, "event"),
		forecast: NewChannelConnect(ch, "forecast"),
	}
}

func NewChannelConnect(ch *amqp.Channel, queueName string) <-chan amqp.Delivery {
	declare, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare '%s' queue: %s", queueName, err.Error())
	}

	channel, err := ch.Consume(
		declare.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer in '%s' queue: %s", queueName, err.Error())
	}

	return channel
}

func (mq *queue) QueueChanListen() {
	defer func() { _ = mq.ch.Close() }()

	var forever chan struct{}

	go func() {
		for msg := range mq.event {
			var event entities.Event

			err := json.Unmarshal(msg.Body, &event)
			if err != nil {
				log.Fatalf("event queue unmarshal message err %s", err)
			}

			mq.bot.Send(eventCheck(event))

			//log.Printf("Received a message: %s", msg.Body)
		}
	}()

	go func() {
		for msg := range mq.forecast {
			var weather entities.Weather

			err := json.Unmarshal(msg.Body, &weather)
			if err != nil {
				log.Fatalf("forecast queue unmarshal message err %s", err)
			}

			cacheFill(weather)

			//fmt.Println(WeatherCache)
		}
	}()

	<-forever
}
