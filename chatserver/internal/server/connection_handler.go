package server

import (
	"chatserver/internal/client"
	"chatserver/internal/db"
	"chatserver/internal/messages"
	"chatserver/internal/messages/chat"
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

// upgrader defines the WebSocket upgrader that allows connections from any origin.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (adjust as needed)
	},
}

// handleConnection upgrades an HTTP request to a WebSocket connection,
// validates the JWT token, registers the user with the hub, and starts
// the read/write pumps for real-time communication.
func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	// Extract JWT Token
	token := r.URL.Query().Get("token")
	if len(token) == 0 {
		log.Println("No token provided, rejecting connection.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, userSub, clientID, err := s.parseAndValidateJWT(token)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("%s connecting through %s", username, clientID)

	banned, err := db.IsUserBanned(s.db, userSub)
	if err != nil {
		http.Error(w, "Could not determine ban status for user.", http.StatusInternalServerError)
		return
	}
	if banned {
		log.Printf("User %s is banned", userSub)
		http.Error(w, "User is banned.", http.StatusUnauthorized)
		return
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
		newConnectionMessage := chat.NewUserStatusMessage(client.Username, client.Sub, true)
		s.hub.SendMessage(newConnectionMessage)
	}

	// Register Client with the Hub
	s.hub.RegisterClient(client, client.ClientID)

	// Send the active channels and cached server messages to the client
	if err := s.sendChannelsAndCachedMessages(conn); err != nil {
		http.Error(w, "Failed to initialize chat data.", http.StatusInternalServerError)
		return
	}

	// Send Connected Users List to the new client
	connectedMsg := chat.NewConnectedUsersMessage(s.hub.GetConnectedUsers())
	client.SendMessage(connectedMsg)

	// Start Read/Write Pumps
	log.Println("Starting read/write pumps")
	go client.ReadPump()
	go client.WritePump()
}

// parseAndValidateJWT parses and validates the JWT token and extracts the username, useerID (sub), and clientID.
func (s *Server) parseAndValidateJWT(token string) (string, string, string, error) {
	parsedToken, err := jwt.Parse(token, s.jwkKeyFunc)
	if err != nil || parsedToken == nil || !parsedToken.Valid {
		return "", "", "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", "", errors.New("invalid token claims")
	}

	username, ok := claims["preferred_username"].(string)
	if !ok {
		return "", "", "", errors.New("username missing from token")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", "", "", errors.New("sub missing from token")
	}

	clientID, ok := claims["azp"].(string)
	if !ok {
		return "", "", "", errors.New("client ID missing from token")
	}

	return username, sub, clientID, nil
}

// sendChannelsAndCachedMessages sends cached chat messages and active channel list to the connected client.
func (s *Server) sendChannelsAndCachedMessages(conn *websocket.Conn) error {
	// Send cached chat messages
	cachedMessages := s.hub.GetCachedChatMessages()
	if len(cachedMessages) > 0 {
		bulkMessage := chat.NewBulkChatMessages(cachedMessages)
		if err := conn.WriteJSON(bulkMessage); err != nil {
			log.Printf("Failed to send bulk chat messages: %v", err)
		} else {
			log.Printf("Sent %d cached messages to client", len(cachedMessages))
		}
	}

	// Fetch and send channels
	channels, err := db.FetchChannels(s.db)
	if err != nil {
		log.Println("Failed to load channels from database:", err)
		return err
	}

	newActiveChannnelsMessage := chat.NewActiveChannelsMessage(channels)
	if err := conn.WriteJSON(newActiveChannnelsMessage); err != nil {
		log.Printf("Failed to send channels to client: %v", err)
		return err
	}

	return nil
}
