package logger

import (
	"github.com/rs/zerolog"
	"os"
)

var Logger zerolog.Logger

func InitLog(){
	consoleWrite := zerolog.ConsoleWriter{Out: os.Stderr}
	multi := zerolog.MultiLevelWriter(consoleWrite, os.Stdout)
	Logger = zerolog.New(multi).With().Timestamp().Logger()
}