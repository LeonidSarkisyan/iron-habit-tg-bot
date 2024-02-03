package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MustCommands(bot *tgbotapi.BotAPI) {
	cmdCfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "Начать",
		},
		tgbotapi.BotCommand{
			Command:     "add_new_habit",
			Description: "Добавить новую привычку",
		},
	)

	_, _ = bot.Send(cmdCfg)
}
