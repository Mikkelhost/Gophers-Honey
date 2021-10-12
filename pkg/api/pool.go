package api

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Heartbeat  chan uint32
	Clients    map[*Client]bool
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Heartbeat:  make(chan uint32),
		Clients:    make(map[*Client]bool),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			log.Logger.Debug().Msgf("Size of Connection Pool: %d", len(pool.Clients))
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			log.Logger.Debug().Msgf("Size of Connection Pool: %d", len(pool.Clients))
			break
		case id := <-pool.Heartbeat:
			log.Logger.Debug().Msg("Sending heartbeat notification to clients")
			for client, _ := range pool.Clients{
				client.Conn.WriteJSON(Message{Type: 2, DeviceID: id})
			}
		}
	}
}
