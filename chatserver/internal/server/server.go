package server

import (
	"chatserver/internal/cache"
	database "chatserver/internal/db"
	"chatserver/internal/interfaces"
	"chatserver/internal/server/handlers"
	"log"
	"net/http"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
Server handles incoming websocket connections and authenticates them.
It validates JWT tokens, registers clients with the hub, and manages connections.
*/
type Server struct {
	HttpServer   *http.Server
	jwkKeyFunc   jwt.Keyfunc
	hub          interfaces.HubInterface
	db           *pgxpool.Pool
	MessageCache *cache.MessageCache
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
func New(addr string, h interfaces.HubInterface, db *pgxpool.Pool, cache *cache.MessageCache) (*Server, error) {
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
