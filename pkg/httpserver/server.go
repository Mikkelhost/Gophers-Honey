package httpserver

import (
	"github.com/Mikkelhost/Gophers-Honey/pkg/api"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/gorilla/mux"
	sLog "log"
	"net/http"
)

var DEV = true

var configured bool

//RunServer
//Runs the http/https server on the mux router
func RunServer() {
	log.Logger.Debug().Msgf("Starting websocket")
	r := mux.NewRouter()
	api.SetupWs(r)
	api.SetupRouters(r)
	if !DEV {
		sLog.Fatal(http.ListenAndServeTLS(":8443", "certs/nginx-selfsigned.crt", "certs/nginx-selfsigned.key", r))
	}
	sLog.Fatal(http.ListenAndServe(":8000", r))
}
