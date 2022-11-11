package telegram

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	q "bot/internal/adapters/queue"
)

func (a *App) Forecast(update *tg.Update) {

	a.MakeResponse(
		update,
		fmt.Sprintf(
			"Сейчас:\n\n%s\nТемпература %d°C\nОщущается как %d°C\n",
			q.WeatherCache["now"]["condition"],
			q.WeatherCache["now"]["temp"],
			q.WeatherCache["now"]["feelsLike"],
		),
	)

	a.MakeResponse(
		update,
		fmt.Sprintf(
			"%s:\n\n%s\nТемпература %d°C\nБудет ощущаться как %d°C\n",
			q.WeatherCache["part1"]["datePart"],
			q.WeatherCache["part1"]["condition"],
			q.WeatherCache["part1"]["temp"],
			q.WeatherCache["part1"]["feelsLike"],
		),
	)

	a.MakeResponse(
		update,
		fmt.Sprintf(
			"%s:\n\n%s\nТемпература %d°C\nБудет ощущаться как %d°C\n",
			q.WeatherCache["part2"]["datePart"],
			q.WeatherCache["part2"]["condition"],
			q.WeatherCache["part2"]["temp"],
			q.WeatherCache["part2"]["feelsLike"],
		),
	)
}
