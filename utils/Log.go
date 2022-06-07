package utils

import (
	"os"

	"github.com/rs/zerolog"
)

// Log is our logging interface
var Log zerolog.Logger = setupLogger()

func setupLogger() zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return logger
}
