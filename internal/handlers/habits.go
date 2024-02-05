package handlers

import (
	"HabitsBot/internal/keyboards"
	"HabitsBot/internal/messages"
	"HabitsBot/internal/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strings"
)

func (h *HabitBot) HandleFSMHabit(update *tgbotapi.Update, done *bool) {
	if update.CallbackQuery == nil {
		state := h.FSM(update).Current()
		switch state {
		case getHabitNameState:
			h.showGettingHabitName(update, done)
		}
	} else {
		callBackData := update.CallbackQuery.Data

		switch {
		case callBackData == "continue":
			h.showHabitNameAndDaysAndAskTime(update, done)
		case callBackData == "time__continue":
			h.handleSaveHabit(update, done)
		case strings.HasPrefix(callBackData, "time__"):
			h.getHabitTime(update, done)
		default:
			h.getHabitDay(update, done)
		}
	}
}

func (h *HabitBot) showGettingHabitName(update *tgbotapi.Update, done *bool) {
	habitName := update.Message.Text
	msgText := fmt.Sprintf("Ваша привычка: %s", habitName)
	h.FSM(update).SetMetadata("habit_name", habitName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	message, _ := h.Bot.Send(msg)

	h.createOrUpdateSliceMetadata(update, "messages_ids", message.MessageID)

	h.FSM(update).SetState(getHabitDaysState)

	inlineKeyboard := keyboards.DaysPickerKeyboard([]string{})

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, messages.InputHabitDaysMsg)

	msg.ReplyMarkup = inlineKeyboard

	message, err := h.Bot.Send(msg)

	h.createOrUpdateSliceMetadata(update, "messages_ids", message.MessageID)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}

	*done = true
}

func (h *HabitBot) getHabitDay(update *tgbotapi.Update, done *bool) {
	callbackData := update.CallbackQuery.Data

	days, existsDays := h.FSM(update).Metadata("days")

	if !existsDays {
		h.FSM(update).SetMetadata("days", []string{callbackData})
		days = []string{callbackData}
	} else {
		dayExists := false

		for i, day := range days.([]string) {
			if day == callbackData {
				dayExists = true
				days = append(days.([]string)[:i], days.([]string)[i+1:]...)
				break
			}
		}

		if !dayExists {
			days = append(days.([]string), callbackData)
		}

		h.FSM(update).SetMetadata("days", days)
	}

	_, exists := h.FSM(update).Metadata("habit_name")

	if !exists {
		log.Info().Msg("habitName нет в FSM Storage")
	}

	inlineKeyboard := keyboards.DaysPickerKeyboard(days.([]string))
	msg := keyboards.EditInlineKeyboard(inlineKeyboard, update)

	message, _ := h.Bot.Send(msg)

	h.createOrUpdateSliceMetadata(update, "messages_ids", message.MessageID)

	h.AnswerCallbackQuery(update)
	*done = true
}

// callBackData: "continue"
func (h *HabitBot) showHabitNameAndDaysAndAskTime(update *tgbotapi.Update, done *bool) {
	habitName, nameExists := h.FSM(update).Metadata("habit_name")
	days, dayExists := h.FSM(update).Metadata("days")

	if nameExists && dayExists {
		msgText := messages.ShowHabitNameAndDaysMsg(habitName.(string), days.([]string))
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msgText)
		message, _ := h.Bot.Send(msg)
		h.createOrUpdateSliceMetadata(update, "messages_ids", message.MessageID)
	} else {
		log.Error().Msg("Не хватает метаданных")
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages.InputHabitTimeMsg)

	msg.ReplyMarkup = keyboards.TimePickerKeyboard([]string{})

	message, _ := h.Bot.Send(msg)

	h.createOrUpdateSliceMetadata(update, "messages_ids", message.MessageID)

	h.AnswerCallbackQuery(update)
	*done = true
}

func (h *HabitBot) getHabitTime(update *tgbotapi.Update, done *bool) {
	callBackData := update.CallbackQuery.Data

	times, exists := h.FSM(update).Metadata("times")

	if !exists {
		h.FSM(update).SetMetadata("times", []string{callBackData})
		times = []string{callBackData}
	} else {
		timeExists := false

		for i, time := range times.([]string) {
			if time == callBackData {
				timeExists = true
				times = append(times.([]string)[:i], times.([]string)[i+1:]...)
				break
			}
		}

		if !timeExists {
			times = append(times.([]string), callBackData)
		}

		h.FSM(update).SetMetadata("times", times)
	}

	inlineKeyboard := keyboards.TimePickerKeyboard(times.([]string))
	msg := keyboards.EditInlineKeyboard(inlineKeyboard, update)

	message, _ := h.Bot.Send(msg)

	h.createOrUpdateSliceMetadata(update, "messages_ids", message.MessageID)

	h.AnswerCallbackQuery(update)

	*done = true
}

func (h *HabitBot) handleSaveHabit(update *tgbotapi.Update, done *bool) {
	habitName, nameExists := h.FSM(update).Metadata("habit_name")
	days, dayExists := h.FSM(update).Metadata("days")
	times, timeExists := h.FSM(update).Metadata("times")

	if nameExists && dayExists && timeExists {
		for i, time := range times.([]string) {
			(times.([]string))[i] = strings.Replace(time, "time__", "", 1)
		}

		h.deleteMessage(update)

		habit := models.Habit{
			Title:  habitName.(string),
			UserID: update.CallbackQuery.From.ID,
		}

		habit.Timestamps = models.NewTimestamps(days.([]string), times.([]string))

		h.Clear(update, "habit_name", "days", "times")

		err := h.HabitStorage.Create(habit)

		if err != nil {
			log.Error().Err(err).Msg(err.Error())

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages.HabitCreateErrorMsg)
			h.Bot.Send(msg)
		} else {
			msgText := messages.ShowSaveHabitMsg(habitName.(string), days.([]string), times.([]string))
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msgText)
			msg.ParseMode = tgbotapi.ModeHTML
			h.Bot.Send(msg)

			h.TimeShedulerChan <- habit
		}

		log.Info().Any("habit", habit).Msg("Сохраняем")
		// TODO: implement saving habit to DB
	}

	h.AnswerCallbackQuery(update)

	*done = true
}
