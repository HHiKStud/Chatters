package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Send     chan []byte
	Username string
}

func (c *Client) ReadPump(h *Hub) {
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()

	// Getting userId after connecting to DB
	var userID int
	err := h.DB.QueryRow("SELECT id FROM USERS WHERE username = $1", c.Username).Scan(&userID)
	if err != nil {
		log.Printf("Error getting userID: %v", err)
		return
	}

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		content := string(messageBytes)

		// Saving message to a DB
		if err := h.SaveMessage(userID, content); err != nil {
			log.Println("Failed to save message:", err) // logs
		}

		// Structure for being sent
		msg := struct {
			Username string    `json:"username"`
			Text     string    `json:"text"`
			Time     time.Time `json:"time"`
		}{
			Username: c.Username,
			Text:     content,
			Time:     time.Now(),
		}

		// Serializing
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Println("JSON marshal error:", err)
			continue
		}

		// Broadcasting
		h.Broadcast <- msgBytes
		log.Printf("Message broadcasted: %s", string(msgBytes)) // logs
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			w.Close()
		}
	}
}
