package server

import (
	database "chatserver/internal/db"
	"chatserver/internal/messages"
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
			channels := r.URL.Query()["channel"]
			keyword := r.URL.Query().Get("keyword")
			limitStr := r.URL.Query().Get("limit")
			offsetStr := r.URL.Query().Get("offset")

			limit := 50 // Default
			offset := 0 // Default

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

			userID := r.URL.Query().Get("user_id")

			// Fetch messages
			search_messages, hasMore, err := database.FetchMessages(db, userID, channels, keyword, limit, offset)
			if err != nil {
				log.Printf("Failed to fetch messages for channels '%v': %v", channels, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Wrap messages in the correct struct
			payload := messages.MessageSearchResultPayload{
				Messages: search_messages, // Fixing struct usage
				HasMore:  hasMore,
			}

			// Use NewMessageSearchResultMessage with correct payload
			responseMessage := messages.NewMessageSearchResultMessage(payload)

			// Send JSON response
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(responseMessage)

		case http.MethodDelete:
			// 1) Parse JSON body
			var body struct {
				IDs []int `json:"ids"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "Invalid JSON body", http.StatusBadRequest)
				return
			}
			log.Printf("Received these IDs to be deleted: %v", body.IDs)

			if len(body.IDs) == 0 {
				http.Error(w, "No IDs provided", http.StatusBadRequest)
				return
			}

			// 2) Call your bulk deletion function
			rowsDeleted, err := database.RemoveMessages(db, body.IDs)
			if err != nil {
				log.Printf("Failed to delete messages: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// 3) If zero rows are deleted, none of those IDs existed
			if rowsDeleted == 0 {
				http.Error(w, "No messages were deleted (IDs not found)", http.StatusNotFound)
				return
			}

			// 4) Respond with success
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func HandleUsers(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		users, err := database.FetchUsers(db, username)
		if err != nil {
			log.Printf("Failed to fetch users: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		payload := messages.UserSearchResultPayload{
			Users: users,
		}

		response := messages.NewUseerSearchResultMessage(payload)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func HandleRecentActivity(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		activity, err := database.FetchSessionActivity(db, userID)
		if err != nil {
			log.Printf("Failed to fetch recent acivity for %s: %v", userID, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		msg := messages.NewSessionActivityMessage(activity)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	}
}

func HandleChannelActivity(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activity, err := database.FetchMessageCountByChannel(db)
		if err != nil {
			log.Printf("Failed to fetch channel message counts: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		msg := messages.NewMessageCountByChannelMessage(activity)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	}
}

func HandleBanUser(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Wrong method for bans")
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// temp owner_id
		ownerID := "ace4e8be-d2a2-46d7-9c9e-57f04f915835"

		var request struct {
			BanishedID string `json:"banished_id"` // The user being banned
			Reason     string `json:"reason"`      // Reason for the ban
			Duration   *int   `json:"duration"`    // Duration in hours (optional, nil means permanent)
		}

		// Parse request body
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Printf("Failed to decode JSON")
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if ownerID == "" || request.BanishedID == "" {
			log.Printf("Failed to read banishedID: %s", request.BanishedID)
			http.Error(w, "Both owner_id and banished_id are required", http.StatusBadRequest)
			return
		}

		// Determine duration (default to 0 for permanent bans)
		duration := 0
		if request.Duration != nil {
			duration = *request.Duration
		}

		// Call the `BanUser` function from your database package
		err := database.BanUser(db, ownerID, request.BanishedID, request.Reason, duration)
		if err != nil {
			log.Printf("Failed to ban user: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Printf("User %s banned by %s. Reason: %s, Duration: %s",
			request.BanishedID, ownerID, request.Reason,
			func() string {
				if duration == 0 {
					return "Permanent"
				}
				return fmt.Sprintf("%d hours", duration)
			}(),
		)

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":  "User banned successfully",
			"banished": request.BanishedID,
			"reason":   request.Reason,
			"duration": func() string {
				if duration == 0 {
					return "Permanent"
				}
				return fmt.Sprintf("%d hours", duration)
			}(),
		})
	}
}

func HandleBanRecords(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit := 50 // Default
		offset := 0 // Default

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

		banRecords, hasMore, err := database.FetchBanRecords(db, limit, offset)
		if err != nil {
			log.Printf("Failed to fetch ban records: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create message payload
		response := messages.NewBanRecordsResultMessage(banRecords, hasMore)

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
