package handlers

import (
	"HabitsBot/internal/filters"
	"HabitsBot/internal/messages"
	"HabitsBot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func NewStatsRouter(habitBot *HabitBot) *Router {
	r := NewRouter(habitBot)

	r.CallBackQuery(habitBot.CompleteHabit, filters.IsCallBackCompleteHabit)
	r.CallBackQuery(habitBot.CancelHabit, filters.IsCallBackCancelHabit)
	r.FSMState(GetTextRejectionState, habitBot.GetTextRejection)

	return r
}

func (h *HabitBot) CompleteHabit(update *tgbotapi.Update) {
	habitID := strings.Replace(update.CallbackQuery.Data, "complete_habit__", "", 1)

	habitIDInt, err := strconv.Atoi(habitID)

	if err != nil {
		log.Info().Err(err).Send()
		return
	}

	e := models.Execution{
		HabitID: habitIDInt,
	}

	err = h.ExecutionStorage.Create(e)

	if err != nil {
		log.Info().Err(err).Send()
		return
	}

	*h.CompleteChanMap[habitIDInt] <- "complete"

	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, messages.RandomCreateCompleteHabitMsg())
	_, _ = h.Bot.Send(msg)

	emptyMarkup := tgbotapi.InlineKeyboardMarkup{}

	emptyMarkup.InlineKeyboard = make([][]tgbotapi.InlineKeyboardButton, 0)

	msgDelete := tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, emptyMarkup)
	_, err = h.Bot.Send(msgDelete)
	if err != nil {
		log.Info().Err(err).Send()
	}
}

func (h *HabitBot) CancelHabit(update *tgbotapi.Update) {
	habitID := strings.Replace(update.CallbackQuery.Data, "cancel_habit__", "", 1)

	h.createOrUpdateSliceMetadata(update, "messages_ids", update.CallbackQuery.Message.MessageID)

	emptyMarkup := tgbotapi.InlineKeyboardMarkup{}

	msg := tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, emptyMarkup)
	_, _ = h.Bot.Send(msg)

	msgAsk := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages.GetTextRejectionMsg)
	_, _ = h.Bot.Send(msgAsk)

	h.FSM(update).SetMetadata("habit_id", habitID)
	h.FSM(update).SetState(GetTextRejectionState)
	h.AnswerCallbackQuery(update)
}

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
