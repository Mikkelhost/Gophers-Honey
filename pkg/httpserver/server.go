package httpserver

import (
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/api"
	"github.com/Mikkelhost/Gophers-Honey/pkg/websocket"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var DEV = true

var configured bool

func RunServer() {
	fmt.Println("Starting websocket")
	r := mux.NewRouter()
	websocket.SetupRouter(r)
	api.SetupRouters(r)
	if !DEV {
		log.Fatal(http.ListenAndServeTLS(":8443", "certs/nginx-selfsigned.crt", "certs/nginx-selfsigned.key", r))
	}
	log.Fatal(http.ListenAndServe(":8000", r))
}