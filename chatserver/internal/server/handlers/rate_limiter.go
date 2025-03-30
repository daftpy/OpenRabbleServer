package handlers

import (
	database "chatserver/internal/db"
	"chatserver/internal/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleRateLimiter(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rateLimiter, err := database.GetRateLimiterByID(db, 1)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to retrieve rate limiter: %v", err), http.StatusInternalServerError)
			}

			var resp struct {
				RateLimiter models.RateLimiter `json:"rate_limiter"`
			}

			resp.RateLimiter = rateLimiter

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)

		case http.MethodPatch:
			var payload struct {
				ID            int `json:"id"`
				MessageLimit  int `json:"message_limit"`
				WindowSeconds int `json:"window_seconds"`
			}

			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
				return
			}

			err := database.UpdateRateLimiter(db, payload.ID, payload.MessageLimit, payload.WindowSeconds)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to update rate limiter: %v", err), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Rate limit updated"})

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
