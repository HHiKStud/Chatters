package database

import (
	"chi/internal/config"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	//importing config as it is required for the function
	cfg := config.LoadConfig()

	_, err := NewDatabase(cfg)
	if err != nil {
		t.Errorf("NewDatabase returned error: %v", err)
	}
}

func TestClose(t *testing.T) {
	cfg := config.LoadConfig()
	db, err := NewDatabase(cfg)
	if err != nil {
		t.Errorf("NewDatabase returned error: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestInit(t *testing.T) {
	cfg := config.LoadConfig()
	db, err := NewDatabase(cfg)
	if err != nil {
		t.Errorf("NewDatabase returned error: %v", err)
	}

	err = db.Init()
	if err != nil {
		t.Errorf("Init returned error: %v", err)
	}

	// Checking if tables exist
	var count int

	// users table
	err = db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Errorf("Error while cheking if table 'users' exist: %v", err)
	} else if count == 0 {
		t.Errorf("Table 'users' doesn't exist.")
	}

	// messages table
	err = db.DB.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
	if err != nil {
		t.Errorf("Error while cheking if table 'messages' exist: %v", err)
	} else if count == 0 {
		t.Errorf("Table 'messages' doesn't exist")
	}
}
