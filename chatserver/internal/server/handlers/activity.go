package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	database "onrabble.com/chatserver/internal/db"
	"onrabble.com/chatserver/internal/messages/api"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleRecentActivity(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		activity, err := database.FetchSessionActivity(db, userID)
		if err != nil {
			log.Printf("Failed to fetch recent acivity for %s: %v", userID, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		msg := api.NewSessionActivityMessage(activity)
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
		msg := api.NewMessageCountByChannelMessage(activity)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	}
}
