package queue

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"scheduler/internal/adapters/storage"
	"scheduler/internal/config"
	"scheduler/internal/entities"
)

type App struct {
	storage         storage.Storage
	ch              *amqp.Channel
	event, forecast amqp.Queue
}

func NewApp(conn *amqp.Connection, storage storage.Storage) *App {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	defer log.Println("Connect to 'event' channel: ok")
	defer log.Println("Connect to 'forecast' channel: ok")

	return &App{
		ch:       ch,
		event:    newQueueDeclare(ch, "event"),
		forecast: newQueueDeclare(ch, "forecast"),
		storage:  storage,
	}
}

func newQueueDeclare(ch *amqp.Channel, queueName string) amqp.Queue {
	queue, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Fatalf("failed to declare a %s: %s", queueName, err.Error())
	}

	return queue
}

func (a *App) Start(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventTicker := time.NewTicker(time.Second * 5)
	defer eventTicker.Stop()

	forecastTicker := time.NewTicker(time.Minute * 40)
	defer forecastTicker.Stop()

	defer func() { _ = a.ch.Close() }()

	var once sync.Once
	once.Do(func() {
		a.SendToQueue(ctx, a.forecast.Name, a.YandexForecastCall(cfg))
	})

	for {
		select {
		case <-eventTicker.C:
			allEvents := a.storage.GetAll(context.Background())

			events := timeCheckEvent(allEvents)
			if len(events) != 0 {
				for _, event := range events {
					eventBytesBuff.Reset()
					json.NewEncoder(eventBytesBuff).Encode(event)

					a.SendToQueue(ctx, a.event.Name, eventBytesBuff.Bytes())
				}
			}

			event := a.ChooseUpcomingEvent(allEvents)
			if event == (entities.Event{}) {
				log.Println("nothing for rabbit...")
				continue
			}

			eventBytesBuff.Reset()
			json.NewEncoder(eventBytesBuff).Encode(event)

			a.SendToQueue(ctx, a.event.Name, eventBytesBuff.Bytes())

		case <-forecastTicker.C:
			resp := a.YandexForecastCall(cfg)
			a.SendToQueue(ctx, a.forecast.Name, resp)
		}
	}
}

func (a *App) SendToQueue(ctx context.Context, queueName string, message []byte) {
	if err := a.ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		}); err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}

	//log.Printf("Sent a message data in %s", queueName)
}
