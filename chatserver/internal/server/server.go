package server

import (
	"chatserver/internal/client"
	database "chatserver/internal/db"
	"chatserver/internal/hub"
	"chatserver/internal/messages"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
Server handles incoming websocket connections and authenticates them.
It validates JWT tokens, registers clients with the hub, and manages connections.
*/
type Server struct {
	HttpServer *http.Server
	jwkKeyFunc jwt.Keyfunc
	hub        hub.HubInterface
	db         *pgxpool.Pool
}

// Upgrades HTTP requests to websocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (adjust as needed)
	},
}

var ServerChannels = []string{"general", "gaming", "tech"}

/*
Handles websocket upgrade requests.
It validates the JWT, registers the client with the hub, and starts read/write pumps.
*/
func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	// Extract JWT Token
	token := r.URL.Query().Get("token")
	if len(token) == 0 {
		log.Println("No token provided, rejecting connection.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse JWT Token
	parsedToken, err := jwt.Parse(token, s.jwkKeyFunc)
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		log.Println("Invalid token, rejecting connection:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract Username from Claims
	var username string
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if u, ok := claims["preferred_username"].(string); ok {
			username = u
		} else {
			log.Println("No username found in token claims")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	log.Println("User connected:", username)

	// Upgrade to WebSocket Connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	// Create Client and Register with Hub
	client := &client.Client{
		Username: username,
		Conn:     conn,
		Send:     make(chan messages.Messager, 256),
		Hub:      s.hub,
	}

	// Notify the other clients that a new client has connected
	newConnectionMessage := messages.NewUserStatusMessage(client.Username, true)
	s.hub.SendMessage(newConnectionMessage)

	// Register Client with the Hub
	s.hub.RegisterClient(client)

	// Send Connected Users List to the new client
	connectedMsg := messages.NewConnectedUsersMessage(s.hub.GetConnectedUsers())
	client.SendMessage(connectedMsg)

	// Fetch channels from the database
	channels, err := s.getChannels()
	if err != nil {
		log.Println("Failed to load channels from database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send active channels to the user
	newActiveChannnelsMessage := messages.NewActiveChannelsMessage(channels)
	conn.WriteJSON(newActiveChannnelsMessage)

	// Start Read/Write Pumps
	log.Println("Starting read/write pumps")
	go client.ReadPump()
	go client.WritePump()
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins, or specify allowed origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
New initializes and returns a new Server instance.
It sets up the WebSocket handler and configures JWT authentication.
*/
func New(addr string, h hub.HubInterface, db *pgxpool.Pool) (*Server, error) {
	mux := http.NewServeMux()

	// Apply CORS Middleware to Allow Cross-Origin Requests
	handler := enableCORS(mux)

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

	// Channels Handler
	mux.HandleFunc("/channels", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Handle GET: Fetch channels
			channels, err := database.FetchChannels(db)
			if err != nil {
				log.Println("Failed to load channels from database:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string][]string{"channels": channels})

		case http.MethodPost:
			// Handle POST: Create a new channel
			var request struct {
				Name string `json:"name"`
			}

			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
				return
			}

			if request.Name == "" {
				http.Error(w, "Channel name is required", http.StatusBadRequest)
				return
			}

			// Insert into the database
			err := createChannel(db, request.Name)
			if err != nil {
				log.Println("Failed to create channel:", err)
				http.Error(w, "Failed to create channel", http.StatusInternalServerError)
				return
			}

			log.Printf("Channel '%s' created successfully", request.Name)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"message": "Channel created", "name": request.Name})

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Load JWKS URL for JWT verification
	jwksURL := "http://keycloak:8080/realms/Chatserver/protocol/openid-connect/certs"
	k, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		log.Printf("failed to create JWK Keyfunc: %v", err)
	}

	// Initialize the server with HubInterface
	srv := &Server{
		HttpServer: &http.Server{Addr: addr, Handler: handler},
		jwkKeyFunc: k.Keyfunc,
		hub:        h,
		db:         db,
	}

	// WebSocket route handled by Server struct
	mux.HandleFunc("/ws", srv.handleConnection)

	return srv, nil
}

func (s *Server) getChannels() ([]string, error) {
	rows, err := s.db.Query(context.Background(), "SELECT name FROM chatserver.channels")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channels: %w", err)
	}
	defer rows.Close()

	var channels []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan channel name: %w", err)
		}
		channels = append(channels, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	log.Printf("Loaded %d channels from database", len(channels))
	return channels, nil
}

func createChannel(db *pgxpool.Pool, name string) error {
	placeholderOwner := "00000000-0000-0000-0000-000000000000"
	query := `
        INSERT INTO chatserver.channels (name, owner_id)
        VALUES ($1, $2)
        ON CONFLICT (name) DO NOTHING
    `
	_, err := db.Exec(context.Background(), query, name, placeholderOwner)
	return err
}
