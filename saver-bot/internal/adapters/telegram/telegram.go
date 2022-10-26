package telegram

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"saver-bot/internal/adapters/queue"
	"saver-bot/internal/domain/usecases"
)

type App struct {
	usecase usecases.Event
	queue   queue.Queue
}

func NewApp(usecase usecases.Event, queue queue.Queue) *App {
	return &App{
		usecase: usecase,
		queue:   queue,
	}
}

func (a *App) Start(updates tg.UpdatesChannel) {

	go a.queue.QueueChanListen()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			a.usecase.CommandHandle(&update)
			continue
		}

		if update.Message.Text == usecases.MainMenu.Keyboard[0][1].Text {
			a.usecase.MainMenuHandle(&update, usecases.ManicState)
			continue
		}

		if update.Message.Text == usecases.MainMenu.Keyboard[0][0].Text {
			a.usecase.MainMenuHandle(&update, usecases.MassageState)
			continue
		}

		manicState, ok := usecases.ManicState[update.Message.From.ID]
		if ok {
			a.usecase.ManicureHandle(&update, manicState)
			continue
		}

		massageState, ok := usecases.MassageState[update.Message.From.ID]
		if ok {
			a.usecase.MassageHandle(&update, massageState)
			continue
		}

		fmt.Println(update.Message.From.UserName, update.Message.Text)
		a.usecase.MakeResponse(&update, "Ой, давай не сейчас...")
	}
}
