package server

import (
	"chatserver/internal/client"
	"chatserver/internal/hub"
	"chatserver/internal/messages"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

/*
ConnectedUsers manages a list of currently connected users.
It provides thread-safe operations for adding and checking connections
*/
type ConnectedUsers struct {
	users map[string]struct{}
	mu    sync.Mutex
}

/*
Server handles incoming websocket connections and authenticates them.
It validates JWT tokens, registers clients with the hub, and manages connections.
*/
type Server struct {
	HttpServer     *http.Server
	jwkKeyFunc     jwt.Keyfunc
	connectedUsers ConnectedUsers
	hub            hub.HubInterface
}

// Upgrades HTTP requests to websocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (adjust as needed)
	},
}

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
	s.addConnectedUser(username)

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

	s.hub.RegisterClient(client) // Register Client with the Hub

	// Send Connected Users List
	connectedMsg := messages.NewConnectedUsersMessage(s.getConnectedUsers())
	client.SendMessage(connectedMsg)

	// Start Read/Write Pumps
	log.Println("Starting read/write pumps")
	go client.ReadPump()
	go client.WritePump()
}

/*
New initializes and returns a new Server instance.
It sets up the WebSocket handler and configures JWT authentication.
*/
func New(addr string, h hub.HubInterface) (*Server, error) {
	mux := http.NewServeMux()

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

	// Load JWKS URL for JWT verification
	jwksURL := "http://keycloak:8080/realms/Chatserver/protocol/openid-connect/certs"
	k, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		log.Printf("failed to create JWK Keyfunc: %v", err)
	}

	// Initialize the server with HubInterface
	srv := &Server{
		HttpServer: &http.Server{Addr: addr, Handler: mux},
		jwkKeyFunc: k.Keyfunc,
		connectedUsers: ConnectedUsers{
			users: make(map[string]struct{}),
		},
		hub: h, // ✅ Use provided HubInterface implementation
	}

	// WebSocket route handled by Server struct
	mux.HandleFunc("/ws", srv.handleConnection)

	return srv, nil
}

// Adds a new user to connectedUsers
func (s *Server) addConnectedUser(username string) {
	s.connectedUsers.mu.Lock()
	defer s.connectedUsers.mu.Unlock()
	s.connectedUsers.users[username] = struct{}{}
	log.Printf("User added: %s", username)
}

// Returns if a particular user is connected
func (s *Server) isConnected(username string) bool {
	s.connectedUsers.mu.Lock()
	defer s.connectedUsers.mu.Unlock()
	_, exists := s.connectedUsers.users[username]
	return exists
}

// Returns a list of connected users
func (s *Server) getConnectedUsers() []string {
	s.connectedUsers.mu.Lock()
	defer s.connectedUsers.mu.Unlock()

	usernames := make([]string, 0, len(s.connectedUsers.users))
	for username := range s.connectedUsers.users {
		usernames = append(usernames, username)
	}
	return usernames
}
