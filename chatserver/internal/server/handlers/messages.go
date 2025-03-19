package handlers

import (
	"chatserver/internal/cache"
	database "chatserver/internal/db"
	"chatserver/internal/messages"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleMessages(db *pgxpool.Pool, cache *cache.MessageCache) http.HandlerFunc {
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
			// Parse JSON body
			var body struct {
				IDs []int `json:"ids"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "Invalid JSON body", http.StatusBadRequest)
				return
			}

			if len(body.IDs) == 0 {
				http.Error(w, "No IDs provided", http.StatusBadRequest)
				return
			}

			// Call RemoveMessages to delete from DB and get cacheIDs
			rowsDeleted, cacheIDs, err := database.RemoveMessages(db, body.IDs)
			if err != nil {
				log.Printf("Failed to delete messages: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if rowsDeleted == 0 {
				http.Error(w, "No messages were deleted (IDs not found)", http.StatusNotFound)
				return
			}

			// Purge messages from the cache
			for _, cacheID := range cacheIDs {
				success := cache.DeleteCachedMessage(cacheID)
				if success {
					log.Printf("Purged message with cacheID %d from cache", cacheID)
				} else {
					log.Printf("Failed to purge message with cacheID %d from cache", cacheID)
				}
			}

			// Respond with success
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
