package models

import (
	"chi/internal/config"
	"chi/internal/database"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMessageHandler(t *testing.T) {
	cfg := config.LoadConfig()
	db, _ := database.NewDatabase(cfg)
	defer db.Close()

	handler := NewMessageHandler(db.DB)
	if handler == nil {
		t.Errorf("Expected not 'nil' value")
	}
}

func TestGetMessageHandler(t *testing.T) {
	cfg := config.LoadConfig()
	db, _ := database.NewDatabase(cfg)
	defer db.Close()

	// initializing a handler with test db(in my case prod db)
	handler := &MessageHandler{db: db.DB}

	// setting up an http server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.GetMessagesHandler(w, r)
	}))
	defer testServer.Close()

	response, err := http.Get(testServer.URL)
	if err != nil {
		t.Fatalf("Error while making a query: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status-code 200 OK, but got %d", response.StatusCode)
	}
}
