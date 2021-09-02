package websocket

import(
  "fmt"
)


type Pool struct {
  Register    chan *Client
  Unregister  chan *Client
  RegisterPi  chan *RPi
  UnregisterPi  chan *RPi
  Clients     map[*Client]bool
  RPis        map[*RPi]bool
}

func NewPool() *Pool {
  return &Pool{
    RegisterPi:  make(chan *RPi),
    UnregisterPi:  make(chan *RPi),
    Register: make(chan *Client),
    Unregister: make(chan *Client),
    Clients: make(map[*Client]bool),
    RPis: make(map[*RPi]bool),
  }
}

func (pool *Pool) Start() {
    for {
        select {
        case client := <-pool.Register:
            pool.Clients[client] = true
            fmt.Println("Size of Connection Pool: ", len(pool.Clients))
            for client, _ := range pool.Clients {
                fmt.Println(client)
                client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
            }
            break
        case client := <-pool.Unregister:
            delete(pool.Clients, client)
            fmt.Println("Size of Connection Pool: ", len(pool.Clients))
            for client, _ := range pool.Clients {
                client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
            }
            break
        }
    }
}
