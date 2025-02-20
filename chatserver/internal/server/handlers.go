package server

import (
	database "chatserver/internal/db"
	"chatserver/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleDiscovery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func HandleChannels(db *pgxpool.Pool, s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			channels, err := database.FetchChannels(db)
			if err != nil {
				log.Println("Failed to load channels from database:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Set the content type and encode the response as JSON
			w.Header().Set("Content-Type", "application/json")
			response := map[string][]models.Channel{"channels": channels}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Println("Failed to encode channels response:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

		case http.MethodPost:
			var request struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}

			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
				return
			}

			if request.Name == "" {
				http.Error(w, "Channel name is required", http.StatusBadRequest)
				return
			}

			err := createChannel(db, request.Name, request.Description)
			if err != nil {
				log.Println("Failed to create channel:", err)
				http.Error(w, "Failed to create channel", http.StatusInternalServerError)
				return
			}

			log.Printf("Channel '%s' created successfully", request.Name)
			// broadcast new channel message s.hub.Broadcast()
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "Channel created", "name": request.Name})

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
