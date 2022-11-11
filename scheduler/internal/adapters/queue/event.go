package queue

import (
	"bytes"
	"context"
	"time"

	"scheduler/internal/entities"
)

var (
	eventBytesBuff = new(bytes.Buffer)
	deduct         = 3
)

// Layout - шаблон формата даты для изменения типа string -> time.Time.
const Layout = "02.01.2006 15:04"

func (a *App) ChooseUpcomingEvent(events []entities.Event) entities.Event {
	for _, event := range events {
		if ((time.Since(convertDate(event.Date)).Minutes() + 180) * -1) < 60 {
			a.storage.Delete(context.Background(), event.ID)
			event.ReminderStatus = 2

			return event
		}
	}

	return entities.Event{}
}

func timeCheckEvent(allEvents []entities.Event) []entities.Event {
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

// convertDate - преобразование входящий даты-строки в дату-time.Time.
func convertDate(date string) time.Time {
	parseDate, err := time.Parse(Layout, date)
	if err != nil {
		return time.Time{}
	}

	return parseDate
}

func morningCheck(date string) bool {
	currentDate := convertDate(date)
	if currentDate.Year() == time.Now().Year() &&
		currentDate.Month() == time.Now().Month() &&
		currentDate.Day() == time.Now().Day() &&
		time.Now().Hour() == 7-deduct && time.Now().Minute() == 40 &&
		(time.Now().Second() == 0 ||
			time.Now().Second() == 1 ||
			time.Now().Second() == 2 ||
			time.Now().Second() == 3 ||
			time.Now().Second() == 4) {
		return true
	}

	return false
}

func eveningCheck(date string) bool {
	currentDate := convertDate(date)
	if currentDate.Year() == time.Now().Year() &&
		currentDate.Month() == time.Now().Month() &&
		currentDate.Day() == time.Now().Day()+1 &&
		time.Now().Hour() == 21-deduct && time.Now().Minute() == 30 &&
		(time.Now().Second() == 0 ||
			time.Now().Second() == 1 ||
			time.Now().Second() == 2 ||
			time.Now().Second() == 3 ||
			time.Now().Second() == 4) {
		return true
	}

	return false
}
