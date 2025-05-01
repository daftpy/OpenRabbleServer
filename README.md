# OnRabble

OnRabble is a modern chat server designed to be secure, scalable, and efficient. It integrates authentication via Keycloak, a WebSocket-based chat server, a admin dashboard, and a reverse proxy for secure access.

## Features

- **Secure Authentication**: Keycloak handles user authentication and authorization.
- **WebSocket Chat Server**: Built in Go, allowing real-time communication.
- **Scalable Database**: PostgreSQL stores user and chat data.
- **Caching Layer**: Valkey (a Redis alternative) provides fast caching.
- **Reverse Proxy**: Caddy serves as a secure reverse proxy.
- **Admin Dashboard**: A web-based client for seamless server administration.

## Services

### `postgres`
- Stores user and chat data.
- Uses environment variables from `./postgres/.env.[mode]`.
- Health checks ensure it's available before dependent services start.

### `keycloak`
- Manages authentication and user management.
- Imports realm configuration from `./keycloak/chat-realm.[mode].json`.

### `chatserver`
- The WebSocket-based chat backend.
- Built from the `./chatserver` directory.
- Relies on PostgreSQL and Valkey.

### `valkey`
- Acts as a caching layer for the chat server.
- Reduces database load for real-time operations.
- See [`internal/cache/README.md`](./internal/cache/README.md) for implementation details.

### `caddy`
- Reverse proxy that handles HTTPS and routing.
- Configuration stored in `./caddy/Caddyfile.[mode]`.

### `web`
- Admin dashboard application built with React.
- Currently runs in development mode and connects to the chat server.

## Future Improvements

- Continue refactor of the chat server message types.
- Enable admins to pardon bans.
- Add messages in the cache, but not yet in the database to search results.
- Enable configuration of the message cache size from the admin dashbaord.
- Enable updating and deletion of Channels.
- Improve documentation.
- Add a getting started section.

## License

TBD
