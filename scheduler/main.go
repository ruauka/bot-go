package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"

	"scheduler/pkg/client/sqlite"
)

var ReqBodyBytes = new(bytes.Buffer)

type Event struct {
	ID       string `db:"id"`
	Date     string `db:"date"`
	Type     string `db:"type"`
	Username string `db:"username"`
	TelegaID int64  `db:"telega_id"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Layout - шаблон формата даты для изменения типа string -> time.Time.
const Layout = "02.01.2006 15:04"

// ConvertApplicationDate - преобразование входящий даты-строки в дату-time.Time.
func ConvertApplicationDate(applicationDate string) time.Time {
	CurrentDate, err := time.Parse(Layout, applicationDate)
	if err != nil {
		return time.Time{}
	}

	return CurrentDate
}

func ChooseUpcomingEvent(events []Event, db *sqlx.DB) Event {

	for _, e := range events {
		if ((time.Since(ConvertApplicationDate(e.Date)).Minutes() + 180) * -1) < 60 {
			DeleteEvent(context.Background(), db, e.ID)
			return e
		}
	}

	return Event{}
}

func DeleteEvent(ctx context.Context, db *sqlx.DB, eventID string) {
	query := fmt.Sprintf("DELETE FROM event WHERE id=$1")
	if _, err := db.ExecContext(ctx, query, eventID); err != nil {
		log.Fatalf("problem delete event")
	}
}

func GetAll(ctx context.Context, db *sqlx.DB) []Event {
	var events []Event

	query := fmt.Sprintf("SELECT id, date, type, username, telega_id FROM event")
	if err := db.SelectContext(ctx, &events, query); err != nil {
		return nil
	}

	return events
}

func SendToQueue(ctx context.Context, ch *amqp.Channel, q amqp.Queue, event Event) {

	ReqBodyBytes.Reset()

	json.NewEncoder(ReqBodyBytes).Encode(event)

	err := ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        ReqBodyBytes.Bytes(),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent a message: %v\n", event)
}

func main() {

	db, err := sqlite.NewSqliteConnect()
	if err != nil {
		fmt.Println(err)
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"my_queue_1", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(time.Second * 5)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			event1 := GetAll(context.Background(), db)
			event := ChooseUpcomingEvent(event1, db)
			if event == (Event{}) {
				log.Println("nothing for rabbit...")
				continue
			}
			SendToQueue(ctx, ch, q, event)
		}
	}

}
