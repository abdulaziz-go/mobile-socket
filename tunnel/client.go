package tunnel

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	hub  *Hub
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		conn: conn,
		hub:  hub,
	}
}

func (c *Client) Listen() {
	defer func() {
		c.hub.RemoveClient(c)
		c.conn.Close()
		log.Println("Client disconnected")
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		c.hub.RelayMessage(c, msg)
	}
}

func (c *Client) Send(msg []byte) {
	err := c.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println("Write error:", err)
	}
}
