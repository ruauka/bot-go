package utils

import (
	"saver-bot/internal/domain/usecases"
)

func IsButton(text string) string {
	if text == usecases.CancelButton.Keyboard[0][0].Text {
		return usecases.CancelButton.Keyboard[0][0].Text
	}

	for _, buttonName := range usecases.MainMenuButtons.Keyboard[0] {
		if buttonName.Text == text {
			return buttonName.Text
		}
	}

	return ""
}

func IsChatState(userID int64) *usecases.State {
	manicState, ok := usecases.ManicState[userID]
	if ok {
		manicState.ChatName = "Маникюр"
		return manicState
	}

	massageState, ok := usecases.MassageState[userID]
	if ok {
		massageState.ChatName = "Массаж"
		return massageState
	}

	return nil
}
