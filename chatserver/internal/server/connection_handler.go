package server

import (
	"chatserver/internal/client"
	"chatserver/internal/db"
	database "chatserver/internal/db"
	"chatserver/internal/messages"
	"chatserver/internal/messages/api"
	"chatserver/internal/messages/chat"
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
			log.Printf("%s connecting through %s", username, clientID)
		} else {
			log.Println("No keycloak client found.")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

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

	// Send Connected Users List to the new client
	connectedMsg := chat.NewConnectedUsersMessage(s.hub.GetConnectedUsers())
	client.SendMessage(connectedMsg)

	// If webclient (admin dash), send analytics
	if clientID == "WebClient" {
		// Send the channel message count analytics
		counts, err := database.FetchMessageCountByChannel(s.db)
		if err != nil {
			log.Printf("Failed to get channel message counts")
		}
		analyticsMsg := api.NewMessageCountByChannelMessage(counts)
		client.SendMessage(analyticsMsg)

		// Send the activity analytics
		activity, err := database.FetchSessionActivity(s.db, "")
		if err != nil {
			log.Printf("Failed to get recent activity: %v", err)
		}
		activityMsg := api.NewSessionActivityMessage(activity)
		client.SendMessage(activityMsg)
	}

	// Bulk send chat history to the new client
	cachedMessages := s.hub.GetCachedChatMessages()
	if len(cachedMessages) > 0 {
		bulkMessage := chat.NewBulkChatMessages(cachedMessages)
		if err := conn.WriteJSON(bulkMessage); err != nil {
			log.Printf("Failed to send bulk chat messages: %v", err)
		} else {
			log.Printf("Sent %d cached messages to client", len(cachedMessages))
		}
	}

	// Fetch channels from the database
	channels, err := db.FetchChannels(s.db)
	if err != nil {
		log.Println("Failed to load channels from database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send active channels to the user
	newActiveChannnelsMessage := chat.NewActiveChannelsMessage(channels)
	conn.WriteJSON(newActiveChannnelsMessage)

	// Start Read/Write Pumps
	log.Println("Starting read/write pumps")
	go client.ReadPump()
	go client.WritePump()
}
