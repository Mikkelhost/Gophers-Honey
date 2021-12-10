package logger

import (
	"github.com/rs/zerolog"
	"os"
)

var Logger zerolog.Logger

//InitLog
//Sets up zerolog for use within the project
func InitLog(debug bool){
	consoleWrite := zerolog.ConsoleWriter{Out: os.Stderr}
	Logger = zerolog.New(consoleWrite).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	if debug {
		Logger = zerolog.New(consoleWrite).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	}
}