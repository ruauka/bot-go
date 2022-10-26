package usecases

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Event interface {
	CommandHandle(update *tg.Update)
	MainMenuHandle(update *tg.Update, stateMap map[int64]*State)
	ManicureHandle(update *tg.Update, state *State)
	MassageHandle(update *tg.Update, state *State)
	MakeResponse(update *tg.Update, text string)
}
