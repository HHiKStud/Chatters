package main

import (
	"chi/internal/config"
	"chi/internal/database"
	"chi/internal/handlers"
	"chi/internal/models"
	"chi/internal/services"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	// Инициализация БД
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Инициализация обработчиков
	handlers.InitAuthHandlers(cfg, db.DB)

	messageHandler := models.NewMessageHandler(db.DB)

	hub := services.NewHub(db.DB)
	go hub.Run()

	// HTTP роуты
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/messages", handlers.AuthMiddleware(messageHandler.GetMessagesHandler))

	http.HandleFunc("/ws", handlers.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWs(hub, w, r)
	}))

	// Раздача статики
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
