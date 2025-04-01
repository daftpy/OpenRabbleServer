package handlers

import (
	database "chatserver/internal/db"
	"chatserver/internal/messages/api"
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
		payload := api.UserSearchResultPayload{
			Users: users,
		}

		response := api.NewUserSearchResultMessage(payload)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func HandleBanUser(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// temp owner_id
			ownerID := "ace4e8be-d2a2-46d7-9c9e-57f04f915835"

			var request struct {
				BanishedID string  `json:"banished_id"` // The user being banned
				Reason     *string `json:"reason"`      // Reason for the ban (optional)
				Duration   *int    `json:"duration"`    // Duration in hours (optional, nil = permanent)
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

			// Unwrap duration
			duration := 0
			if request.Duration != nil {
				duration = *request.Duration
			}

			// Unwrap reason
			reason := ""
			if request.Reason != nil {
				reason = *request.Reason
			}

			// Ban the user
			err := database.BanUser(db, ownerID, request.BanishedID, reason, duration)
			if err != nil {
				log.Printf("Failed to ban user: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			log.Printf("User %s banned by %s. Reason: %v, Duration: %s",
				request.BanishedID, ownerID, reason,
				func() string {
					if duration == 0 {
						return "Permanent"
					}
					return fmt.Sprintf("%d hours", duration)
				}(),
			)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message":  "User banned successfully",
				"banished": request.BanishedID,
				"reason":   reason,
				"duration": func() string {
					if duration == 0 {
						return "Permanent"
					}
					return fmt.Sprintf("%d hours", duration)
				}(),
			})

		case http.MethodDelete:
			banIDStr := r.URL.Query().Get("ban_id")
			if banIDStr == "" {
				http.Error(w, "Missing ban_id parameter", http.StatusBadRequest)
				return
			}

			banID, err := strconv.Atoi(banIDStr)
			if err != nil {
				http.Error(w, "Invalid ban_id", http.StatusBadRequest)
				return
			}

			err = database.PardonUser(db, banID)
			if err != nil {
				log.Printf("Failed to pardon user: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			log.Printf("Ban ID %d pardoned", banID)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "User pardoned successfully",
				"ban_id":  strconv.Itoa(banID),
			})

		default:
			log.Printf("Wrong method for bans: %s", r.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
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
		response := api.NewBanRecordsResultMessage(banRecords, hasMore)

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
