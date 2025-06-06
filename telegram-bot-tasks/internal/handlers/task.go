package handlers

import (
	"encoding/json"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (h *Handler) GetAllTasks(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	userID := bot.Self.ID
	userName := bot.Self.UserName

	res, err := h.serv.GetAllTasks(userID, userName)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(
			msg.Chat.ID,
			err.Error(),
		))

		return
	}

	bot.Send(tgbotapi.NewMessage(
		msg.Chat.ID,
		Response(res)),
	)
}

func (h *Handler) NewTask(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	description := strings.TrimSpace(msg.CommandArguments())

	if description == "" {
		bot.Send(tgbotapi.NewMessage(
			msg.Chat.ID,
			"specify the task",
		))

		return
	}

	assignName := bot.Self.UserName
	ownerID := bot.Self.ID

	res, err := h.serv.NewTask(description, assignName, ownerID)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(
			msg.Chat.ID,
			err.Error(),
		))

		return
	}

	bot.Send(tgbotapi.NewMessage(
		msg.Chat.ID,
		Response(res),
	))
}

func Response(res map[int]string) string {
	resp, _ := json.Marshal(res)

	return string(resp)
}
