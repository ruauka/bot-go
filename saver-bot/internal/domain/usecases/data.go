package usecases

import (
	"regexp"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	StateDate = iota
	StateTime
)

const (
	SignDate        = "Укажи дату. Формат: dd.mm.yyyy 🗓"
	SignTime        = "Укажи время. Формат: hh:mm 🕔"
	SaveUpdate      = "Cохранил. Напомню тебе 👌"
	DBProblem       = "Проблема с БД ❌"
	WrongDateFormat = "Некорректный формат даты ❌🗓"
	WrongTimeFormat = "Некорректный формат времени ❌🕔"
	EmptyManic      = "Пока у тебя нет записей на маникюр 🤷‍♀"
	EmptyMassage    = "Пока у тебя нет записей на массаж 🤷‍♀"
)

type State struct {
	State int // 0 - date, 1 - time
	Date  string
	Time  string
}

var (
	dateRe = regexp.MustCompile(`^\s*(3[01]|[12][0-9]|0?[1-9])\.(1[012]|0?[1-9])\.((?:19|20)\d{2})\s*$`)
	timeRe = regexp.MustCompile(`^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$`)
)

var (
	MassageState = NewMassageState()
	ManicState   = NewManicState()

	MainMenu = tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("💆‍♀ Массаж"),
			tg.NewKeyboardButton("💅 Маникюр"),
		),
	)
)

func NewMassageState() map[int64]*State {
	return make(map[int64]*State)
}

func NewManicState() map[int64]*State {
	return make(map[int64]*State)
}
