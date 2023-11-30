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
	CbProblem                = "Не достучался до сайта ЦБ ❌"
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
			tg.NewKeyboardButton("🐧 Linux"),
		),
		tg.NewKeyboardButtonRow(
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
		"git push origin :branch_name - удаление remote ветки",
	}

	Docker = []string{
		"-----------------Image-----------------",
		"docker build --tag=image_name .",
		"docker pull image_name",
		"docker images",
		"docker rmi image_name",
		"docker rmi -f $(docker images -a -q)",

		"-----------------Container-----------------",
		"docker run -d --rm/--restart=always --name=container_name -p 8080:8000 image_name",
		"docker run -d --rm --name=container_name -e VAR_ENV_NAME_IN_CODE=env_name image_name",
		"docker run -d --rm --name=container_name -v db:/var/lib/postgresql/data postgres:latest",
		"docker ps (-a)",
		"docker stop container_name",
		"docker rm container_name",
		"docker logs container_name",
		"docker logs -f --tail=100 container_name",
		"docker inspect container_name",
		"docker exec -it container_name sh/bash",
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
		"apt install docker.io",
		"apt install docker-compose",
		"docker system df",
		"docker system prune -a --volumes",
		"docker image prune --all",
	}

	Kuber = []string{
		"-----------------Minikube-----------------",
		"minikube start",
		"minikube status",
		"minikube stop",
		"minikube delete",

		"-----------------General-----------------",
		"kubectl cluster-info - проверка подключения к кластеру",
		"minikube dashboard - ui миникуба",
		"kubectl port-forward --address 0.0.0.0 pod/pod-name 8000:8000 - прокинуть порты на локалку (вместо services)",
		"kubectl get replicaset replicaset-name -o=jsonpath='{.spec.template.spec.containers[0].image}' - проверка образа, указанного в ReplicaSet",
		"kubectl get pods -l app=replicaset-name -o=jsonpath='{.items[0:3].spec.containers[0].image}' - проверка образа, из которого запустились pods",

		"-----------------Pod-----------------",
		"kubectl apply -f manifest.yaml - применение манифеста pod",
		"kubectl apply -f manifest.yaml && kubectl get pods -w - запуск и просмотр в консоле",
		"kubectl get pods - посмотреть все pods",
		"kubectl describe pod pod-name - описание событий и атрибутов pod",
		"kubectl exec -it pod-name --container=container_name -- sh - провалиться в контейнер внутри pod",
		"kubectl delete pod pod-name - удалить pod",

		"-----------------Deployment-----------------",
		"kubectl get deployments - просмотр всех deployment",
		"kubectl get pods -l app=deployment-name - просмотр в консоле по конкретному тегу",
		"kubectl apply -f manifest.yaml | kubectl get pods -l app=deployment-name -w - запуск и просмотр в консоле по конкретному тегу",
		"kubectl rollout status deployment/deployment-name - проверка прохождения readinessProbe",
		"kubectl delete -n default deployment deployment-name - удалить deployment со всеми pods",
		"kubectl delete pods -l app=deployment-name - удаление pods по конкретному тегу",
		"kubectl delete pods -l app=deployment-name | kubectl get pods -l app=deployment-name -w - удаление pods по конкретному тегу и просморт",

		"-----------------Ad-hoc-----------------",
		"kubectl run frontend --image=ruauka/frontend:latest --restart=Never - альтернативный способ запуска ресурса (pod здесь)",
		"kubectl run frontend --image=ruauka/frontend:latest --restart=Never --dry-run -o yaml > frontend-pod.yaml - собрать мнифест не поднимаю pod (dry-run - режим)",
		"kubectl scale replicaset frontend --replicas=3 - увеличение реплик",
	}

	Linux = []string{
		"df -h --- разбивка по занятому месту",
		"du -hx --max-depth=15 / | grep \"[[:digit:]]\\.*G\" --- самые большие директории",
		"journalctl --vacuum-time=1d --- очистка логов в var/log/journal до 1 дня",
		"sudo lsof -i -P | grep LISTEN | grep :$PORT --- какой порт заянт",
		"sudo kill -9 <PID>",
	}
)
