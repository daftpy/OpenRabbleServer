package server

import (
	database "chatserver/internal/db"
	"chatserver/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleDiscovery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		KChostname := os.Getenv("KC_HOSTNAME")
		chatClientName := os.Getenv("CHAT_CLIENT_NAME")
		realmName := os.Getenv("REALM_NAME")
		hostname := os.Getenv("PUBLIC_HOSTNAME")

		log.Printf("KC_HOSTNAME: %s", KChostname)
		log.Printf("REALM_NAME: %s", realmName)

		url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect", KChostname, realmName)
		response := map[string]string{
			"auth_url":    url + "/auth",
			"chat_client": chatClientName,
			"chat_url":    "wss://chat." + hostname + "/ws",
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

func HandleMessages(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Extract query parameters
			channel := r.URL.Query().Get("channel")
			keyword := r.URL.Query().Get("keyword") // Optional search keyword
			limitStr := r.URL.Query().Get("limit")
			offsetStr := r.URL.Query().Get("offset")

			limit := 50 // Default: 50 messages per request
			offset := 0 // Default: No offset (start from latest messages)

			if limitStr != "" {
				parsedLimit, err := strconv.Atoi(limitStr)
				if err != nil || parsedLimit <= 0 {
					http.Error(w, "Invalid 'limit' query parameter", http.StatusBadRequest)
					return
				}
				limit = parsedLimit
			}

			if offsetStr != "" {
				parsedOffset, err := strconv.Atoi(offsetStr)
				if err != nil || parsedOffset < 0 {
					http.Error(w, "Invalid 'offset' query parameter", http.StatusBadRequest)
					return
				}
				offset = parsedOffset
			}

			// Fetch messages from the database
			messages, err := database.FetchMessages(db, channel, keyword, limit, offset)
			if err != nil {
				log.Printf("Failed to fetch messages for channel '%s': %v", channel, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string][]models.ChatMessage{"messages": messages})

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
