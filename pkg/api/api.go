package api

import (
	"fmt"
	"github.com/gorilla/mux"
	logs "log"
	"os"
)

/*
Here the main router will be passed to each of the different type of API
subrouters.
*/

var (
	SECRET_KEY = getenv("SECRET_KEY", "UWKvPGDYd2zmAmbYQB2K")
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		logs.Println(fmt.Sprintf("Environment variable %s not set. Defaulting to %s", key, fallback))
		return fallback
	}
	return value
}

func SetupRouters(r *mux.Router) {
	devicesSubrouter(r)
	usersSubrouter(r)
	configSubrouter(r)
	logsSubrouter(r)
	imageSubRouter(r)
}
