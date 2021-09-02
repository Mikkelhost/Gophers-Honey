package api

import (
	"github.com/gorilla/mux"
)

/*
Here the main router will be passed to each of the different type of API
subrouters.
 */

func SetupRouters(r *mux.Router) {
	devicesSubrouter(r)
	usersSubrouter(r)
}

