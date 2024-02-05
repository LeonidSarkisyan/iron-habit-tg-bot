package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (h *HabitBot) createOrUpdateSliceMetadata(update *tgbotapi.Update, key string, value any) {
	sliceValue, exists := h.FSM(update).Metadata(key)
	if exists {
		sliceValue = append(sliceValue.([]any), value)
		h.FSM(update).SetMetadata(key, sliceValue)
	} else {
		h.FSM(update).SetMetadata(key, []any{value})
	}
}
