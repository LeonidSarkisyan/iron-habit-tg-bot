package handlers

import (
	"HabitsBot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func (h *HabitBot) deleteMessage(update *tgbotapi.Update) {
	var chatID int64

	switch {
	case update.Message != nil:
		chatID = update.Message.Chat.ID
	case update.CallbackQuery != nil:
		chatID = update.CallbackQuery.Message.Chat.ID
	}

	messagesIDS, exists := h.FSM(update).Metadata("messages_ids")

	if exists {
		log.Info().Any("messages_ids", messagesIDS).Msg("нужно удалить старые сообщения")

		messagesIDS = utils.ConvertToIntSlice(messagesIDS.([]any))
		messagesIDS := utils.RemoveDuplicates(messagesIDS.([]int))

		for _, messageID := range messagesIDS {
			msgToDelete := tgbotapi.DeleteMessageConfig{
				ChatID:    chatID,
				MessageID: messageID,
			}

			_, err := h.Bot.Request(msgToDelete)

			if err != nil {
				log.Error().Err(err).Msg(err.Error())
			}
		}
	}
}
