package telegram

import (
	"fmt"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivangurin/cbrf-go"
)

func (a *App) Currency(update *tg.Update) {
	yesterday := make([]float64, 2, 2)
	today := make([]float64, 2, 2)
	tomorrow := make([]float64, 2, 2)

	var err error

	yesterday[0], err = a.usecases.Currency.Get(cbrf.CurrencyUSD, time.Now().Add(time.Hour*-24))
	yesterday[1], err = a.usecases.Currency.Get(cbrf.CurrencyEUR, time.Now().Add(time.Hour*-24))
	today[0], err = a.usecases.Currency.Get(cbrf.CurrencyUSD, time.Now())
	today[1], err = a.usecases.Currency.Get(cbrf.CurrencyEUR, time.Now())
	tomorrow[0], err = a.usecases.Currency.Get(cbrf.CurrencyUSD, time.Now().Add(time.Hour*24))
	tomorrow[1], err = a.usecases.Currency.Get(cbrf.CurrencyEUR, time.Now().Add(time.Hour*24))

	if err != nil {
		a.MakeResponse(update, CbProblem)
	}

	if today[0] == tomorrow[0] && today[1] == tomorrow[1] {
		a.MakeResponse(update, fmt.Sprintf("Вчера:\n\n💵 Доллар: %.4f\n💶 Евро: %.4f", yesterday[0], yesterday[1]))
		a.MakeResponse(update, fmt.Sprintf("Сейчас:\n\n💵 Доллар: %.4f\n💶 Евро: %.4f", today[0], today[1]))
		a.MakeResponse(update, "На завтра:\n\n🤷‍♂ Курс пока не назначен")
		return
	}

	a.MakeResponse(update, fmt.Sprintf("Вчера:\n\n💵 Доллар: %.4f\n💶 Евро: %.4f", yesterday[0], yesterday[1]))
	a.MakeResponse(update, fmt.Sprintf("Сейчас:\n\n💵 Доллар: %.4f\n💶 Евро: %.4f", today[0], today[1]))
	a.MakeResponse(update, fmt.Sprintf("На завтра:\n\n💵 Доллар: %.4f\n💶 Евро: %.4f", tomorrow[0], tomorrow[1]))
}
