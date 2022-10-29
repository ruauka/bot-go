package usecases

import (
	"context"
	"fmt"
	"log"
	"regexp"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"saver-bot/internal/adapters/storage"
	"saver-bot/internal/domain/entities"
)

type storageUsecase struct {
	bot     *tg.BotAPI
	storage storage.Storage
}

func NewStorageUsecase(storage storage.Storage, bot *tg.BotAPI) StorageUsecase {
	return &storageUsecase{
		bot:     bot,
		storage: storage,
	}
}

func (e *storageUsecase) CommandHandle(update *tg.Update) {
	cmd := update.Message.Command()

	if cmd == "start" {
		e.MakeResponse(update, HelloMsg)
	}

	if cmd == "menu" {
		msg := tg.NewMessage(update.Message.Chat.ID, MainMenu)
		msg.ReplyMarkup = MainMenuButtons
		defer func() { _, _ = e.bot.Send(msg) }()
	}

	if cmd == "all" {
		e.getAllEvents(update)
	}
}

func (e *storageUsecase) ButtonsHandle(update *tg.Update, button string) {
	switch {
	case button == CancelButton.Keyboard[0][0].Text:
		delete(ManicState, update.Message.From.ID)
		delete(MassageState, update.Message.From.ID)
		e.MakeMarkupResponse(update, MainMenu, "", MainMenuButtons)
		return
	case button == MainMenuButtons.Keyboard[0][0].Text:
		MassageState[update.Message.From.ID] = new(State)
		MassageState[update.Message.From.ID].State = StateDate
	case button == MainMenuButtons.Keyboard[0][1].Text:
		ManicState[update.Message.From.ID] = new(State)
		ManicState[update.Message.From.ID].State = StateDate
	}

	e.MakeMarkupResponse(update, SignDate, "", CancelButton)
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
}

func (e *storageUsecase) ChatStateHandle(update *tg.Update, state *State) {
	switch state.State {
	case StateDate:
		if !e.regExpCheck(dateRe, update, WrongDateFormat, SignDate) {
			break
		}

		state.Date = update.Message.Text
		e.MakeMarkupResponse(update, SignTime, "", CancelButton)
		state.State = StateTime
	default:
		if !e.regExpCheck(timeRe, update, WrongTimeFormat, SignTime) {
			break
		}

		state.Time = update.Message.Text

		payload := &entities.Event{}

		switch state.ChatName {
		case Manic:
			defer delete(ManicState, update.Message.From.ID)

			payload = &entities.Event{
				Date:     concat(update, ManicState),
				Type:     state.ChatName,
				Username: update.Message.From.UserName,
				TelegaID: update.Message.From.ID,
			}
		case Massage:
			defer delete(MassageState, update.Message.From.ID)

			payload = &entities.Event{
				Date:     concat(update, MassageState),
				Type:     state.ChatName,
				Username: update.Message.From.UserName,
				TelegaID: update.Message.From.ID,
			}
		}

		err := e.storage.Save(context.Background(), payload)
		if err != nil {
			e.MakeResponse(update, DBProblem)
		}

		e.MakeMarkupResponse(update, MainMenu, SaveUpdate, MainMenuButtons)
	}
}

func (e *storageUsecase) MakeResponse(update *tg.Update, text string) {
	msg := tg.NewMessage(update.Message.Chat.ID, text)
	defer func() { _, _ = e.bot.Send(msg) }()
}

func (e *storageUsecase) MakeMarkupResponse(update *tg.Update, text, additionText string, reply tg.ReplyKeyboardMarkup) {
	msg := tg.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = reply
	if additionText != "" {
		msg.Text = additionText
	}
	defer func() { _, _ = e.bot.Send(msg) }()
}

func (e *storageUsecase) getAllEvents(update *tg.Update) {
	events, err := e.storage.GetAll(context.Background())
	if err != nil {
		log.Println(err)
	}

	var manics []entities.EventAll
	var massages []entities.EventAll

	for _, event := range events {
		e := entities.EventAll{}

		if event.Type == Manic {
			e.Date = event.Date
			manics = append(manics, e)
		}

		if event.Type == Massage {
			e.Date = event.Date
			massages = append(massages, e)
		}
	}

	if len(manics) == 0 {
		e.MakeResponse(update, EmptyManic)
	}
	if len(massages) == 0 {
		e.MakeResponse(update, EmptyMassage)
	}

	for _, v := range manics {
		e.MakeResponse(update, fmt.Sprintf("💅 Запись на %s", v.Date))
	}

	for _, v := range massages {
		e.MakeResponse(update, fmt.Sprintf("💆‍♀ Запись на %s", v.Date))
	}
}

func (e *storageUsecase) regExpCheck(pattern *regexp.Regexp, update *tg.Update, incorrect, correct string) bool {
	matched, err := regexp.MatchString(pattern.String(), update.Message.Text)
	if !matched || err != nil {
		log.Printf("ошибка регулярки: %s", err)
		e.MakeResponse(update, incorrect)
		e.MakeResponse(update, correct)
		return false
	}

	return true
}

func concat(update *tg.Update, stateMap map[int64]*State) string {
	return fmt.Sprintf("%s %s",
		stateMap[update.Message.From.ID].Date,
		stateMap[update.Message.From.ID].Time,
	)
}

//msg.ReplyMarkup = tg.NewRemoveKeyboard(false)
