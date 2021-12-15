package httpserver

import (
	"github.com/Mikkelhost/Gophers-Honey/pkg/api"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/gorilla/mux"
	sLog "log"
	"net/http"
)

var configured bool

//RunServer
//Runs the http/https server on the mux router
func RunServer() {
	log.Logger.Debug().Msgf("Starting websocket")
	r := mux.NewRouter()
	api.SetupWs(r)
	api.SetupRouters(r)
	sLog.Fatal(http.ListenAndServe(":8000", r))
}
