# OnRabble Chat Server

The `chatserver` module provides a fully self-hostable, real-time chat backend written in Go. It integrates WebSocket-based messaging, JWT authentication via Keycloak, PostgreSQL for persistence, and Valkey (Redis-compatible) for caching and rate limiting. The system is modular and designed to be observable, scalable, and easy to deploy behind a reverse proxy.


## Packages

| Package         | Description                                                                 |
|----------------|-----------------------------------------------------------------------------|
| `server`        | WebSocket server, handles JWT auth, client registration, route handling     |
| `client`        | Represents an individual WebSocket connection                               |
| `hub`           | Message router and client registry                                          |
| `cache`         | Valkey-backed caching system for chat messages and rate limiting            |
| `db`            | PostgreSQL queries and data access layer                                    |
| `handlers`      | HTTP route handlers for REST endpoints used by the admin dashboard          |
| `interfaces`    | Defines shared interfaces to reduce package coupling                        |
| `models`        | Core data structures shared between packages                                |
| `messages`      | Defines chat and API message types (including subtype folders `chat`, `api`) |


## ✅ Features

- JWT-based authentication using [Keycloak](https://www.keycloak.org/)
- WebSocket support for real-time messaging (`/ws`)
- Modular message types and routing
- Client-to-client and private messaging
- Dashboard analytics for usage and moderation
- Message caching and periodic batch database flushing
- User banning system
- Rate limiting per user
- REST API endpoints for administrative access


## 🚀 Getting Started

### Requirements

- Go 1.21+
- PostgreSQL 16+
- Valkey (or Redis)
- Keycloak (used for JWT auth)
- Caddy (reverse proxy with TLS)


### Local Dev Setup

This project is configured for local development using Docker Compose and Caddy.

1. Trust Caddy’s root CA (`caddyfile.dev` uses TLS internally)
2. Run `docker-compose up --build`
3. Open [https://chat.localhost](https://chat.localhost) in your browser

> ⚠️ Make sure you trust Caddy’s root certificate or your browser will block HTTPS requests locally.


## 🔐 Security

- All communication is authenticated with JWT tokens.
- Connections must use `wss://` in production (Caddy handles this).
- Rate limiting is enforced using Valkey scripts.
- Admin dashboard (`WebClient`) has scoped permissions via `azp` in JWT.


## Routes

- WebSocket: `/ws`
- Admin/API:
  - `/discovery`
  - `/channels`
  - `/messages`
  - `/users`, `/users/ban`, `/users/bans`
  - `/activity/sessions`, `/activity/channels`
  - `/ratelimits`


## Deployment Notes

- TLS is handled automatically by [Caddy](https://caddyserver.com)
- `caddyfile.dev` and `caddyfile.prod` are included in the project
- Server identity and rate limiting config are initialized from PostgreSQL
- JWT keys are loaded dynamically from Keycloak’s JWKS endpoint

> In development, you must trust Caddy's local CA for HTTPS to work properly on `localhost`


## 📝 TODO

- [ ] **Graceful shutdown support**  
      Flush user session durations and close WebSocket connections cleanly on shutdown.

- [ ] **Ping/pong keepalive and idle timeout**  
      Detect and clean up stale client connections automatically.

- [ ] **Throttle dashboard analytics**  
      Prevent `WebClient` connections from overwhelming the browser by staggering analytics updates.

- [ ] **Retry failed whispers**  
      Implement a retry queue for private messages if the target client is momentarily disconnected.

- [ ] **Split whisper/broadcast queues**  
      Process `Broadcast()` and `Whisper()` flows in separate background goroutines to isolate errors and improve throughput.

- [ ] **Support multi-realm JWTs**  
      Allow parsing/validation of tokens from multiple Keycloak realms or issuers.

- [ ] **Expose Prometheus metrics (optional)**  
      Emit hub/client/caching counters and durations for system health visibility.

- [ ] **Improve dashboard load UX**  
      Introduce loading states and staged delivery to avoid UI stalls.

- [ ] **Multi-source analytics expansion**  
      Add deeper analytics endpoints (e.g., per-user usage, time series)
