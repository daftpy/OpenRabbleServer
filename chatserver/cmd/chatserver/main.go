package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (adjust as needed)
	},
}

func main() {
	mux := http.NewServeMux()

	// Discovery Handler
	mux.HandleFunc("/discovery", func(w http.ResponseWriter, r *http.Request) {
		KChostname := os.Getenv("KC_HOSTNAME")
		chatClientName := os.Getenv("CHAT_CLIENT_NAME")
		realmName := os.Getenv("REALM_NAME")
		log.Printf("KC_HOSTNAME: %s", KChostname)
		log.Printf("REALM_NAME: %s", realmName)

		url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect", KChostname, realmName)
		response := map[string]string{
			"auth_url":    url + "/auth",
			"chat_client": chatClientName,
			"chat_url":    "wss://chat.localhost/ws",
			"token_url":   url + "/token",
			"server_name": "OnRabble",
			"server_id":   "Placeholder",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// WebSocket Handler
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Check for a token
		token := r.URL.Query().Get("token")

		if len(token) == 0 {
			log.Printf("No token, no connection.")
			return
		}
		log.Println("Token:", token)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}
		defer conn.Close()

		log.Println("New WebSocket connection established")
		for {
			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocket read error:", err)
				break
			}
			log.Printf("Received message: %s\n", msg)

			// Echo the message back
			if err := conn.WriteMessage(messageType, msg); err != nil {
				log.Println("WebSocket write error:", err)
				break
			}
		}
	})

	// Start the HTTP server
	log.Println("Starting server on :8080")
	http.ListenAndServe("0.0.0.0:8080", mux)
}
