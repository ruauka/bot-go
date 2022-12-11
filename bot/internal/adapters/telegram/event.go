package telegram

import (
	"fmt"
	"log"
	"regexp"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"bot/internal/domain/entities"
)

func (a *App) ChatStateHandle(update *tg.Update, state *State) {
	if state.DeleteMode {
		if !a.TryToRemoveEvent(update, state) {
			return
		}
		a.MakeMarkupResponse(update, MainMenu, DeleteUpdate, MashaMenuButtons)
		deleteChatState(update)
		return
	}

	switch state.State {
	case StateDate:
		if !a.RegExpCheck(dateRe, update, WrongDateFormat, SignDate) {
			break
		}

		if !a.DateCheck(update, SignDate) {
			break
		}

		// если создание Встречи
		if state.ChatName == Meeting {
			state.Date = update.Message.Text
			a.MakeMarkupResponse(update, MeetingSignTime, "", CancelButton)
			state.State = StateTime
			break
		}

		state.Date = update.Message.Text
		a.MakeMarkupResponse(update, SignTime, "", CancelButton)
		state.State = StateTime
	case StateTime:
		if !a.RegExpCheck(timeRe, update, WrongTimeFormat, SignTime) {
			break
		}

		if !a.TimeCheck(update, state, SignTime) {
			break
		}

		// если создание Встречи
		if state.ChatName == Meeting {
			state.Time = update.Message.Text
			a.MakeMarkupResponse(update, MeetingSignWithWhom, "", CancelButton)
			state.State = StateMeeting
			break
		}

		state.Time = update.Message.Text

		payload := entities.Event{}

		switch state.ChatName {
		case Manic:
			defer delete(ManicState, update.Message.From.ID)
			payload = makePayload(update, state.ChatName, "", ManicState)
		case Massage:
			defer delete(MassageState, update.Message.From.ID)
			payload = makePayload(update, state.ChatName, "", MassageState)
		case Sport:
			defer delete(SportState, update.Message.From.ID)
			payload = makePayload(update, state.ChatName, "", SportState)
		}

		err := a.usecases.Event.Save(&payload)
		if err != nil {
			a.MakeResponse(update, DBProblem)
		}

		a.MakeMarkupResponse(update, MainMenu, SaveUpdate, MashaMenuButtons)
	case StateMeeting:
		payload := entities.Event{}

		defer delete(MeetingState, update.Message.From.ID)
		payload = makePayload(update, state.ChatName, update.Message.Text, MeetingState)

		err := a.usecases.Event.Save(&payload)
		if err != nil {
			a.MakeResponse(update, DBProblem)
		}

		a.MakeMarkupResponse(update, MainMenu, SaveUpdate, MashaMenuButtons)
	}
}

func (a *App) GetAllEvents(update *tg.Update) {
	events, err := a.usecases.Event.GetAll()
	if err != nil {
		a.MakeResponse(update, DBProblem)
	}

	var manics []entities.EventAll
	var massages []entities.EventAll
	var sports []entities.EventAll
	var meetings []entities.EventMeeting

	for _, event := range events {
		e := entities.EventAll{}
		m := entities.EventMeeting{}

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
			m.Date = event.Date
			m.Whom = event.Whom
			meetings = append(meetings, m)
		}
	}

	if len(manics) == 0 {
		a.MakeResponse(update, EmptyManic)
	}
	if len(massages) == 0 {
		a.MakeResponse(update, EmptyMassage)
	}
	if len(sports) == 0 {
		a.MakeResponse(update, EmptySport)
	}
	if len(meetings) == 0 {
		a.MakeResponse(update, EmptyMeeting)
	}

	for _, v := range manics {
		a.MakeResponse(update, fmt.Sprintf("💅 Запись на %s", v.Date))
	}
	for _, v := range massages {
		a.MakeResponse(update, fmt.Sprintf("💆‍♀ Запись на %s", v.Date))
	}
	for _, v := range sports {
		a.MakeResponse(update, fmt.Sprintf("🏃‍♀ Запись на %s", v.Date))
	}
	for _, v := range meetings {
		a.MakeResponse(update, fmt.Sprintf("🗓 Встреча %s на %s", v.Whom, v.Date))
	}
}

func (a *App) MakeSliceResponse(update *tg.Update, cmdSlice []string) {
	for _, cmd := range cmdSlice {
		msg := tg.NewMessage(update.Message.Chat.ID, cmd)
		a.bot.Send(msg)
	}
}

func (a *App) MakeMarkupResponse(update *tg.Update, text, additionText string, reply tg.ReplyKeyboardMarkup) {
	msg := tg.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = reply
	if additionText != "" {
		msg.Text = additionText
	}
	defer func() { _, _ = a.bot.Send(msg) }()
}

