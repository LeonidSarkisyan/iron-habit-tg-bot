package handlers

import (
	"HabitsBot/internal/callbackdata"
	"HabitsBot/internal/filters"
	"HabitsBot/internal/keyboards"
	"HabitsBot/internal/messages"
	"HabitsBot/internal/pagination"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func NewReminderRouter(habitBot *HabitBot) *Router {
	r := NewRouter(habitBot)

	r.CallBackQuery(habitBot.ReminderList, filters.IsReminderListCallBackData)
	r.CallBackQuery(habitBot.ReminderListWithPage, filters.IsReminderPage)

	return r
}

func (h *HabitBot) ReminderList(update *tgbotapi.Update) {
	data := strings.Split(strings.Replace(update.CallbackQuery.Data, "reminder_list__", "", 1), "__")

	habitID := data[0]
	pageHabit := data[1]

	habitIDInt, err := strconv.Atoi(habitID)

	if err != nil {
		log.Error().Err(err).Msg("Невозможно получить ID привычки")
		return
	}

	habitName, err := h.HabitStorage.Name(habitIDInt, update.CallbackQuery.From.ID)

	if err != nil {
		log.Error().Err(err).Send()
		msgError := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages.ErrorGetReminders)
		_, _ = h.Bot.Send(msgError)
		return
	}

	offset := 0

	reminders, existsMore, err := h.TimestampStorage.GetByHabitID(habitIDInt, offset)

	log.Info().Any("Reminders", reminders).Msg("Список напоминаний")

	if err != nil {
		log.Error().Err(err).Send()
		msgError := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messages.ErrorGetReminders)
		_, _ = h.Bot.Send(msgError)
		return
	}

	rh := pagination.NewReminderPagination(existsMore)

	update.CallbackQuery.Message.Text = messages.ReminderListMsg(habitName)
	msg := keyboards.EditInlineKeyboard(
		keyboards.ReminderListKeyboard(habitIDInt, reminders, rh, pageHabit), update)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = h.Bot.Send(msg)

	if err != nil {
		log.Error().Err(err).Send()
	}
}

func (h *HabitBot) ReminderListWithPage(update *tgbotapi.Update) {
	callBackData := callbackdata.NewCallBackDataFromData(update.CallbackQuery.Data)
	data := callBackData.IntData()

	habitID := data[0]
	pageHabit := data[1]
	page := data[2]

	log.Info().Int("HabitID", habitID).Send()
	log.Info().Int("PageHabit", pageHabit).Send()
	log.Info().Int("Page", page).Send()
}
