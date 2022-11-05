package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"scheduler/internal/adapters/storage"
	"scheduler/internal/entities"
)

type App struct {
	ch      *amqp.Channel
	queue   amqp.Queue
	storage storage.Storage
}

func NewApp(conn *amqp.Connection, storage storage.Storage) *App {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	log.Println("Connect to channel: ok")

	queue, err := ch.QueueDeclare(
		"my_queue_1", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	return &App{
		ch:      ch,
		queue:   queue,
		storage: storage,
	}
}

func (a *App) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	defer func() { _ = a.ch.Close() }()

	for {
		select {
		case <-ticker.C:
			allEvents := a.storage.GetAll(context.Background())

			events := a.timeCheckEvent(allEvents)
			if len(events) != 0 {
				for _, event := range events {
					a.sendToQueue(ctx, event)
				}
			}

			event := a.chooseUpcomingEvent(allEvents)
			if event == (entities.Event{}) {
				log.Println("nothing for rabbit...")
				continue
			}

			a.sendToQueue(ctx, event)
		}
	}
}

func (a *App) timeCheckEvent(allEvents []entities.Event) []entities.Event {
	var events []entities.Event

	for _, event := range allEvents {
		if morningCheck(event.Date) {
			event.ReminderStatus = 0
			events = append(events, event)
		}
		if eveningCheck(event.Date) {
			event.ReminderStatus = 1
			events = append(events, event)
		}
	}

	return events
}

func (a *App) chooseUpcomingEvent(events []entities.Event) entities.Event {
	for _, event := range events {
		if ((time.Since(convertDate(event.Date)).Minutes() + 180) * -1) < 60 {
			a.storage.Delete(context.Background(), event.ID)
			event.ReminderStatus = 2

			return event
		}
	}

	return entities.Event{}
}

func (a *App) sendToQueue(ctx context.Context, event entities.Event) {
	ReqBodyBytes.Reset()

	json.NewEncoder(ReqBodyBytes).Encode(event)

	if err := a.ch.PublishWithContext(ctx,
		"",           // exchange
		a.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        ReqBodyBytes.Bytes(),
		}); err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}

	log.Printf("Sent a message: %v\n", event)
}
