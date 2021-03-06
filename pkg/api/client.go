package api

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type     int    `json:"type"`
	Body     string `json:"body,omitempty"`
	DeviceID uint32 `json:"device_id,omitempty"`
}

//Read
//keeps the ws connection alive
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
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
