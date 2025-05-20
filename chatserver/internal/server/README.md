# Server Package

The `server` package handles incoming HTTP and WebSocket connections, authenticates users with JWTs, upgrades validated requests to persistent WebSocket sessions, and registers connected clients with the central hub.


## Architecture Overview

### Responsibilities

- Serve the public WebSocket entrypoint at `/ws`
- Authenticate clients using Keycloak-issued JWTs
- Reject unauthorized or banned users
- Upgrade eligible HTTP requests to WebSocket connections
- Register new clients with the `hub` instance
- Launch per-client `ReadPump` and `WritePump` goroutines
- Handle admin and discovery API routes


### Main Struct

- **`Server`**:
  - `HttpServer`: Embedded `http.Server`
  - `jwkKeyFunc`: JWKS-based JWT validator
  - `hub`: The central hub instance
  - `db`: PostgreSQL database pool
  - `MessageCache`: Shared Valkey-backed message cache


## 🔐 Authentication Flow

1. Client connects to `/ws?token=<JWT>`
2. The server:
   - Parses the JWT using `jwkKeyFunc`
   - Extracts required claims: `preferred_username`, `sub`, and `azp`
   - Validates the token and checks ban status from DB
3. On success:
   - Upgrades the HTTP request to WebSocket
   - Registers the client with the hub
   - Optionally sends analytics and cached messages based on client type


## Connection Lifecycle

- **Client Connects**:
  - Validates the JWT
  - Constructs a `Client` instance
  - Sends welcome data (connected users, cached messages)
  - Starts `ReadPump` and `WritePump` goroutines on the client

- **Goroutines**:
  - `ReadPump`: Listens for inbound messages and forwards them to the hub
  - `WritePump`: Delivers outbound messages from the hub via a channel


## API Routes

| Route                 | Purpose                                   |
|----------------------|-------------------------------------------|
| `/ws`                | WebSocket entrypoint                      |
| `/discovery`         | Public discovery info                     |
| `/channels`          | Channel metadata                         |
| `/messages`          | Chat message operations                   |
| `/users`             | User metadata                            |
| `/users/ban`         | Issue user bans                          |
| `/users/bans`        | Retrieve ban history                     |
| `/activity/sessions` | View user session analytics              |
| `/activity/channels` | View message frequency by channel        |
| `/ratelimits`        | View/update message rate limiter settings|


## Initialization

```go
New(addr string, hub HubInterface, db *pgxpool.Pool, cache *MessageCache)
```
- Sets up HTTP mux and CORS
- Loads JWKS from Keycloak for JWT validation
- Applies rate limiter values from the database
- Registers REST and WebSocket endpoints


## 🛡️ HTTPS & Reverse Proxy Configuration

This server is designed to run **behind a reverse proxy**, which handles HTTPS termination.


### Deployment with Caddy

Caddy is used for both local development and production. The repo already includes:
- `caddyfile.dev` for internal CA on `localhost`
- `caddyfile.prod` for public HTTPS using Let’s Encrypt

> Caddy automatically enables HTTPS — no `tls` block is needed unless customizing behavior.

> ⚠️ In development, you must trust Caddy's local certificate authority to allow HTTPS to work properly on `localhost`.

### WebSocket Security

Ensure WebSocket clients use `wss://` when connecting through Caddy:

```caddyfile
chat.example.com {
    reverse_proxy chatserver:8080 {
        header_up Host {host}
        header_up X-Forwarded-Proto {scheme}
    }
    log {
        output stdout
        level debug
    }
    # TLS is automatic in Caddy unless explicitly disabled
}
```


> ⚠️ Failing to use `wss://` in production will expose tokens over plaintext HTTP.


## 📝 TODO

- [ ] **Graceful shutdown support**  
      On termination, user session data to the database and close WebSocket connections cleanly.

- [ ] **Throttle dashboard analytics**  
      When `WebClient` connects, batch or delay analytics messages to avoid overwhelming low-powered clients.