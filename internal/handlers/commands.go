package handlers

import (
	"HabitsBot/internal/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MustCommands(bot *tgbotapi.BotAPI) {
	cmdCfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     commands.StartCommand,
			Description: "Начать",
		},
		tgbotapi.BotCommand{
			Command:     commands.AddNewHabitCommand,
			Description: "Добавить новую привычку",
		},
		tgbotapi.BotCommand{
			Command:     commands.MyHabitsCommand,
			Description: "Мои привычки",
		},
		tgbotapi.BotCommand{
			Command:     cancelCommand,
			Description: "Отменить действие",
		},
	)

	_, _ = bot.Send(cmdCfg)
}
