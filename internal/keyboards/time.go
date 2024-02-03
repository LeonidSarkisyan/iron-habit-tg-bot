package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func TimePickerKeyboard(times []string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for h := 0; h <= 23; h++ {
		for m := 0; m < 60; m += 30 {
			buttonText := fmt.Sprintf("%02d:%02d", h, m)
			callbackData := fmt.Sprintf("%02d:%02d", h, m)

			// Создаем кнопку и добавляем в текущую строку
			for _, time := range times {
				if strings.Replace(time, "time__", "", 1) == callbackData {
					buttonText += "  " + DoneEmoji
					break
				}
			}

			btn := tgbotapi.NewInlineKeyboardButtonData(buttonText, "time__"+callbackData)
			rowIndex := len(rows) - 1
			if rowIndex == -1 || len(rows[rowIndex]) == 4 {
				// Если текущая строка полна, создаем новую строку
				rows = append(rows, []tgbotapi.InlineKeyboardButton{btn})
			} else {
				// Иначе добавляем к текущей строке
				rows[rowIndex] = append(rows[rowIndex], btn)
			}
		}
	}

	if len(times) > 0 {
		btn := tgbotapi.NewInlineKeyboardButtonData("Продолжить  ⏩", "time__continue")
		rows = append(rows, []tgbotapi.InlineKeyboardButton{btn})
	}

	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}