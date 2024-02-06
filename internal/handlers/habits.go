package handlers

import (
	"HabitsBot/internal/keyboards"
	"HabitsBot/internal/messages"
	"HabitsBot/internal/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func (h *HabitBot) HandleFSMHabit(update *tgbotapi.Update, done *bool) {
	if update.CallbackQuery == nil {
		state := h.FSM(update).Current()
		switch state {
		case getHabitNameState:
			h.showGettingHabitName(update, done)
		case getWarningTimeState:
			h.getHabitWarningTime(update, done)
		case getCompletedTimeState:
			h.getHabitCompletedTime(update, done)
		case getTextRejectionState:
			h.getTextRejection(update, done)
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
		case strings.HasPrefix(callBackData, "cancel_habit__"):
			h.cancelHabit(update, done)
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

	msgText := messages.InputWarningTimeMsg
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msgText)
	h.Bot.Send(msg)
	h.FSM(update).SetState(getWarningTimeState)

	h.AnswerCallbackQuery(update)

	*done = true
}

func (h *HabitBot) getHabitWarningTime(update *tgbotapi.Update, done *bool) {
	defer func() {
		*done = true
	}()

	warningTimeStr := update.Message.Text

	warningTime, err := strconv.Atoi(warningTimeStr)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.InvalidInputMsg)
		h.Bot.Send(msg)
		return
	}

	if warningTime < 0 || warningTime > 60 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.InvalidRangeInputWarningTimeMsg)
		h.Bot.Send(msg)
		return
	}

	h.FSM(update).SetMetadata("warning_time", warningTime)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.GetWarningTimeMsg)
	h.Bot.Send(msg)

	h.FSM(update).SetState(getCompletedTimeState)
}

// FSM STATE = getCompletedTimeState
func (h *HabitBot) getHabitCompletedTime(update *tgbotapi.Update, done *bool) {
	defer func() {
		*done = true
	}()

	completedTimeStr := update.Message.Text

	completedTime, err := strconv.Atoi(completedTimeStr)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.InvalidInputMsg)
		h.Bot.Send(msg)
		return
	}

	if completedTime < 0 || completedTime > 60 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.InvalidRangeInputWarningTimeMsg)
		h.Bot.Send(msg)
		return
	}

	habitName, nameExists := h.FSM(update).Metadata("habit_name")
	days, dayExists := h.FSM(update).Metadata("days")
	times, timeExists := h.FSM(update).Metadata("times")
	warningTime, warningTimeExists := h.FSM(update).Metadata("warning_time")

	if nameExists && dayExists && timeExists && warningTimeExists {
		for i, time := range times.([]string) {
			(times.([]string))[i] = strings.Replace(time, "time__", "", 1)
		}

		h.deleteMessage(update)

		habit := models.Habit{
			Title:         habitName.(string),
			UserID:        update.Message.From.ID,
			WarningTime:   warningTime.(int),
			CompletedTime: completedTime,
		}

		habit.Timestamps = models.NewTimestamps(days.([]string), times.([]string))

		h.Clear(update, "habit_name", "days", "times")

		id, err := h.HabitStorage.Create(habit)

		if err != nil {
			log.Error().Err(err).Msg(err.Error())

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.HabitCreateErrorMsg)
			h.Bot.Send(msg)
		} else {
			msgText := messages.ShowSaveHabitMsg(habitName.(string), days.([]string), times.([]string))
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			msg.ParseMode = tgbotapi.ModeHTML
			h.Bot.Send(msg)

			habit.ID = id

			h.TimeShedulerChan <- habit
		}

		log.Info().Any("habit", habit).Msg("Сохраняем")
	}
}

// CALL_BACK_DATA = cancel_habit__{habit_id}
func (h *HabitBot) cancelHabit(update *tgbotapi.Update, done *bool) {
	defer func() { *done = true }()

	habitID := strings.Replace(update.CallbackQuery.Data, "cancel_habit__", "", 1)

	h.createOrUpdateSliceMetadata(update, "messages_ids", update.CallbackQuery.Message.MessageID)

	msgAsk := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages.GetTextRejectionMsg)
	h.Bot.Send(msgAsk)

	h.FSM(update).SetMetadata("habit_id", habitID)
	h.FSM(update).SetState(getTextRejectionState)
	h.AnswerCallbackQuery(update)
}

// FSM STATE = getTextRejectionState
func (h *HabitBot) getTextRejection(update *tgbotapi.Update, done *bool) {
	defer func() { *done = true }()

	if update.Message.Text == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.InvalidRejectionMsg)
		h.Bot.Send(msg)
		return
	}

	habitID, exists := h.FSM(update).Metadata("habit_id")

	if !exists {
		log.Error().Msg("почему-то habit_id не существует")
		return
	}

	text := update.Message.Text

	habitIDInt, err := strconv.Atoi(habitID.(string))

	if err != nil {
		log.Error().Err(err).Msg("не удалось преобразовать habit_id в int")
		return
	}

	rejection := models.Rejection{
		Text:    text,
		HabitID: habitIDInt,
	}

	err = h.RejectionStorage.Create(rejection)

	if err != nil {
		msgText := tgbotapi.NewMessage(update.Message.Chat.ID, messages.RejectionCreateErrorMsg)
		h.Bot.Send(msgText)
		return
	}

	*h.ControlChanMap[habitIDInt] <- "cancel"

	habitName, err := h.HabitStorage.Name(habitIDInt, update.Message.From.ID)

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.From.ID, messages.CancelHabitErrorMsg)
		h.Bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.From.ID, messages.CancelHabitMsg(habitName, text))
	h.Bot.Send(msg)

	h.deleteMessage(update)
	h.FSM(update).DeleteMetadata("habit_id")
	h.FSM(update).SetState(startState)
}
