package tunnel

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients []*Client
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{clients: make([]*Client, 0, 2)}
}

func (h *Hub) HandleConnection(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.clients) >= 2 {
		log.Println("Max clients reached. Closing new connection.")
		conn.WriteMessage(websocket.TextMessage, []byte("Max clients connected. Try again later."))
		conn.Close()
		return
	}

	client := NewClient(conn, h)
	h.clients = append(h.clients, client)
	go client.Listen()
}

func (h *Hub) RelayMessage(from *Client, msg []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range h.clients {
		if client != from {
			client.Send(msg)
		}
	}
}

func (h *Hub) RemoveClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i, client := range h.clients {
		if client == c {
			h.clients = append(h.clients[:i], h.clients[i+1:]...)
			break
		}
	}
}
