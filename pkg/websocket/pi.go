package websocket

import (
  "fmt"
  "log"
  "github.com/gorilla/websocket"
)

type RPi struct {
  ID string
  IP string
  Services string
  Configurable string
  Conn *websocket.Conn
  Pool *Pool
}

func (c *RPi) Read() {
    defer func() {
        c.Pool.UnregisterPi <- c
        c.Conn.Close()
    }()

    for {
        messageType, p, err := c.Conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        message := Message{Type: messageType, Body: string(p)}
        //c.Pool.Broadcast <- message
        fmt.Printf("Message Received: %+v\n", message)
    }
}
