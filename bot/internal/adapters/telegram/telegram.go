package telegram

import (
	"fmt"
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"bot/internal/domain/usecases"
)

type App struct {
	usecases *usecases.Usecases
	bot      *tg.BotAPI
}

func NewApp(usecases *usecases.Usecases, bot *tg.BotAPI) *App {
	return &App{
		usecases: usecases,
		bot:      bot,
	}
}

func (a *App) Start(updates tg.UpdatesChannel) {

	go a.usecases.Queue.QueueChanListen()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			a.CommandHandle(&update)
			continue
		}

		if button := isButton(update.Message.Text); button != "" {
			a.ButtonsHandle(&update, button)
			continue
		}

		if chatState := isChatState(update.Message.From.ID); chatState != nil {
			if chatState.State != StateQuestion {
				a.ChatStateHandle(&update, chatState)
				continue
			}
		}

		fmt.Println(update.Message.From.UserName, update.Message.Text)
		a.MakeResponse(&update, OtherMessagesPlug)
	}
}

func (a *App) CommandHandle(update *tg.Update) {
	cmd := update.Message.Command()

	switch cmd {
	case StartCmd:
		a.MakeResponse(update, HelloMsg)
	case MenuCmd:
		a.MakeMarkupResponse(update, MainMenu, "", MainMenuButtons)
		deleteChatState(update)
	}
}

func (a *App) ButtonsHandle(update *tg.Update, button string) {
	switch {
	// cancel
	case CancelButton.Keyboard[0][0].Text == button:
		deleteChatState(update)
		changeBackButtonStatus(update, BackButtonMashaMenu)
		a.MakeMarkupResponse(update, MashaMenu, "", MashaMenuButtons)

	// Main menu
	case MainMenuButtons.Keyboard[0][0].Text == button:
		// Саша
		a.MakeMarkupResponse(update, SashaMenu, "", SashaMenuButtons)
		changeBackButtonStatus(update, BackButtonSashaMenu)
		// Маша
	case MainMenuButtons.Keyboard[1][0].Text == button:
		changeBackButtonStatus(update, BackButtonMashaMenu)
		a.MakeMarkupResponse(update, MashaMenu, "", MashaMenuButtons)
		// Погода
	case MainMenuButtons.Keyboard[2][0].Text == button:
		//a.MakeResponse(update, "в разработке")
		a.Forecast(update)
		// Курсы валюты
	case MainMenuButtons.Keyboard[2][1].Text == button:
		a.Currency(update)

	// Sasha menu
	// Git
	case SashaMenuButtons.Keyboard[0][0].Text == button:
		a.MakeSliceResponse(update, Git)
		// Docker
	case SashaMenuButtons.Keyboard[1][0].Text == button:
		a.MakeSliceResponse(update, Docker)
		// Kuber
	case SashaMenuButtons.Keyboard[2][0].Text == button:
		a.MakeSliceResponse(update, Kuber)
		// Назад
	case SashaMenuButtons.Keyboard[3][0].Text == button && BackButtonStatus[update.Message.From.ID] == BackButtonSashaMenu:
		a.MakeMarkupResponse(update, MainMenu, "", MainMenuButtons)

	// Masha menu
	// Массаж
	case MashaMenuButtons.Keyboard[0][0].Text == button:
		createStateChat(update, MassageState)
		changeBackButtonStatus(update, BackButtonMashaOrderMenu)
		a.MakeMarkupResponse(update, MassageQuestion, "", OrderButtons)
		// Маникюр
	case MashaMenuButtons.Keyboard[0][1].Text == button:
		createStateChat(update, ManicState)
		changeBackButtonStatus(update, BackButtonMashaOrderMenu)
		a.MakeMarkupResponse(update, ManicQuestion, "", OrderButtons)
		// Спорт
	case MashaMenuButtons.Keyboard[1][0].Text == button:
		createStateChat(update, SportState)
		changeBackButtonStatus(update, BackButtonMashaOrderMenu)
		a.MakeMarkupResponse(update, SportQuestion, "", OrderButtons)
		// Встреча
	case MashaMenuButtons.Keyboard[1][1].Text == button:
		createStateChat(update, MeetingState)
		changeBackButtonStatus(update, BackButtonMashaOrderMenu)
		a.MakeMarkupResponse(update, MeetingQuestion, "", OrderButtons)
		// Назад
	case MashaMenuButtons.Keyboard[3][0].Text == button && BackButtonStatus[update.Message.From.ID] == BackButtonMashaMenu:
		a.MakeMarkupResponse(update, MainMenu, "", MainMenuButtons)

	// Order menu
	// Создать
	case OrderButtons.Keyboard[0][0].Text == button:
		changeState(update, StateDate)
		changeBackButtonStatus(update, BackButtonMashaMenu)
		a.MakeMarkupResponse(update, SignDate, "", CancelButton)
		// Отменить
	case OrderButtons.Keyboard[1][0].Text == button:
		setDeleteModeState(update)
		changeBackButtonStatus(update, BackButtonMashaMenu)
		a.MakeMarkupResponse(update, DeleteEvent, "", CancelButton)
		// Назад
	case OrderButtons.Keyboard[2][0].Text == button && BackButtonStatus[update.Message.From.ID] == BackButtonMashaOrderMenu:
		changeBackButtonStatus(update, BackButtonMashaMenu)
		deleteChatState(update)
		a.MakeMarkupResponse(update, MashaMenu, "", MashaMenuButtons)

	// all events
	case MashaMenuButtons.Keyboard[2][0].Text == button:
		a.GetAllEvents(update)
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
}

func (a *App) MakeResponse(update *tg.Update, text string) {
	msg := tg.NewMessage(update.Message.Chat.ID, text)
	defer func() { _, _ = a.bot.Send(msg) }()
}
