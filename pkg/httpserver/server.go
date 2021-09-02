package httpserver

import (
	"fmt"
	"github.com/GeekMuch/GoHoney/pkg/api"
	"github.com/GeekMuch/GoHoney/pkg/websocket"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var configured bool

func RunServer() {
	fmt.Println("Starting websocket")
	r := mux.NewRouter()
	websocket.SetupRouter(r)
	api.SetupRouters(r)
	log.Fatal(http.ListenAndServeTLS(":8443", "certs/nginx-selfsigned.crt", "certs/nginx-selfsigned.key", r))
}