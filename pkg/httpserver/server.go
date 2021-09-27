package httpserver

import (
	"github.com/Mikkelhost/Gophers-Honey/pkg/api"
	"github.com/Mikkelhost/Gophers-Honey/pkg/config"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/websocket"
	"github.com/gorilla/mux"
	sLog "log"
	"net/http"
)

var DEV = true

var configured bool

func RunServer(c *config.Config) {
	log.Logger.Debug().Msgf("Starting websocket")
	r := mux.NewRouter()
	websocket.SetupRouter(r)
	api.SetupRouters(r, c)
	if !DEV {
		sLog.Fatal(http.ListenAndServeTLS(":8443", "certs/nginx-selfsigned.crt", "certs/nginx-selfsigned.key", r))
	}
	sLog.Fatal(http.ListenAndServe(":8000", r))
}
