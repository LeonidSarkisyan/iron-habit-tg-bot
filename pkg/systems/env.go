package systems

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func MustLoadEnv(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Fatal().Err(err).Msg("ошибка загрузки .env файла")
	}
}
