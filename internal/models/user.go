package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"` // Use hash ONLY in prod
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
