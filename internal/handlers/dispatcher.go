package handlers

import (
	"HabitsBot/internal/filters"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type Filter func(update *tgbotapi.Update) bool

type HandlerFilter struct {
	Filter  Filter
	Handler HandlerFunc
}

type Dispatcher struct {
	HabitBot        *HabitBot
	handlersFilters []HandlerFilter
	Routers         []*Router
}

func NewDispatcher(bot *HabitBot) *Dispatcher {
	return &Dispatcher{HabitBot: bot, handlersFilters: make([]HandlerFilter, 0), Routers: make([]*Router, 0)}
}

func (d *Dispatcher) Message(handler HandlerFunc, filters_ ...filters.Filter) {
	messageFilter := filters.F(func(update *tgbotapi.Update) bool {

		f := filters.F(filters_...)

		return f(update) && update.Message != nil
	})

	d.register(messageFilter, handler)
}

func (d *Dispatcher) CallBackQuery(handler HandlerFunc, filters_ ...filters.Filter) {
	callBackFilter := filters.F(func(update *tgbotapi.Update) bool {

		f := filters.F(filters_...)

		return f(update) && update.CallbackQuery != nil
	})

	d.register(callBackFilter, handler)
}

func (d *Dispatcher) FSMState(state string, handler HandlerFunc, filters_ ...filters.Filter) {
	filterWithFSMState := func(update *tgbotapi.Update) bool {
		FSMState := d.HabitBot.FSM(update).Current()

		f := filters.F(filters_...)

		return f(update) && state == FSMState
	}

	d.register(filterWithFSMState, handler)
}

func (d *Dispatcher) PassHandlers(update *tgbotapi.Update) {
	done := false

	for _, handlerFilter := range d.handlersFilters {
		if handlerFilter.Filter(update) {
			handlerFilter.Handler(update)
			done = true
			break
		}
	}

	if !done {
		for _, router := range d.Routers {
			done = router.PassHandlers(update)
			if done {
				break
			}
		}
	}
}

func (d *Dispatcher) IncludeRouter(router *Router) {
	d.Routers = append(d.Routers, router)
}

func (d *Dispatcher) Polling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := d.HabitBot.Bot.GetUpdatesChan(u)

	for update := range updates {
		log.Info().Msgf("%+v", update)

		if update.CallbackQuery != nil {
			log.Info().Str("call back data", update.CallbackQuery.Data).Send()
		}

		d.PassHandlers(&update)
	}
}

func (d *Dispatcher) register(filter func(update *tgbotapi.Update) bool, handler HandlerFunc) {
	d.handlersFilters = append(d.handlersFilters, HandlerFilter{Filter: filter, Handler: handler})
}
