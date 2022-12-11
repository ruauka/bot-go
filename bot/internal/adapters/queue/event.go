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
	Massage      = "Массаж"
	Manic        = "Маникюр"
	Sport        = "Спорт"
	Meeting      = "Встреча"
)

func eventCheck(event entities.Event) tg.MessageConfig {
	if event.Type == Manic {
		return reminderRespCreate(event, ManicSmall, "💅")
	}

	if event.Type == Massage {
		return reminderRespCreate(event, MassageSmall, "💆‍♀")
	}

	if event.Type == Sport {
		return reminderRespCreate(event, SportSmall, "🏃‍♀")
	}

	if event.Type == Meeting {
		return reminderRespCreate(event, "", "🗓")
	}

	return tg.MessageConfig{}
}

func reminderRespCreate(event entities.Event, eventName, badge string) tg.MessageConfig {
	if event.Type == Meeting {
		switch event.ReminderStatus {
		case 0:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Доброе утро. У тебя сегодня встреча %s в %s %s", event.Whom, event.Date[11:], badge),
			)
		case 1:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Напоминаю. У тебя завтра встреча %s в %s %s", event.Whom, event.Date[11:], badge),
			)
		case 2:
			return tg.NewMessage(event.TelegaID,
				fmt.Sprintf("Напоминаю. У тебя через час встреча %s в %s %s", event.Whom, event.Date[11:], badge),
			)
		}
	}

	switch event.ReminderStatus {
	case 0:
		return tg.NewMessage(event.TelegaID,
			fmt.Sprintf("Доброе утро. У тебя сегодня %s в %s %s", eventName, event.Date[11:], badge),
		)
	case 1:
		return tg.NewMessage(event.TelegaID,
			fmt.Sprintf("Напоминаю. У тебя завтра %s в %s %s", eventName, event.Date[11:], badge),
		)
	case 2:
		return tg.NewMessage(event.TelegaID,
			fmt.Sprintf("Напоминаю. У тебя через час %s в %s %s", eventName, event.Date[11:], badge),
		)
	}

	return tg.MessageConfig{}
}
