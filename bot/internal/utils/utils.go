package utils

import (
	"bot/internal/domain/usecases"
)

func IsButton(text string) string {
	if text == usecases.CancelButton.Keyboard[0][0].Text {
		return usecases.CancelButton.Keyboard[0][0].Text
	}

	for _, buttonSlice := range usecases.MainMenuButtons.Keyboard {
		for _, buttonName := range buttonSlice {
			if buttonName.Text == text {
				return buttonName.Text
			}
		}
	}

	for _, buttonName := range usecases.SashaMenuButtons.Keyboard {
		if buttonName[0].Text == text {
			return buttonName[0].Text
		}
	}

	for _, buttonSlice := range usecases.MashaMenuButtons.Keyboard {
		for _, buttonName := range buttonSlice {
			if buttonName.Text == text {
				return buttonName.Text
			}
		}
	}

	for _, buttonName := range usecases.OrderButtons.Keyboard {
		if buttonName[0].Text == text {
			return buttonName[0].Text
		}
	}

	return ""
}
