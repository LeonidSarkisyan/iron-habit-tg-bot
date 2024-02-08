package handlers

import (
	"HabitsBot/internal/filters"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router Dispatcher

func NewRouter(habitBot *HabitBot) *Router {
	return &Router{
		habitBot:        habitBot,
		handlersFilters: make([]HandlerFilter, 0),
		Routers:         make([]*Router, 0),
	}
}

func (r *Router) Message(handler HandlerFunc, filters_ ...filters.Filter) {
	messageFilter := filters.F(func(update *tgbotapi.Update) bool {

		f := filters.F(filters_...)

		return f(update) && update.Message != nil
	})

	r.register(messageFilter, handler)
}

func (r *Router) CallBackQuery(handler HandlerFunc, filters_ ...filters.Filter) {
	callBackFilter := filters.F(func(update *tgbotapi.Update) bool {

		f := filters.F(filters_...)

		return f(update) && update.CallbackQuery != nil
	})

	r.register(callBackFilter, handler)
}

func (r *Router) FSMState(state string, handler HandlerFunc, filters_ ...filters.Filter) {
	filterWithFSMState := func(update *tgbotapi.Update) bool {
		FSMState := r.habitBot.FSM(update).Current()

		f := filters.F(filters_...)

		return f(update) && state == FSMState && update.CallbackQuery == nil
	}

	r.register(filterWithFSMState, handler)
}

func (r *Router) PassHandlers(update *tgbotapi.Update) bool {
	done := false

	for _, handlerFilter := range r.handlersFilters {
		if handlerFilter.Filter(update) {
			handlerFilter.Handler(update)
			done = true
			break
		}
	}

	if !done {
		for _, router := range r.Routers {
			done = router.PassHandlers(update)
			if done {
				break
			}
		}
	}

	return done
}

func (r *Router) register(filter func(update *tgbotapi.Update) bool, handler HandlerFunc) {
	r.handlersFilters = append(r.handlersFilters, HandlerFilter{Filter: filter, Handler: handler})
}
