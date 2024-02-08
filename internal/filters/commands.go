package filters

import (
	"HabitsBot/internal/commands"
	"HabitsBot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type Filter func(update *tgbotapi.Update) bool

func F(filters ...Filter) func(update *tgbotapi.Update) bool {
	return func(update *tgbotapi.Update) bool {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Err(r.(error)).Send()
			}
		}()

		var bools []bool

		for _, filter := range filters {
			bools = append(bools, filter(update))
		}

		return utils.All(bools)
	}
}

func IsCommandAddNewHabit(update *tgbotapi.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == commands.AddNewHabitCommand
}

func IsCommandMyHabits(update *tgbotapi.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == commands.MyHabitsCommand
}

func IsCommandStart(update *tgbotapi.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == commands.StartCommand
}

func IsCommandCancel(update *tgbotapi.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == commands.CancelCommand
}