func (a *App) DateCheck(update *tg.Update, correct string) bool {
	currentDate, err := time.Parse(DatePointTimeLayout, fmt.Sprintf("%s 23:59", update.Message.Text))
	if currentDate.Before(time.Now()) || err != nil {
		log.Printf("ошибка даты: %s", err)
		a.MakeResponse(update, DateBeforeNow)
		a.MakeResponse(update, correct)
		return false
	}

	return true
}

func (a *App) TimeCheck(update *tg.Update, state *State, correct string) bool {
	currentTime, err := time.Parse(DatePointLayout, state.Date)
	if currentTime.After(time.Now()) {
		return true
	}

	currentTime, err = time.Parse(TimeDashTimeLayout, fmt.Sprintf("%s %s", time.Now().String()[:10], update.Message.Text))
	if currentTime.Before(time.Now().UTC().Add(time.Hour*3)) || err != nil {
		log.Printf("ошибка времени: %s", err)
		a.MakeResponse(update, TimeBeforeNow)
		a.MakeResponse(update, correct)
		return false
	}

	return true
}

func (a *App) TryToRemoveEvent(update *tg.Update, state *State) bool {
	err := a.usecases.Event.Remove(update.Message.Text, state.ChatName, update.Message.From.UserName)
	if err != nil {
		a.MakeResponse(update, EventNotFound)
		a.MakeResponse(update, DeleteEvent)
		return false
	}

	return true
}

func (a *App) RegExpCheck(pattern *regexp.Regexp, update *tg.Update, incorrect, correct string) bool {
	matched, err := regexp.MatchString(pattern.String(), update.Message.Text)
	if !matched || err != nil {
		log.Printf("ошибка регулярки: %s", err)
		a.MakeResponse(update, incorrect)
		a.MakeResponse(update, correct)
		return false
	}

	return true
}

func changeBackButtonStatus(update *tg.Update, status string) {
	BackButtonStatus[update.Message.From.ID] = status
}

func isButton(text string) string {
	if text == CancelButton.Keyboard[0][0].Text {
		return CancelButton.Keyboard[0][0].Text
	}

	for _, buttonSlice := range MainMenuButtons.Keyboard {
		for _, buttonName := range buttonSlice {
			if buttonName.Text == text {
				return buttonName.Text
			}
		}
	}

	for _, buttonSlice := range SashaMenuButtons.Keyboard {
		for _, buttonName := range buttonSlice {
			if buttonName.Text == text {
				return buttonName.Text
			}
		}
	}

	for _, buttonSlice := range MashaMenuButtons.Keyboard {
		for _, buttonName := range buttonSlice {
			if buttonName.Text == text {
				return buttonName.Text
			}
		}
	}

	for _, buttonName := range OrderButtons.Keyboard {
		if buttonName[0].Text == text {
			return buttonName[0].Text
		}
	}

	return ""
}

func isChatState(userID int64) *State {
	for index, chat := range Chats {
		state, ok := chat[userID]
		if ok {
			state.ChatName = EventArr[index]
			return state
		}
	}

	return nil
}

func meetingChatCheck(userID int64) bool {
	_, ok := MeetingState[userID]
	if ok {
		return true
	}

	return false
}

func setDeleteModeState(update *tg.Update) {
	chatState := isChatState(update.Message.From.ID)
	chatState.DeleteMode = true
	chatState.ChatName = update.Message.Text
	chatState.State = StateDate
}

func changeState(update *tg.Update, state int) {
	chatState := isChatState(update.Message.From.ID)
	chatState.State = state
}

func createStateChat(update *tg.Update, state map[int64]*State, chatName string) {
	state[update.Message.From.ID] = new(State)
	state[update.Message.From.ID].State = StateQuestion
	state[update.Message.From.ID].ChatName = chatName
}

func makePayload(update *tg.Update, chatName, whom string, state map[int64]*State) entities.Event {
	return entities.Event{
		Date:     concat(update, state),
		Whom:     whom,
		Type:     chatName,
		Username: update.Message.From.UserName,
		TelegaID: update.Message.From.ID,
	}
}

func concat(update *tg.Update, stateMap map[int64]*State) string {
	return fmt.Sprintf("%s %s",
		stateMap[update.Message.From.ID].Date,
		stateMap[update.Message.From.ID].Time,
	)
}

func deleteChatState(update *tg.Update) {
	delete(ManicState, update.Message.From.ID)
	delete(MassageState, update.Message.From.ID)
	delete(SportState, update.Message.From.ID)
	delete(MeetingState, update.Message.From.ID)
}
