package handlers

import (
	"HabitsBot/internal/filters"
	"HabitsBot/internal/keyboards"
	"HabitsBot/internal/messages"
	"HabitsBot/internal/models"
	"HabitsBot/pkg/utils"
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
	r.CallBackQuery(habitBot.CreateReminder, filters.IsCallBackDataTimeContinue)
	r.CallBackQuery(habitBot.GetHabitTime, filters.IsCallBackDataStartWithTime)

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
	message, _ := h.Bot.Send(msg)

	h.createOrUpdateSliceMetadata(update, "messages_ids", message.MessageID)
	h.AnswerCallbackQuery(update)
}

func (h *HabitBot) GetHabitDay(update *tgbotapi.Update) {
	dayFromCallback := update.CallbackQuery.Data

	h.FSM(update).SetMetadata("day", dayFromCallback)

	inlineKeyboard := keyboards.DaysPickerKeyboard(&dayFromCallback)
	msg := keyboards.EditInlineKeyboard(inlineKeyboard, update)

	message, _ := h.Bot.Send(msg)

	h.createOrUpdateSliceMetadata(update, "messages_ids", message.MessageID)
	h.AnswerCallbackQuery(update)
}

func (h *HabitBot) AskHabitTime(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages.AskHabitTimeMsg)

	msg.ReplyMarkup = keyboards.TimePickerKeyboard(nil)

	message, _ := h.Bot.Send(msg)

	h.createOrUpdateSliceMetadata(update, "messages_ids", message.MessageID)
	h.AnswerCallbackQuery(update)
}

// CreateReminder callBackData == "continue"
func (h *HabitBot) CreateReminder(update *tgbotapi.Update) {
	habitID, habitIDExists := h.FSM(update).Metadata("habit_id")
	day, dayExists := h.FSM(update).Metadata("day")
	time_, timeExists := h.FSM(update).Metadata("time")

	if habitIDExists && dayExists && timeExists {

		day = strings.Replace(day.(string), "day__", "", 1)
		time_ = strings.Replace(time_.(string), "time__", "", 1)

		ts := models.Timestamp{
			HabitID: habitID.(int),
			Day:     day.(string),
			Time:    time_.(string),
		}

		err := h.TimestampStorage.Create(ts)

		if err != nil {
			log.Error().Err(err).Msg(err.Error())
			msgTextError := messages.ErrorCreateReminder
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msgTextError)
			_, _ = h.Bot.Send(msg)
		}

		habit, err := h.HabitStorage.Get(update.CallbackQuery.From.ID, habitID.(int))

		if err != nil {
			log.Error().Err(err).Send()
		}

		h.TimeShedulerChan <- TimeShedulerData{habit, ts}

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
			messages.CreateReminderMsg(habit.Title, day.(string), time_.(string)))
		msg.ParseMode = tgbotapi.ModeHTML
		_, _ = h.Bot.Send(msg)

		hour, minute, err := utils.ExtractHoursAndMinutes("time__" + time_.(string))

		if err != nil {
			log.Error().Err(err).Send()
		}

		day, err := strconv.Atoi(day.(string))

		if err != nil {
			log.Error().Err(err).Send()
		}

		dayWeekDay, err := utils.WeekDayFromInt(day)

		if err != nil {
			log.Error().Err(err).Send()
		}

		dayUntil, hourUntil, minuteUntil := utils.TimeUntil(dayWeekDay, hour, minute)

		msgText := messages.TimeWhenDoMsg(dayUntil, hourUntil, minuteUntil)
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msgText)
		_, _ = h.Bot.Send(msg)
	}

	h.deleteMessage(update)
	h.AnswerCallbackQuery(update)
}

func (h *HabitBot) GetHabitTime(update *tgbotapi.Update) {
	callBackData := update.CallbackQuery.Data

	log.Info().Any("callBackData", callBackData).Msg("callBackData")

	inlineKeyboard := keyboards.TimePickerKeyboard(&callBackData)
	msg := keyboards.EditInlineKeyboard(inlineKeyboard, update)

	message, _ := h.Bot.Send(msg)

	h.FSM(update).SetMetadata("time", callBackData)
	h.createOrUpdateSliceMetadata(update, "messages_ids", message.MessageID)
	h.AnswerCallbackQuery(update)
}
