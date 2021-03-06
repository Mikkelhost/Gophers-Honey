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

//getenv
//Gets environment var or sets to default if not exists
func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		logs.Println(fmt.Sprintf("Environment variable %s not set. Defaulting to %s", key, fallback))
		return fallback
	}
	return value
}

//SetupRouters
//Sets up the different routers for the different api endpoints
func SetupRouters(r *mux.Router) {
	devicesSubrouter(r)
	usersSubrouter(r)
	configSubrouter(r)
	logsSubrouter(r)
	imageSubRouter(r)
}
