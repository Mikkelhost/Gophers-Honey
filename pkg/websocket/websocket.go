package websocket

import (

	//"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{
		Conn: conn,
		Pool: pool,
	}

	log.Println("Client successfully connected...")
	pool.Register <- client
	client.Read()
}

func SetupRouter(r *mux.Router) {
	pool := NewPool()
	go pool.Start()

	ws := r.PathPrefix("/ws").Subrouter()
	pi := r.PathPrefix("/pi").Subrouter()

	ws.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})

	pi.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}