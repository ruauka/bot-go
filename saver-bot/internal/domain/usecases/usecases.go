package usecases

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Event interface {
	CommandHandle(update *tg.Update)
	MenuButtonsHandle(update *tg.Update, button string)
	ChatStateHandle(update *tg.Update, state *State)
	MakeResponse(update *tg.Update, text string)
}
