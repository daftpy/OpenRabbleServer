package handlers

import (
	database "chatserver/internal/db"
	"chatserver/internal/messages"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

		response := messages.NewUserSearchResultMessage(payload)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
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
