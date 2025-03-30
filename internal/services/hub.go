package services

import (
	"database/sql"
	"log"
)

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	DB         *sql.DB
}

func NewHub(db *sql.DB) *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		DB:         db,
	}
}

func (h *Hub) SaveMessage(UserID int, Content string) error {
	_, err := h.DB.Exec("INSERT INTO messages (user_id, content) VALUES ($1, $2)",
		UserID, Content)
	if err != nil {
		log.Printf("SaveMessage error: %v", err)
	}
	return err
}

func (h *Hub) BroadcastMessage(message []byte) {
	log.Printf("Broadcasting message: %s", string(message))
	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.Clients, client)
			log.Println("Client disconnected due to slow connection")
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("Clients registered. Total: ", len(h.Clients))

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Println("Clients UNregistered. Total: ", len(h.Clients))
			}

		case message := <-h.Broadcast:
			log.Printf("Broadcasting to %d clients: %s", len(h.Clients), string(message))
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
					log.Println("Client disconnected (slow connection)")
				}
			}
		}
	}
}
