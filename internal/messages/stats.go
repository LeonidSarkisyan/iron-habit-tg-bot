package messages

import (
	"HabitsBot/pkg/utils"
	"github.com/rs/zerolog/log"
	"math/rand"
)

const (
	HabitCancelMsg = "Привычка на сегодня отменена  ✅"
)

var congratulations = map[int]string{
	0: "Вы успешно выполнили привычку. Поздравляем! 😊",
	1: "Отличная работа! Привычка успешно выполнена. 👍",
	2: "Мы рады сообщить вам, что привычка выполнена. 🎉",
	3: "Поздравляем с достижением цели! Привычка выполнена. 🌟",
	4: "Ваш труд увенчался успехом. Привычка выполнена. 🚀",
	5: "Мы выражаем вам признание за выполнение привычки. 👏",
	6: "Вы продемонстрировали высокий уровень самодисциплины. Привычка выполнена. 💪",
	7: "Отличный результат! Привычка успешно выполнена. 🥳",
	8: "Мы благодарим вас за ваше усердие. Привычка выполнена. 🙌",
	9: "Поздравляем с выполнением привычки! Вы на верном пути к успеху. 🎊",
}

func RandomCreateCompleteHabitMsg() string {
	seed, err := utils.GenerateRandomSeed()

	if err != nil {
		log.Error().Err(err).Msg("ошибка при генерация сида")
	}

	r := rand.New(rand.NewSource(seed))

	randIndex := r.Intn(len(congratulations))
	congratulation := congratulations[randIndex]

	return congratulation
}
