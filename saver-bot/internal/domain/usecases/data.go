package usecases

import (
	"regexp"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const DatePointLayout = "02.01.2006"
const DatePointTimeLayout = "02.01.2006 15:04"
const TimeDashTimeLayout = "2006-01-02 15:04"

const (
	StateQuestion = iota
	StateDate
	StateTime
)

const (
	StartCmd          = "start"
	MenuCmd           = "menu"
	Massage           = "Массаж"
	Manic             = "Маникюр"
	Sport             = "Спорт"
	Meeting           = "Встреча"
	HelloMsg          = "⬇ Привет, выбери пункт меню"
	MainMenu          = "Главное меню"
	SignDate          = "Шаг [1/2]\n\nУкажи дату. Формат: dd.mm.yyyy 🗓"
	SignTime          = "Шаг [2/2]\n\nУкажи время. Формат: hh:mm 🕔"
	DeleteEvent       = "Укажи дату и время записи.\nФормат: dd.mm.yyyy hh:mm 🗓"
	SaveUpdate        = "Cохранил. Напомню тебе 👌"
	DeleteUpdate      = "Отменил 👌"
	DBProblem         = "Проблема с БД ❌"
	WrongDateFormat   = "Некорректный формат даты ❌🗓"
	WrongTimeFormat   = "Некорректный формат времени ❌🕔"
	DateBeforeNow     = "Этот день уже прошел ❌"
	TimeBeforeNow     = "Это время уже прошло ❌"
	EmptyManic        = "Пока нет записей на маникюр 🤷‍♀"
	EmptyMassage      = "Пока нет записей на массаж 🤷‍♀"
	EmptySport        = "Пока нет записей на спорт 🤷‍♀"
	EmptyMeeting      = "Пока нет встреч 🤷‍♀"
	OtherMessagesPlug = "Ой, давай не сейчас..."
	MashaMenu         = "Чем займемся?"
	MassageQuestion   = "Что делаем с массажем?"
	ManicQuestion     = "Что делаем с маникюром?"
	SportQuestion     = "Что делаем со спортом?"
	MeetingQuestion   = "Что делаем со встречей?"
	EventNotFound     = "Не нашел такого 🤷‍♀"
)

type State struct {
	State      int // 0 - question, 1 - date, 2 - time
	ChatName   string
	DeleteMode bool
	Date       string
	Time       string
}

var (
	dateRe = regexp.MustCompile(`^\s*(3[01]|[12][0-9]|0?[1-9])\.(1[012]|0?[1-9])\.((?:19|20)\d{2})\s*$`)
	timeRe = regexp.MustCompile(`^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$`)
)

var (
	MassageState = NewState()
	ManicState   = NewState()
	SportState   = NewState()
	MeetingState = NewState()

	Chats    = []map[int64]*State{MassageState, ManicState, SportState, MeetingState}
	EventArr = [4]string{Massage, Manic, Sport, Meeting}

	MainMenuButtons = tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("Саша"),
		),
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("Маша"),
		),
		// погода
		// курсы валют
	)

	MashaMenuButtons = tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("💆‍♀ Массаж"),
			tg.NewKeyboardButton("💅 Маникюр"),
		),
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("🏃‍♀ Спорт"),
			tg.NewKeyboardButton("🗓 Встреча"),
		),
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("Все мои записи"),
		),
	)

	OrderButtons = tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("🙋‍♀ Создать"),
		),
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("🙅‍♀ Отменить"),
		),
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("Назад"),
		),
	)

	CancelButton = tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("Отмена"),
		),
	)
)

func NewState() map[int64]*State {
	return make(map[int64]*State)
}
