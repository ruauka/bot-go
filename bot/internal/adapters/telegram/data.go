package telegram

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
	StateMeeting
)

const (
	StartCmd                 = "start"
	MenuCmd                  = "menu"
	Massage                  = "Массаж"
	Manic                    = "Маникюр"
	Sport                    = "Спорт"
	Meeting                  = "Встреча"
	HelloMsg                 = "⬇ Привет, выбери пункт меню"
	MainMenu                 = "Главное меню"
	SignDate                 = "Шаг [1/2]\n\nУкажи дату. Формат: dd.mm.yyyy 🗓"
	SignTime                 = "Шаг [2/2]\n\nУкажи время. Формат: hh:mm 🕔"
	MeetingSignDate          = "Шаг [1/3]\n\nУкажи дату. Формат: dd.mm.yyyy 🗓"
	MeetingSignTime          = "Шаг [2/3]\n\nУкажи время. Формат: hh:mm 🕔"
	MeetingSignWithWhom      = "Шаг [3/3]\n\nУкажи с кем встреча.\nФормат: с тем-то 💃"
	DeleteEvent              = "Укажи дату и время записи.\nФормат: dd.mm.yyyy hh:mm 🗓"
	SaveUpdate               = "Cохранил. Напомню тебе 👌"
	DeleteUpdate             = "Отменил 👌"
	DBProblem                = "Проблема с БД ❌"
	WrongDateFormat          = "Некорректный формат даты ❌🗓"
	WrongTimeFormat          = "Некорректный формат времени ❌🕔"
	DateBeforeNow            = "Этот день уже прошел ❌"
	TimeBeforeNow            = "Это время уже прошло ❌"
	EmptyManic               = "Пока нет записей на маникюр 🤷‍♀"
	EmptyMassage             = "Пока нет записей на массаж 🤷‍♀"
	EmptySport               = "Пока нет записей на спорт 🤷‍♀"
	EmptyMeeting             = "Пока нет встреч 🤷‍♀"
	OtherMessagesPlug        = "Ой, давай не сейчас..."
	MashaMenu                = "Чем займемся?"
	SashaMenu                = "Чего напомнить?"
	MassageQuestion          = "Что делаем с массажем?"
	ManicQuestion            = "Что делаем с маникюром?"
	SportQuestion            = "Что делаем со спортом?"
	MeetingQuestion          = "Что делаем со встречей?"
	EventNotFound            = "Не нашел такого 🤷‍♀"
	BackButtonMashaOrderMenu = "masha order menu"
	BackButtonMashaMenu      = "masha menu"
	BackButtonSashaMenu      = "sasha menu"
	USD                      = "USD"
	EURO                     = "EUR"
	CBProblem                = "Не достучался до сайта ЦБ ❌"
)

type State struct {
	State      int // 0 - question, 1 - date, 2 - time
	ChatName   string
	DeleteMode bool
	Date       string
	Time       string
}

var BackButtonStatus = make(map[int64]string)

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
			tg.NewKeyboardButton("👦 Саша"),
		),
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("👩 Маша"),
		),
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("🌦 Погода"),
			tg.NewKeyboardButton("💵 Курсы валюты"),
		),
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
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("Назад"),
		),
	)

	SashaMenuButtons = tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("🗜 Git"),
			tg.NewKeyboardButton("🐳 Docker"),
			tg.NewKeyboardButton("🕸 Kuber"),
		),
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("Назад"),
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

var (
	Git = []string{
		"-----------------Branch-----------------",
		"git branch -d branch_name",
		"git checkout -b new_branch",
		"-----------------Others-----------------",
		"rm -rf .git",
		"git rm -r --cached instance_name",
		"git reset --hard HEAD && git clean -f",
		"git commit --amend [-m 'new message'] -> :wq",
	}

	Docker = []string{
		"-----------------Image-----------------",
		"docker build --tag=image_name .",
		"docker pull image_name",
		"docker images",
		"docker rmi image_name",
		"docker rmi -f $(docker images -a -q)",
		"-----------------Container-----------------",
		"docker run -d --rm --name container_name -p 8080:8000 image_name",
		"docker run -d --name container_name -e VAR_ENV_NAME_IN_CODE=env_name image_name",
		"docker run -d --name container_name -v db:/var/lib/postgresql/data postgres:latest",
		"docker ps",
		"docker stop container_name",
		"docker rm container_name",
		"docker logs container_name",
		"docker logs -f container_name (live logs)",
		"docker inspect container_name",
		"docker exec -it container_name sh",
		"docker rm -f $(docker ps -a -q)",
		"-----------------Volume-----------------",
		"docker volume create volume_name",
		"docker volume rm volume_name",
		"docker volume ls",
		"docker inspect volume volume_name",
		"-----------------Network-----------------",
		"docker network create --driver=bridge test-net",
		"docker run -d -it --name container_name_1 test-net alpine",
		"docker run -it --name container_name_2 test-net alpine",
		"ping container_name_1",
		"-----------------Others-----------------",
		"docker system prune -a --volumes",
	}

	Kuber = []string{
		"minikube start",
		"minikube status",
		"minikube stop",
		"minikube delete",
	}
)
