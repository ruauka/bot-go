package utils

import (
	"saver-bot/internal/domain/usecases"
)

func IsButton(text string) string {
	if text == usecases.CancelButton.Keyboard[0][0].Text {
		return usecases.CancelButton.Keyboard[0][0].Text
	}

	for _, buttonName := range usecases.MainMenuButtons.Keyboard {
		if buttonName[0].Text == text {
			return buttonName[0].Text
		}
	}

	return ""
}

func IsChatState(userID int64) *usecases.State {
	manicState, ok := usecases.ManicState[userID]
	if ok {
		manicState.ChatName = usecases.Manic
		return manicState
	}

	massageState, ok := usecases.MassageState[userID]
	if ok {
		massageState.ChatName = usecases.Massage
		return massageState
	}

	sportState, ok := usecases.SportState[userID]
	if ok {
		sportState.ChatName = usecases.Sport
		return sportState
	}

	meetingState, ok := usecases.MeetingState[userID]
	if ok {
		meetingState.ChatName = usecases.Meeting
		return meetingState
	}

	return nil
}
