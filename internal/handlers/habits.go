package handlers

import (
	"HabitsBot/internal/filters"
	"HabitsBot/internal/keyboards"
	"HabitsBot/internal/messages"
	"HabitsBot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func NewHabitsFSMRouter(habitBot *HabitBot) *Router {
	r := NewRouter(habitBot)

	r.FSMState(GetHabitNameState, habitBot.ShowGettingHabitName, filters.IsCallbackQueryEmpty)
	r.FSMState(
		CreateHabitState,
		habitBot.CreateHabit,
		filters.IsCreateHabitText,
	)
	r.FSMState(
		CreateHabitState,
		habitBot.CancelCreateHabit,
		filters.IsCancelCreateHabitText,
	)
	r.FSMState(
		GetCompletedTimeState,
		habitBot.GetHabitCompletedTime,
	)

	return r
}

func NewHabitsCallBackRouter(habitBot *HabitBot) *Router {
	r := NewRouter(habitBot)

	r.CallBackQuery(habitBot.AskDayReminder, filters.IsCallBackDataAddReminder)
	r.CallBackQuery(habitBot.AskHabitTime, filters.IsCallBackDataContinue)
	r.CallBackQuery(habitBot.GetHabitDay, filters.IsCallBackDataStartWithDay)

	return r
}

// ShowGettingHabitName FSM State = GetHabitNameState
func (h *HabitBot) ShowGettingHabitName(update *tgbotapi.Update) {
	habitName := update.Message.Text

	h.FSM(update).SetMetadata("habit_name", habitName)

	msgText := messages.BeforeCreateHabitMsg(habitName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = keyboards.BeforeCreateHabitReplyKeyboard()
	_, _ = h.Bot.Send(msg)

	h.FSM(update).SetState(CreateHabitState)
}

func (h *HabitBot) CreateHabit(update *tgbotapi.Update) {
	_, exists := h.FSM(update).Metadata("habit_name")

	if !exists {
		log.Info().Msg("habitName нет в FSM Storage")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.AskTimeCompleted)
	msg.ParseMode = tgbotapi.ModeHTML
	_, _ = h.Bot.Send(msg)

	h.FSM(update).SetState(GetCompletedTimeState)
}

func (h *HabitBot) CancelCreateHabit(update *tgbotapi.Update) {
	msgText := messages.CancelCreateHabitMsg
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	_, _ = h.Bot.Send(msg)

	h.Clear(update, "habit_name")
}

// GetHabitCompletedTime FSM STATE = GetCompletedTimeState
func (h *HabitBot) GetHabitCompletedTime(update *tgbotapi.Update) {
	completedTimeStr := update.Message.Text

	completedTime, err := strconv.Atoi(completedTimeStr)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.InvalidInputMsg)
		_, _ = h.Bot.Send(msg)
		return
	}

	if completedTime < 15 || completedTime > 300 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.InvalidRangeInputWarningTimeMsg)
		_, _ = h.Bot.Send(msg)
		return
	}

	habitName, _ := h.FSM(update).Metadata("habit_name")

	habit := models.Habit{
		Title:         habitName.(string),
		UserID:        update.Message.From.ID,
		CompletedTime: completedTime,
	}

	habitID, err := h.HabitStorage.Create(habit)

	if err != nil {
		msgText := messages.ErrorCreateHabitMsg
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		_, _ = h.Bot.Send(msg)
	}

	msgText := messages.HabitCreatedMsg(habitName.(string))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = keyboards.CreateDayTimeInlineKeyboard(habitID)
	_, _ = h.Bot.Send(msg)

	h.Clear(update)
}

func (h *HabitBot) AskDayReminder(update *tgbotapi.Update) {
	habitID := strings.Replace(update.CallbackQuery.Data, "add_reminder__", "", 1)

	habitIDInt, err := strconv.Atoi(habitID)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return
	}

	h.FSM(update).SetMetadata("habit_id", habitIDInt)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages.AskHabitDaysMsg)
	msg.ReplyMarkup = keyboards.DaysPickerKeyboard(nil)
	_, _ = h.Bot.Send(msg)

	h.AnswerCallbackQuery(update)
}

func (h *HabitBot) GetHabitDay(update *tgbotapi.Update) {
	dayFromCallback := update.CallbackQuery.Data

	h.FSM(update).SetMetadata("day", dayFromCallback)

	inlineKeyboard := keyboards.DaysPickerKeyboard(&dayFromCallback)
	msg := keyboards.EditInlineKeyboard(inlineKeyboard, update)

	_, _ = h.Bot.Send(msg)

	h.AnswerCallbackQuery(update)
}

func (h *HabitBot) AskHabitTime(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages.AskHabitTimeMsg)

	msg.ReplyMarkup = keyboards.TimePickerKeyboard(nil)

	_, _ = h.Bot.Send(msg)

	h.AnswerCallbackQuery(update)
}

// CreateReminder callBackData == "continue"
func (h *HabitBot) CreateReminder(update *tgbotapi.Update) {
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

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages.AskHabitTimeMsg)

	msg.ReplyMarkup = keyboards.TimePickerKeyboard([]string{})

	message, _ := h.Bot.Send(msg)

	h.createOrUpdateSliceMetadata(update, "messages_ids", message.MessageID)

	h.AnswerCallbackQuery(update)
}

func (h *HabitBot) GetHabitTime(update *tgbotapi.Update) {
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
}

func (h *HabitBot) HandleSaveHabit(update *tgbotapi.Update) {

	msgText := messages.InputWarningTimeMsg
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msgText)
	h.Bot.Send(msg)
	h.FSM(update).SetState(GetWarningTimeState)

	h.AnswerCallbackQuery(update)
}

func (h *HabitBot) GetHabitWarningTime(update *tgbotapi.Update) {
	warningTimeStr := update.Message.Text

	warningTime, err := strconv.Atoi(warningTimeStr)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.InvalidInputMsg)
		h.Bot.Send(msg)
		return
	}

	if warningTime < 5 || warningTime > 60 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.InvalidRangeInputWarningTimeMsg)
		h.Bot.Send(msg)
		return
	}

	h.FSM(update).SetMetadata("warning_time", warningTime)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.GetWarningTimeMsg)
	h.Bot.Send(msg)

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, messages.InputCompleteTimeMsg)
	h.Bot.Send(msg)

	h.FSM(update).SetState(GetCompletedTimeState)
}

// CancelHabit CALL_BACK_DATA = cancel_habit__{habit_id}
func (h *HabitBot) CancelHabit(update *tgbotapi.Update) {
	habitID := strings.Replace(update.CallbackQuery.Data, "cancel_habit__", "", 1)

	h.createOrUpdateSliceMetadata(update, "messages_ids", update.CallbackQuery.Message.MessageID)

	msgAsk := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages.GetTextRejectionMsg)
	h.Bot.Send(msgAsk)

	h.FSM(update).SetMetadata("habit_id", habitID)
	h.FSM(update).SetState(GetTextRejectionState)
	h.AnswerCallbackQuery(update)
}

// FSM STATE = GetTextRejectionState
func (h *HabitBot) GetTextRejection(update *tgbotapi.Update) {
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
	h.FSM(update).SetState(StartState)
}
