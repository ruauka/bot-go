package telegram

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"saver-bot/internal/domain/usecases"
	"saver-bot/internal/utils"
)

type App struct {
	usecase *usecases.Usecases
}

func NewApp(usecase *usecases.Usecases) *App {
	return &App{
		usecase: usecase,
	}
}

func (a *App) Start(updates tg.UpdatesChannel) {

	go a.usecase.Queue.QueueChanListen()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			a.usecase.Storage.CommandHandle(&update)
			continue
		}

		if button := utils.IsButton(update.Message.Text); button != "" {
			a.usecase.Storage.ButtonsHandle(&update, button)
			continue
		}

		if chatState := utils.IsChatState(update.Message.From.ID); chatState != nil {
			a.usecase.Storage.ChatStateHandle(&update, chatState)
			continue
		}

		fmt.Println(update.Message.From.UserName, update.Message.Text)
		a.usecase.Storage.MakeResponse(&update, usecases.OtherMessagesPlug)
	}
}
