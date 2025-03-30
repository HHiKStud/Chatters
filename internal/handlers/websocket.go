package handlers

import (
	"log"
	"net/http"
	"time"

	"chi/internal/services"

	"github.com/gorilla/websocket"
)

// Глобальные переменные
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			log.Println("WebSocket connection attempt from:", origin)
			return true // Для разработки
		},
		HandshakeTimeout: 109 * time.Second,
	}
)

func ServeWs(h *services.Hub, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	username, ok := r.Context().Value("username").(string)
	if !ok {
		log.Println("Username not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("New WebSocket connection from %s", username) // Добавлено логирование

	client := &services.Client{
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Username: username,
	}

	h.Register <- client
	log.Printf("Client %s registered", username) // Добавлено логирование

	go client.WritePump()
	go client.ReadPump(h)
}
