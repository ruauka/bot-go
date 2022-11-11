package queue

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"bot/internal/domain/entities"
)

var (
	ManicSmall   = "маникюр"
	MassageSmall = "массаж"
	SportSmall   = "спорт"
	MeetingSmall = "встреча"
	Massage      = "Массаж"
	Manic        = "Маникюр"
	Sport        = "Спорт"
	Meeting      = "Встреча"
)

func eventCheck(event entities.Event) tg.MessageConfig {
	if event.Type == Manic {
		switch event.ReminderStatus {
		case 0:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Доброе утро. У тебя сегодня %s в %s 💅", ManicSmall, event.Date[11:]),
			)
		case 1:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Напоминаю. У тебя завтра %s в %s 💅", ManicSmall, event.Date[11:]),
			)
		case 2:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Напоминаю. У тебя через час %s в %s 💅", ManicSmall, event.Date[11:]),
			)
		}
	}

	if event.Type == Massage {
		switch event.ReminderStatus {
		case 0:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Доброе утро. У тебя сегодня %s в %s 💆‍♀", MassageSmall, event.Date[11:]),
			)
		case 1:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Напоминаю. У тебя завтра %s в %s 💆‍♀", MassageSmall, event.Date[11:]),
			)
		case 2:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Напоминаю. У тебя через час %s в %s 💆‍♀", MassageSmall, event.Date[11:]),
			)
		}
	}

	if event.Type == Sport {
		switch event.ReminderStatus {
		case 0:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Доброе утро. У тебя сегодня %s в %s 🏃‍♀", SportSmall, event.Date[11:]),
			)
		case 1:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Напоминаю. У тебя завтра %s в %s 🏃‍♀", SportSmall, event.Date[11:]),
			)
		case 2:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Напоминаю. У тебя через час %s в %s 🏃‍♀", SportSmall, event.Date[11:]),
			)
		}
	}

	if event.Type == Meeting {
		switch event.ReminderStatus {
		case 0:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Доброе утро. У тебя сегодня %s в %s 🗓", MeetingSmall, event.Date[11:]),
			)
		case 1:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Напоминаю. У тебя завтра %s в %s 🗓", MeetingSmall, event.Date[11:]),
			)
		case 2:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Напоминаю. У тебя через час %s в %s 🗓", MeetingSmall, event.Date[11:]),
			)
		}
	}

	return tg.MessageConfig{}
}
