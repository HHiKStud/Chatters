package config

import (
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	cfg := LoadConfig()

	// testing if it even loads
	if cfg == nil {
		t.Errorf("Expected non-nil config")
	}

	/// field tests
	if cfg.ServerAddress == "" {
		t.Errorf("Expected non-empty server address")
	}

	if cfg.ServerAddress != ":8080" {
		t.Errorf("Expected address: 8080, got %s", cfg.ServerAddress)
	}

	if cfg.JWTSecret != "your-secret-key" {
		t.Errorf("Expected key: 'your-secret-key', got %s", cfg.JWTSecret)
	}

	if cfg.TokenExpiration != 24*time.Hour {
		t.Errorf("Expected token expiration: 24 hours, got %v", cfg.TokenExpiration)
	}

	if cfg.DBHost != "127.0.0.1" {
		t.Errorf("Expected host '127.0.0.1 or localhost', got %s", cfg.DBHost)
	}

	if cfg.DBName != "chat_app" {
		t.Errorf("Expected db name to be 'chat_app', got %s", cfg.DBName)
	}

	if cfg.DBPort != "5432" {
		t.Errorf("Expected default db port - (5432), got %s", cfg.DBPort)
	}

	if cfg.DBUser != "postgres" {
		t.Errorf("Expected default db user - (postgres), got %s", cfg.DBUser)
	}

	if cfg.DBPassword != "password" {
		t.Errorf("Expected password - 'password', got %s", cfg.DBPassword)
	}
}
