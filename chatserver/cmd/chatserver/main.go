package main

import (
	"log"
	"net/http"

	"onrabble.com/chatserver/internal/cache"
	"onrabble.com/chatserver/internal/db"
	"onrabble.com/chatserver/internal/hub"
	"onrabble.com/chatserver/internal/server"

	"github.com/gorilla/websocket"
	"github.com/valkey-io/valkey-go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (adjust as needed)
	},
}

func main() {
	// Connect to the database
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Database connection established.")

	// Initialize Valkey client
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"valkey:6379"},
	})
	if err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}

	// Initialize the message cache
	messageCache := &cache.MessageCache{
		ValkeyClient: client,
		DB:           conn,
	}
	messageCache.StartPeriodicFlush()

	// Create a new Hub instance
	h := hub.NewHub(conn, messageCache)

	// Start the Hub in a separate goroutine
	go h.Run()

	if err != nil {
		log.Fatalf("Failed to connected to database")
	}

	// Create the Server instance and pass the Hub
	srv, err := server.New("0.0.0.0:8080", h, conn, messageCache)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Start the HTTP server
	log.Println("Starting server on :8080")
	if err := srv.HttpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
