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

		if button := usecases.IsMenuButton(update.Message.Text); button != "" {
			a.usecase.MenuButtonsHandle(&update, button)
			continue
		}

		if chatState := usecases.IsChatState(update.Message.From.ID); chatState != nil {
			a.usecase.ChatStateHandle(&update, chatState)
			continue
		}

		fmt.Println(update.Message.From.UserName, update.Message.Text)
		a.usecase.MakeResponse(&update, usecases.OtherMessagesPlug)
	}
}
