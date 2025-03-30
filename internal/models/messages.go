package models

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Message struct {
	ID      int       `json:"id" db:"id"`
	UserID  int       `json:"user_id" db:"user_id"`
	Content string    `json:"content" db:"content"`
	SentAt  time.Time `json:"sent_at" db:"sent_at"`

	// Это поле не хранится в бд, нужно для удобства
	Username string `json:"username,omitempty" db:""`
}

type MessageHandler struct {
	db *sql.DB
}

func NewMessageHandler(db *sql.DB) *MessageHandler {
	return &MessageHandler{db: db}
}

func (h *MessageHandler) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем правильный Content-Type
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Получаем username текущего пользователя из контекста
	currentUser, _ := r.Context().Value("username").(string)

	rows, err := h.db.Query(`
        SELECT m.id, u.username, m.content, m.sent_at 
        FROM messages m
        JOIN users u ON m.user_id = u.id
        ORDER BY m.sent_at ASC
        LIMIT 100
    `)
	if err != nil {
		log.Printf("DB query error: %v", err)
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var messages []struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Text     string `json:"text"`
		Time     string `json:"time"`
		IsMine   bool   `json:"is_mine"`
	}

	for rows.Next() {
		var msg struct {
			ID       int
			Username string
			Content  string
			SentAt   time.Time
		}

		if err := rows.Scan(&msg.ID, &msg.Username, &msg.Content, &msg.SentAt); err != nil {
			log.Printf("DB scan error: %v", err)
			continue
		}

		messages = append(messages, struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Text     string `json:"text"`
			Time     string `json:"time"`
			IsMine   bool   `json:"is_mine"`
		}{
			ID:       msg.ID,
			Username: msg.Username,
			Text:     msg.Content,
			Time:     msg.SentAt.Format(time.RFC3339),
			IsMine:   msg.Username == currentUser, // Check if the message is from the current user
		})
	}

	// Кодируем JSON с обработкой ошибок
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("JSON encode error: %v", err)
		http.Error(w, `{"error": "JSON encoding failed"}`, http.StatusInternalServerError)
	}
}
