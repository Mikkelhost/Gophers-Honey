package api

import (
	"github.com/Mikkelhost/Gophers-Honey/pkg/config"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/gorilla/mux"
	"os"
)

/*
Here the main router will be passed to each of the different type of API
subrouters.
*/
var (
	SECRET_KEY = getenv("SECRET_KEY", "UWKvPGDYd2zmAmbYQB2K")
)

var conf *config.Config

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	log.Logger.Debug().Msgf("Env %s not set, using default of %s", key, fallback)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func SetupRouters(r *mux.Router, c *config.Config) {
	conf = c
	devicesSubrouter(r)
	usersSubrouter(r)
	configSubrouter(r)
	logsSubrouter(r)
}
