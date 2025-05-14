package server

import (
	"net/http"

	"onrabble.com/chatserver/internal/cache"
	"onrabble.com/chatserver/internal/db"
	"onrabble.com/chatserver/internal/server/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterRoutes binds all HTTP endpoints, including the WebSocket route.
func RegisterRoutes(srv *Server, mux *http.ServeMux, db *pgxpool.Pool, cache *cache.MessageCache, identity db.ServerIdentity) {
	mux.HandleFunc("/ws", srv.handleConnection)
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
