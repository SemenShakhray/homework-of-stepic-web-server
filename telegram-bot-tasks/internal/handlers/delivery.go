package handlers

import (
	"fmt"
	"taskbot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewBot(handler *Handler, conf config.Config) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(conf.TelegramAPIToken)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания бота")
	}

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(conf.WebhookURL))
	if err != nil {
		return nil, fmt.Errorf("ошибка установки вебхука")
	}

	updates := bot.ListenForWebhook("/")

	// go func() {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "tasks":
				handler.GetAllTasks(bot, update.Message)
			case "new":
				handler.NewTask(bot, update.Message)
			}
		}
	}
	// }()

	return bot, nil
}
