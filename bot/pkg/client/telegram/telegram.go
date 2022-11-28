package telegram

import (
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"bot/internal/config"
)

func NewTelegramBot(cfg *config.Config) (*tg.BotAPI, tg.UpdatesChannel, error) {
	bot, err := tg.NewBotAPI(cfg.TelegaToken)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	return bot, updates, nil
}
