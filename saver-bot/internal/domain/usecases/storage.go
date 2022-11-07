package usecases

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

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

	switch cmd {
	case StartCmd:
		e.MakeResponse(update, HelloMsg)
	case MenuCmd:
		e.makeMarkupResponse(update, MainMenu, "", MainMenuButtons)
		e.deleteChatState(update)
	}
}

func (e *storageUsecase) ButtonsHandle(update *tg.Update, button string) {
	switch {
	// cancel
	case CancelButton.Keyboard[0][0].Text == button:
		e.deleteChatState(update)
		e.changeBackButtonStatus(update, BackButtonMashaMenu)
		e.makeMarkupResponse(update, MashaMenu, "", MashaMenuButtons)

	// Main menu
	case MainMenuButtons.Keyboard[0][0].Text == button:
		// Саша
		e.makeMarkupResponse(update, MashaMenu, "", SashaMenuButtons)
		e.changeBackButtonStatus(update, BackButtonSashaMenu)
		// Маша
	case MainMenuButtons.Keyboard[1][0].Text == button:
		e.changeBackButtonStatus(update, BackButtonMashaMenu)
		e.makeMarkupResponse(update, MashaMenu, "", MashaMenuButtons)
		// Погода
	case MainMenuButtons.Keyboard[2][0].Text == button:
		e.MakeResponse(update, "в разработке")
		// Курсы валюты
	case MainMenuButtons.Keyboard[2][1].Text == button:
		e.MakeResponse(update, "в разработке")

	// Sasha menu
	// Git
	// Docker
	// Kuber

	// Назад
	case SashaMenuButtons.Keyboard[3][0].Text == button && BackButtonStatus[update.Message.From.ID] == BackButtonSashaMenu:
		e.makeMarkupResponse(update, MainMenu, "", MainMenuButtons)

	// Masha menu
	// Массаж
	case MashaMenuButtons.Keyboard[0][0].Text == button:
		e.createStateChat(update, MassageState)
		e.changeBackButtonStatus(update, BackButtonMashaOrderMenu)
		e.makeMarkupResponse(update, MassageQuestion, "", OrderButtons)
	// Маникюр
	case MashaMenuButtons.Keyboard[0][1].Text == button:
		e.createStateChat(update, ManicState)
		e.changeBackButtonStatus(update, BackButtonMashaOrderMenu)
		e.makeMarkupResponse(update, ManicQuestion, "", OrderButtons)
	// Спорт
	case MashaMenuButtons.Keyboard[1][0].Text == button:
		e.createStateChat(update, SportState)
		e.changeBackButtonStatus(update, BackButtonMashaOrderMenu)
		e.makeMarkupResponse(update, SportQuestion, "", OrderButtons)
	// Встреча
	case MashaMenuButtons.Keyboard[1][1].Text == button:
		e.createStateChat(update, MeetingState)
		e.changeBackButtonStatus(update, BackButtonMashaOrderMenu)
		e.makeMarkupResponse(update, MeetingQuestion, "", OrderButtons)
	// Назад
	case MashaMenuButtons.Keyboard[3][0].Text == button && BackButtonStatus[update.Message.From.ID] == BackButtonMashaMenu:
		e.makeMarkupResponse(update, MainMenu, "", MainMenuButtons)

	// Order menu
	// Создать
	case OrderButtons.Keyboard[0][0].Text == button:
		e.changeState(update, StateDate)
		e.changeBackButtonStatus(update, BackButtonMashaMenu)
		e.makeMarkupResponse(update, SignDate, "", CancelButton)
	// Отменить
	case OrderButtons.Keyboard[1][0].Text == button:
		e.setDeleteModeState(update)
		e.changeBackButtonStatus(update, BackButtonMashaMenu)
		e.makeMarkupResponse(update, DeleteEvent, "", CancelButton)
	// Назад
	case OrderButtons.Keyboard[2][0].Text == button && BackButtonStatus[update.Message.From.ID] == BackButtonMashaOrderMenu:
		e.changeBackButtonStatus(update, BackButtonMashaMenu)
		e.deleteChatState(update)
		e.makeMarkupResponse(update, MashaMenu, "", MashaMenuButtons)

	// all events
	case MashaMenuButtons.Keyboard[2][0].Text == button:
		e.getAllEvents(update)
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
}

func (e *storageUsecase) ChatStateHandle(update *tg.Update, state *State) {
	if state.DeleteMode {
		if !e.tryToRemoveEvent(update, state) {
			return
		}
		e.makeMarkupResponse(update, MainMenu, DeleteUpdate, MashaMenuButtons)
		e.deleteChatState(update)
		return
	}

	switch state.State {
	case StateDate:
		if !e.regExpCheck(dateRe, update, WrongDateFormat, SignDate) {
			break
		}

		if !e.dateCheck(update, SignDate) {
			break
		}

		state.Date = update.Message.Text
		e.makeMarkupResponse(update, SignTime, "", CancelButton)
		state.State = StateTime
	default:
		if !e.regExpCheck(timeRe, update, WrongTimeFormat, SignTime) {
			break
		}

		if !e.timeCheck(update, state, SignTime) {
			break
		}

		state.Time = update.Message.Text

		payload := entities.Event{}

		switch state.ChatName {
		case Manic:
			defer delete(ManicState, update.Message.From.ID)
			payload = e.makePayload(update, state.ChatName, ManicState)
		case Massage:
			defer delete(MassageState, update.Message.From.ID)
			payload = e.makePayload(update, state.ChatName, MassageState)
		case Sport:
			defer delete(SportState, update.Message.From.ID)
			payload = e.makePayload(update, state.ChatName, SportState)
		case Meeting:
			defer delete(MeetingState, update.Message.From.ID)
			payload = e.makePayload(update, state.ChatName, MeetingState)
		}

		err := e.storage.Save(context.Background(), &payload)
		if err != nil {
			log.Println(err)
			e.MakeResponse(update, DBProblem)
		}

		e.makeMarkupResponse(update, MainMenu, SaveUpdate, MashaMenuButtons)
	}
}

func (e *storageUsecase) IsChatState(userID int64) *State {
	for index, chat := range Chats {
		state, ok := chat[userID]
		if ok {
			state.ChatName = EventArr[index]
			return state
		}
	}

	return nil
}

func (e *storageUsecase) MakeResponse(update *tg.Update, text string) {
	msg := tg.NewMessage(update.Message.Chat.ID, text)
	defer func() { _, _ = e.bot.Send(msg) }()
}

func (e *storageUsecase) setDeleteModeState(update *tg.Update) {
	chatState := e.IsChatState(update.Message.From.ID)
	chatState.DeleteMode = true
	chatState.ChatName = update.Message.Text
	chatState.State = StateDate
}

func (e *storageUsecase) changeState(update *tg.Update, state int) {
	chatState := e.IsChatState(update.Message.From.ID)
	chatState.State = state
}

func (e *storageUsecase) createStateChat(update *tg.Update, state map[int64]*State) {
	state[update.Message.From.ID] = new(State)
	state[update.Message.From.ID].State = StateQuestion
}

func (e *storageUsecase) makeMarkupResponse(update *tg.Update, text, additionText string, reply tg.ReplyKeyboardMarkup) {
	msg := tg.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = reply
	if additionText != "" {
		msg.Text = additionText
	}
	defer func() { _, _ = e.bot.Send(msg) }()
}

func (e *storageUsecase) makePayload(update *tg.Update, chatName string, state map[int64]*State) entities.Event {
	return entities.Event{
		Date:     e.concat(update, state),
		Type:     chatName,
		Username: update.Message.From.UserName,
		TelegaID: update.Message.From.ID,
	}
}

func (e *storageUsecase) dateCheck(update *tg.Update, correct string) bool {
	currentDate, err := time.Parse(DatePointTimeLayout, fmt.Sprintf("%s 23:59", update.Message.Text))
	if currentDate.Before(time.Now()) || err != nil {
		log.Printf("ошибка даты: %s", err)
		e.MakeResponse(update, DateBeforeNow)
		e.MakeResponse(update, correct)
		return false
	}

	return true
}

func (e *storageUsecase) timeCheck(update *tg.Update, state *State, correct string) bool {
	currentTime, err := time.Parse(DatePointLayout, state.Date)
	if currentTime.After(time.Now()) {
		return true
	}

	currentTime, err = time.Parse(TimeDashTimeLayout, fmt.Sprintf("%s %s", time.Now().String()[:10], update.Message.Text))
	if currentTime.Before(time.Now().UTC().Add(time.Hour*3)) || err != nil {
		log.Printf("ошибка времени: %s", err)
		e.MakeResponse(update, TimeBeforeNow)
		e.MakeResponse(update, correct)
		return false
	}

	return true
}

func (e *storageUsecase) dateParse(date string) time.Time {
	currentDate, _ := time.Parse(DatePointLayout, date)
	return currentDate
}

func (e *storageUsecase) tryToRemoveEvent(update *tg.Update, state *State) bool {
	err := e.storage.Remove(context.Background(), update.Message.Text, state.ChatName, update.Message.From.UserName)
	if err != nil {
		e.MakeResponse(update, EventNotFound)
		e.MakeResponse(update, DeleteEvent)
		return false
	}

	return true
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

func (e *storageUsecase) concat(update *tg.Update, stateMap map[int64]*State) string {
	return fmt.Sprintf("%s %s",
		stateMap[update.Message.From.ID].Date,
		stateMap[update.Message.From.ID].Time,
	)
}

func (e *storageUsecase) deleteChatState(update *tg.Update) {
	delete(ManicState, update.Message.From.ID)
	delete(MassageState, update.Message.From.ID)
	delete(SportState, update.Message.From.ID)
	delete(MeetingState, update.Message.From.ID)
}

func (e *storageUsecase) changeBackButtonStatus(update *tg.Update, status string) {
	BackButtonStatus[update.Message.From.ID] = status
}

func (e *storageUsecase) getAllEvents(update *tg.Update) {
	events, err := e.storage.GetAll(context.Background())
	if err != nil {
		log.Println(err)
	}

	var manics []entities.EventAll
	var massages []entities.EventAll
	var sports []entities.EventAll
	var meetings []entities.EventAll

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
		if event.Type == Sport {
			e.Date = event.Date
			sports = append(sports, e)
		}
		if event.Type == Meeting {
			e.Date = event.Date
			meetings = append(meetings, e)
		}
	}

	if len(manics) == 0 {
		e.MakeResponse(update, EmptyManic)
	}
	if len(massages) == 0 {
		e.MakeResponse(update, EmptyMassage)
	}
	if len(sports) == 0 {
		e.MakeResponse(update, EmptySport)
	}
	if len(meetings) == 0 {
		e.MakeResponse(update, EmptyMeeting)
	}

	for _, v := range manics {
		e.MakeResponse(update, fmt.Sprintf("💅 Запись на %s", v.Date))
	}
	for _, v := range massages {
		e.MakeResponse(update, fmt.Sprintf("💆‍♀ Запись на %s", v.Date))
	}
	for _, v := range sports {
		e.MakeResponse(update, fmt.Sprintf("🏃‍♀ Запись на %s", v.Date))
	}
	for _, v := range meetings {
		e.MakeResponse(update, fmt.Sprintf("🗓 Запись на %s", v.Date))
	}
}
