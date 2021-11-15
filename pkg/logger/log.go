package logger

import (
	"github.com/rs/zerolog"
	"log"
	"os"
)

var Logger zerolog.Logger

func InitLog(debug bool) {
	consoleWrite := zerolog.ConsoleWriter{Out: os.Stderr}
	Logger = zerolog.New(consoleWrite).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	if debug {
		Logger = zerolog.New(consoleWrite).With().Caller().Timestamp().Logger().Level(zerolog.DebugLevel)

	}
	stdLogger := Logger
	log.SetFlags(0)
	log.SetOutput(stdLogger)
}

func GetLogger() *zerolog.Logger {
	return &Logger
}
