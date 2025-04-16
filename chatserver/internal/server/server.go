package server

import (
	"chatserver/internal/cache"
	"chatserver/internal/client"
	"chatserver/internal/db"
	database "chatserver/internal/db"
	"chatserver/internal/hub"
	"chatserver/internal/messages"
	"chatserver/internal/messages/api"
	"chatserver/internal/messages/chat"
	"chatserver/internal/server/handlers"
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
	HttpServer   *http.Server
	jwkKeyFunc   jwt.Keyfunc
	hub          hub.HubInterface
	db           *pgxpool.Pool
	MessageCache *cache.MessageCache
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

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins, or specify allowed origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, DELETE")
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
func New(addr string, h hub.HubInterface, db *pgxpool.Pool, cache *cache.MessageCache) (*Server, error) {
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

	// Set the rate limiter
	rl, err := database.GetRateLimiterByID(db, 1)
	if err != nil {
		return nil, err
	}
	cache.UpdateRateLimitSettings(rl.MessageLimit, rl.WindowSeconds)

	serverIdentity := database.RegisterOrLoadServer(db)

	// Register handlers
	RegisterRoutes(srv, mux, db, cache, serverIdentity)

	return srv, nil
}

func RegisterRoutes(srv *Server, mux *http.ServeMux, db *pgxpool.Pool, cache *cache.MessageCache, identity database.ServerIdentity) {
	// WebSocket route handled by Server struct
	mux.HandleFunc("/ws", srv.handleConnection)

	// Register the other endpoints
	mux.HandleFunc("/discovery", handlers.HandleDiscovery(identity))
	mux.HandleFunc("/channels", handlers.HandleChannels(db))
	mux.HandleFunc("/messages", handlers.HandleMessages(db, cache))
	mux.HandleFunc("/users", handlers.HandleUsers(db))
	mux.HandleFunc("/users/ban", handlers.HandleBanUser(db))
	mux.HandleFunc("/users/bans", handlers.HandleBanRecords(db))
	mux.HandleFunc("/activity/sessions", handlers.HandleRecentActivity(db))
	mux.HandleFunc("/activity/channels", handlers.HandleChannelActivity(db))
	mux.HandleFunc("/ratelimits", handlers.HandleRateLimiter(db, cache))
}
