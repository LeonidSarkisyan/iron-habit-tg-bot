package messages

import "fmt"

const (
	HabitCreateErrorMsg = "Ошибка сохранения привычки. Попробуйте ещё раз"

	InvalidInputMsg                 = "Введено некорректное значение. Попробуйте ещё раз"
	InvalidRangeInputWarningTimeMsg = "Число должно быть в диапазоне от 0 до 60"

	CancelHabitErrorMsg = "Не удалось отменить выполнение привычки. Попробуйте ещё раз"

	InvalidRejectionMsg = "Введено некорректное значение. Попробуйте ещё раз"

	RejectionCreateErrorMsg = "Ошибка при создании отмены привычки. Попробуйте ещё раз"
)

func CancelHabitMsg(habitName, rejectionText string) string {
	return fmt.Sprintf(
		"Вы отменили выполнение привычки <b>\"%s\"</b> 😤\n\n<b>Причина:</b> %s", habitName, rejectionText)
}
