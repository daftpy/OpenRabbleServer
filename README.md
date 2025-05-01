# OnRabble

**OnRabble** is a modern, self-hostable real-time chat platform designed for security, scalability, and observability. It combines WebSocket messaging, Keycloak-based authentication, PostgreSQL for persistence, Valkey for caching, and a React-based admin dashboard — all orchestrated behind a Caddy reverse proxy.


## 🚀 Features

- **Secure Authentication** – Keycloak provides robust user management and JWT-based authentication.
- **Real-Time Messaging** – Go-based WebSocket server for fast, reliable chat delivery.
- **Scalable Storage** – PostgreSQL stores messages, users, and analytics data.
- **High-Speed Caching** – Valkey (a Redis-compatible engine) caches recent messages and enforces rate limits.
- **Reverse Proxy (HTTPS)** – Caddy terminates TLS and securely routes traffic.
- **Admin Dashboard** – A web-based React app for user moderation and analytics.


## 🧩 Services Overview

### `postgres`
- Relational database backing for users, messages, channels, and sessions.
- Configured via environment variables from `./postgres/.env.[mode]`.
- Health checks used for startup synchronization.

### `keycloak`
- Manages authentication, user roles, and client credentials.
- Auto-imports realm config from `./keycloak/chat-realm.[mode].json`.

### `chatserver`
- Go-based WebSocket backend, located in `./chatserver`.
- Communicates with PostgreSQL and Valkey.
- Handles routing, caching, rate limiting, and moderation.

### `valkey`
- Redis-compatible caching system.
- Supports message buffering, rate limiting, and recent history.
- See [`chatserver/internal/cache/README.md`](chatserver/internal/cache/README.md) for implementation details.

### `caddy`
- Handles HTTPS and request routing.
- Uses automatic TLS via Let's Encrypt in production.
- Configuration in `./caddy/Caddyfile.dev` and `Caddyfile.prod`.

### `web`
- React-based admin dashboard.
- Communicates with the `chatserver` over WebSockets and REST.
- Provides analytics, message history, server settings, and user moderation tools.


## 🔧 Development Setup

1. Trust Caddy’s internal root certificate if running locally.
2. Run: `docker compose -f docker-compose.dev.yml up -d --build`
3. Open the dashboard at: [https://localhost](https://localhost)

> ⚠️ Ensure your browser trusts Caddy's local CA or HTTPS will be blocked.


## Future Improvements

- [ ] Refactor message types for better modularity
- [x] Add pardon/unban functionality for admins
- [ ] Include unsaved cached messages in search results
- [ ] Allow configuration of cache size from the admin dashboard
- [x] Enable updating and deleting of chat channels
- [x] Expand developer documentation
- [ ] Add a full “Getting Started” guide


## Licensing

This project is licensed under the MIT License ([LICENSE](./LICENSE)).

The services defined in `docker-compose.dev.yml` and `docker-compose.prod.yml` use external software (e.g., Keycloak, PostgreSQL, Valkey, etc.), documented in [`THIRD_PARTY_SERVICES.md`](./THIRD_PARTY_SERVICES.md).

Each subproject may have additional third-party dependencies:

- `chatserver/`: see [`chatserver/THIRD_PARTY_LICENSES.md`](./chatserver/THIRD_PARTY_LICENSES.md)
- `web/`: see [`web/THIRD_PARTY_LICENSES.md`](./web/THIRD_PARTY_LICENSES.md)
