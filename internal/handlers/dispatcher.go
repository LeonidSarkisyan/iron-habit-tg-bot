package handlers

import (
	"HabitsBot/internal/filters"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Filter func(update *tgbotapi.Update) bool

type HandlerFilter struct {
	Filter  Filter
	Handler HandlerFunc
}

type Dispatcher struct {
	habitBot        *HabitBot
	handlersFilters []HandlerFilter
	Routers         []*Router
}

func NewDispatcher(bot *HabitBot) *Dispatcher {
	return &Dispatcher{habitBot: bot, handlersFilters: make([]HandlerFilter, 0), Routers: make([]*Router, 0)}
}

func (d *Dispatcher) Message(filter Filter, handler HandlerFunc) {
	messageFilter := filters.F(func(update *tgbotapi.Update) bool {
		return filter(update) && update.Message != nil
	})

	d.register(messageFilter, handler)
}

func (d *Dispatcher) CallBackQuery(filter Filter, handler HandlerFunc) {
	callBackFilter := filters.F(func(update *tgbotapi.Update) bool {
		return filter(update) && update.CallbackQuery != nil
	})

	d.register(callBackFilter, handler)
}

func (d *Dispatcher) FSMState(state string, filter Filter, handler HandlerFunc) {
	filterWithFSMState := func(update *tgbotapi.Update) bool {
		FSMState := d.habitBot.FSM(update).Current()

		return filter(update) && state == FSMState
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

func (d *Dispatcher) register(filter func(update *tgbotapi.Update) bool, handler HandlerFunc) {
	d.handlersFilters = append(d.handlersFilters, HandlerFilter{Filter: filter, Handler: handler})
}
