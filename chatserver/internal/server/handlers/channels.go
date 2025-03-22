package handlers

import (
	database "chatserver/internal/db"
	"chatserver/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// HandleChannels handles creating and fetching channels.
func HandleChannels(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			channels, err := database.FetchChannels(db)
			if err != nil {
				log.Println("Failed to load channels from database:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			response := map[string][]models.Channel{"channels": channels}
			json.NewEncoder(w).Encode(response)

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
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "Channel created", "name": request.Name})

		case http.MethodPatch:
			var request struct {
				ID          *int    `json:"id"`
				Name        *string `json:"name,omitempty"`
				Description *string `json:"description,omitempty"`
			}

			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
				return
			}

			if request.ID == nil {
				http.Error(w, "A Channel ID is required", http.StatusBadRequest)
				return
			}

			if request.Name == nil && request.Description == nil {
				http.Error(w, "Nothing to update", http.StatusBadRequest)
				return
			}

			if err := database.UpdateChannel(db, *request.ID, request.Name, request.Description); err != nil {
				log.Println("Failed to update channel:", err)
				http.Error(w, "Failed to update channel", http.StatusInternalServerError)
				return
			}

			log.Printf("Channel ID '%d' updated successfully", request.ID)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Channel updated"})

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func createChannel(db *pgxpool.Pool, name string, description string) error {
	placeholderOwner := "00000000-0000-0000-0000-000000000000"
	query := `
        INSERT INTO chatserver.channels (name, description, owner_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (name) DO NOTHING
    `
	_, err := db.Exec(context.Background(), query, name, description, placeholderOwner)
	return err
}
