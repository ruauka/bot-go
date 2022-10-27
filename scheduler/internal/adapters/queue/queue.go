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

type EventQueue struct {
	ch      *amqp.Channel
	queue   amqp.Queue
	storage storage.Storage
}

func NewEventQueue(conn *amqp.Connection, storage storage.Storage) *EventQueue {
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

	return &EventQueue{
		ch:      ch,
		queue:   queue,
		storage: storage,
	}
}

func (e *EventQueue) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	defer func() { _ = e.ch.Close() }()

	for {
		select {
		case <-ticker.C:
			events := e.storage.GetAll(context.Background())
			event := e.chooseUpcomingEvent(events)
			if event == (entities.Event{}) {
				log.Println("nothing for rabbit...")
				continue
			}
			e.sendToQueue(ctx, event)
		}
	}
}

func (e *EventQueue) chooseUpcomingEvent(events []entities.Event) entities.Event {
	for _, event := range events {
		if ((time.Since(convertDate(event.Date)).Minutes() + 180) * -1) < 60 {
			e.storage.Delete(context.Background(), event.ID)
			return event
		}
	}

	return entities.Event{}
}

func (e *EventQueue) sendToQueue(ctx context.Context, event entities.Event) {
	ReqBodyBytes.Reset()

	json.NewEncoder(ReqBodyBytes).Encode(event)

	if err := e.ch.PublishWithContext(ctx,
		"",           // exchange
		e.queue.Name, // routing key
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
