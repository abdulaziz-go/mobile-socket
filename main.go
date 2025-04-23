package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"mobile-socket/tunnel"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	hub := tunnel.NewHub()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}
		hub.HandleConnection(conn)
	})

	fmt.Println("WebSocket tunnel server started on :8080")
	http.ListenAndServe(":8080", nil)
}
