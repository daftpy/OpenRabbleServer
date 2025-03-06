package server

import (
	"chatserver/internal/client"
	"chatserver/internal/db"
	"chatserver/internal/hub"
	"chatserver/internal/messages"
	"chatserver/internal/models"
	"context"
	"fmt"
	"log"
	"net/http"

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
	var userSub string
	var clientID string
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if u, ok := claims["preferred_username"].(string); ok {
			username = u
		} else {
			log.Println("No username found in token claims")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// sub (User ID)
		if s, ok := claims["sub"].(string); ok {
			userSub = s
			log.Println("Sub found", userSub)
		} else {
			log.Println("No sub found in token claims.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if c, ok := claims["azp"].(string); ok {
			clientID = c
			log.Printf("%s connected through %s", username, clientID)
		} else {
			log.Println("No keycloak client found.")
			http.Error(w, "Bad Request", http.StatusBadRequest)
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
		Send:     make(chan messages.BaseMessage, 256),
		Hub:      s.hub,
		Sub:      userSub,
		ClientID: clientID,
	}

	// Notify the other clients that a new client has connected (if they did not connect through dashboard)
	if clientID != "WebClient" {
		newConnectionMessage := messages.NewUserStatusMessage(client.Username, true)
		s.hub.SendMessage(newConnectionMessage)
	}

	// Register Client with the Hub
	s.hub.RegisterClient(client, client.ClientID)

	// Send Connected Users List to the new client
	connectedMsg := messages.NewConnectedUsersMessage(s.hub.GetConnectedUsers())
	client.SendMessage(connectedMsg)

	// If webclient (admin dash), send analytics
	if clientID == "WebClient" {
		// Send the channel message count analytics
		counts, err := db.FetchMessageCountByChannel(s.db)
		if err != nil {
			log.Printf("Failed to get channel message counts")
		}
		analyticsMsg := messages.NewMessageCountByChannelMessage(counts)
		client.SendMessage(analyticsMsg)

		// Send the activity analytics
		activity, err := db.FetchSessionActivity(s.db)
		if err != nil {
			log.Printf("Failed to get recent activity: %v", err)
		}
		activityMsg := messages.NewSessionActivityMessage(activity)
		client.SendMessage(activityMsg)
	}

	// Bulk send chat history to the new client
	cachedMessages := s.hub.GetCachedChatMessages()
	if len(cachedMessages) > 0 {
		bulkMessage := messages.NewBulkChatMessages(cachedMessages)
		if err := conn.WriteJSON(bulkMessage); err != nil {
			log.Printf("Failed to send bulk chat messages: %v", err)
		} else {
			log.Printf("Sent %d cached messages to client", len(cachedMessages))
		}
	}

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

	// Register handlers
	mux.HandleFunc("/discovery", HandleDiscovery())
	mux.HandleFunc("/channels", HandleChannels(db, srv))
	mux.HandleFunc("/messages", HandleMessages(db))
	mux.HandleFunc("/users", HandleUsers(db))

	return srv, nil
}

func (s *Server) getChannels() ([]models.Channel, error) {
	// Update query to select both name and description
	rows, err := s.db.Query(context.Background(), "SELECT name, description FROM chatserver.channels")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channels: %w", err)
	}
	defer rows.Close()

	var channels []models.Channel
	for rows.Next() {
		var channel models.Channel
		if err := rows.Scan(&channel.Name, &channel.Description); err != nil {
			return nil, fmt.Errorf("failed to scan channel row: %w", err)
		}
		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	log.Printf("Loaded %d channels from database", len(channels))
	return channels, nil
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
