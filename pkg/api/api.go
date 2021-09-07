package api

import (
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

/*
Here the main router will be passed to each of the different type of API
subrouters.
*/
var (
	SECRET_KEY = getenv("SECRET_KEY", "UWKvPGDYd2zmAmbYQB2K")
	DEBUG      = false
)

func getenv(key, fallback string) string {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	value := os.Getenv(key)
	log.Debug().Msgf("Env %s not set, using default of %s", key, fallback)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func SetupRouters(r *mux.Router) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	devicesSubrouter(r)
	usersSubrouter(r)
}
