package systems

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func SetupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.DebugLevel).With().Caller().Logger()

	log.Info().Msg("логирование включено")
}
