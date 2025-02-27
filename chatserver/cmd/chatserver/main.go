package main

import (
	"chatserver/internal/db"
	"chatserver/internal/hub"
	"chatserver/internal/server"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (adjust as needed)
	},
}

func main() {
	// Connect to the database
	conn, err := db.Connect()

	// Create a new Hub instance
	h := hub.NewHub(conn)

	// Start the Hub in a separate goroutine
	go h.Run()

	if err != nil {
		log.Fatalf("Failed to connected to database")
	}

	// Create the Server instance and pass the Hub
	srv, err := server.New("0.0.0.0:8080", h, conn)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Start the HTTP server
	log.Println("Starting server on :8080")
	if err := srv.HttpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
