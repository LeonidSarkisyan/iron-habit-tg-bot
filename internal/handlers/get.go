package handlers

import (
	"HabitsBot/internal/filters"
	"HabitsBot/internal/keyboards"
	"HabitsBot/internal/messages"
	"HabitsBot/internal/models"
	"HabitsBot/internal/pagination"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"sort"
	"strconv"
	"strings"
)

type ByID []models.Habit

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Habits COMMAND = /my_habits
func (h *HabitBot) Habits(u *tgbotapi.Update) {
	userID := u.Message.From.ID

	habits, existsMore, err := h.HabitStorage.GetAll(userID, 0)

	log.Info().Any("habits", habits).Send()

	sort.Sort(ByID(habits))

	if err != nil {
		msgError := tgbotapi.NewMessage(userID, messages.ErrorGetHabits)
		h.Bot.Send(msgError)
		return
	}

	hp := pagination.NewHabitPagination(existsMore)

	msg := tgbotapi.NewMessage(userID, messages.HabitListMsg(habits))
	msg.ReplyMarkup = keyboards.HabitsListKeyboard(habits, hp)
	msg.ParseMode = "HTML"
	h.Bot.Send(msg)
}

func NewHabitListRouter(habitBot *HabitBot) *Router {
	r := NewRouter(habitBot)

	r.CallBackQuery(habitBot.HabitsWithPage, filters.IsPage)

	return r
}

func (h *HabitBot) HabitsWithPage(update *tgbotapi.Update) {
	page := strings.Replace(update.CallbackQuery.Data, "habits__page__", "", 1)

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	userID := update.CallbackQuery.From.ID

	offset := (pageInt - 1) * 5

	habits, existsMore, err := h.HabitStorage.GetAll(userID, offset)

	sort.Sort(ByID(habits))

	if err != nil {
		msgError := tgbotapi.NewMessage(userID, messages.ErrorGetHabits)
		h.Bot.Send(msgError)
		return
	}

	hp := pagination.NewHabitPagination(existsMore)
	hp.Page = pageInt

	msg := keyboards.EditInlineKeyboard(keyboards.HabitsListKeyboard(habits, hp), update)
	msg.ParseMode = "HTML"

	_, _ = h.Bot.Send(msg)
}
