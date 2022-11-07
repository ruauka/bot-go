package telegram

import (
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewTelegramBot() (*tg.BotAPI, tg.UpdatesChannel, error) {
	//master
	bot, err := tg.NewBotAPI("5719088062:AAF4tRC7pzzxjcViHmf4BPXVwg5qErdM4zA")
	// dev
	//bot, err := tg.NewBotAPI("5756888484:AAHQDc0aaLrVfO1z8JZRF_dy8DmnK-_bBMc")
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	return bot, updates, nil
}
