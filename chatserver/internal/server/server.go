package server

import (
	"log"
	"net/http"

	database "onrabble.com/chatserver/internal/db"

	"onrabble.com/chatserver/internal/cache"
	"onrabble.com/chatserver/internal/interfaces"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server handles WebSocket connections, authentication, and client coordination.
type Server struct {
	HttpServer   *http.Server            // The underlying HTTP server.
	jwkKeyFunc   jwt.Keyfunc             // Key function for JWT validation.
	hub          interfaces.HubInterface // Central hub for managing client communication.
	db           *pgxpool.Pool           // PostgreSQL database pool.
	MessageCache *cache.MessageCache     // Shared message cache for recent chat messages.
}

// New initializes and returns a new Server instance.
// It configures JWT authentication, rate limiting, and registers all HTTP routes.
func New(addr string, h interfaces.HubInterface, db *pgxpool.Pool, cache *cache.MessageCache) (*Server, error) {
	mux := http.NewServeMux()
	handler := enableCORS(mux)

	// Load JWKS for JWT validation
	jwksURL := "http://keycloak:8080/realms/Chatserver/protocol/openid-connect/certs"
	k, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		log.Printf("failed to create JWK Keyfunc: %v", err)
	}

	// Create server instance
	srv := &Server{
		HttpServer: &http.Server{Addr: addr, Handler: handler},
		jwkKeyFunc: k.Keyfunc,
		hub:        h,
		db:         db,
	}

	// Load rate limiting settings from DB
	rl, err := database.GetRateLimiterByID(db, 1)
	if err != nil {
		return nil, err
	}
	cache.UpdateRateLimitSettings(rl.MessageLimit, rl.WindowSeconds)

	// Register this server instance and HTTP routes
	serverIdentity := database.RegisterOrLoadServer(db)
	RegisterRoutes(srv, mux, db, cache, serverIdentity)

	return srv, nil
}
