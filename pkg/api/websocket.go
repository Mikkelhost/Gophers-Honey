package api

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	//"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

var ClientPool *Pool

func serveWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Logger.Warn().Msgf("%s", err)
	}

	client := &Client{
		Conn: conn,
		Pool: pool,
	}

	log.Logger.Debug().Msg("Client successfully connected...")
	pool.Register <- client
	client.Read()
}

func SetupWs(r *mux.Router) {
	ClientPool = NewPool()
	go ClientPool.Start()

	ws := r.PathPrefix("/ws").Subrouter()

	ws.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		serveWs(ClientPool, w, r)
	})
}
