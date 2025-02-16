package server

import (
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

type ConnectedUsers struct {
	users map[string]struct{}
	mu    sync.Mutex
}

type Server struct {
	HttpServer     *http.Server
	jwkKeyFunc     jwt.Keyfunc
	connectedUsers ConnectedUsers
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (adjust as needed)
	},
}

func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	// Check for a token
	token := r.URL.Query().Get("token")

	if len(token) == 0 {
		log.Printf("No token, no connection.")
		return
	}
	log.Println("Token:", token)
	parsedToken, err := jwt.Parse(token, s.jwkKeyFunc)
	if err != nil {
		log.Println("Parsing token failed:", err)
	}
	log.Println("Parsed token:", parsedToken)

	// Access claims
	var username string
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		username = claims["preferred_username"].(string)
	}
	log.Println("User connected: ", username)
	s.addConnectedUser(username)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket connection established")

	// connectedMsg := messages.ConnectedUsersMessage{
	// 	Type:  "status_message",
	// 	Users: s.getConnectedUsers(),
	// }
	connectedMsg := messages.NewConnectedUsersMessage(s.getConnectedUsers())
	conn.WriteJSON(connectedMsg)

	for {
		// Read incoming message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		log.Printf("Received message: %s\n", msg)

		// Create JSON response
		// readMessage := messages.ChatMessage{
		// 	Type:    "chat_message",
		// 	Message: string(msg),
		// 	User:    username, // Include sender info
		// 	Channel: "default",
		// }
		readMessage := messages.NewChatMessage(string(msg), username, "default")

		// Send the JSON message
		if err := conn.WriteJSON(readMessage); err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}
	}
}

func New(addr string) (*Server, error) {
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

	// Initialize the server
	srv := &Server{
		HttpServer: &http.Server{Addr: ":8080", Handler: mux},
		jwkKeyFunc: k.Keyfunc,
		connectedUsers: ConnectedUsers{
			users: make(map[string]struct{}),
		},
	}

	// WebSocket route handled by Server struct
	mux.HandleFunc("/ws", srv.handleConnection)

	return srv, nil
}

func (s *Server) addConnectedUser(username string) {
	s.connectedUsers.mu.Lock()
	defer s.connectedUsers.mu.Unlock()
	s.connectedUsers.users[username] = struct{}{}
	log.Printf("User added: %s", username)
}

func (s *Server) isConnected(username string) bool {
	s.connectedUsers.mu.Lock()
	defer s.connectedUsers.mu.Unlock()
	_, exists := s.connectedUsers.users[username]
	return exists
}

func (s *Server) getConnectedUsers() []string {
	s.connectedUsers.mu.Lock()
	defer s.connectedUsers.mu.Unlock()

	usernames := make([]string, 0, len(s.connectedUsers.users))
	for username := range s.connectedUsers.users {
		usernames = append(usernames, username)
	}
	return usernames
}
