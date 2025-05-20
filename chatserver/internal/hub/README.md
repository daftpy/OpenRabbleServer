# Hub Package

The `hub` package is the central message coordinator in the chat server. It manages client connections, message routing, broadcasting, private messaging, and integration with both caching and persistence layers.


## Architecture

### Responsibilities

- Track all active WebSocket clients.
- Broadcast chat messages and user status updates to connected clients.
- Whisper private messages between specific users.
- Integrate with the `cache` package to temporarily store and rate-limit messages.
- Record session data to the database via the `db` package.


### Components

- **`Hub` Struct**: Core state manager with:
  - `Connections`: map of active clients.
  - `Register`, `Unregister`: channels for client lifecycle.
  - `Messages`: channel for incoming messages.
  - `MessageCache`: reference to the Valkey-backed message cache.
  - `db`: PostgreSQL pool for recording session data.

- **Message Types**: Supports:
  - Public chat messages
  - Private (whisper) messages
  - User connect/disconnect events
  - Requests for current user list


## Workflow

1. **Client Registers**:
   - On connect, a `Client` sends itself to the hub's `Register` channel.
   - The hub stores the client and begins tracking session time.

2. **Message Handling**:
   - Chat messages are received via the `Messages` channel.
   - The hub delegates by:
     - Checking message type.
     - Adding a `cacheID` (via `MessageCache`).
     - Broadcasting to all clients or sending privately.

3. **Client Disconnects**:
   - A client sends itself to the `Unregister` channel.
   - The hub removes the client, closes its channel, and writes session info to the database.


## Configuration

- The `Hub` is created via:
  ```go
  NewHub(db *pgxpool.Pool, cache *MessageCache)
  ```

- Message cache limits, flush intervals, and rate limits are configured in the `cache` package.


## Setup & Usage

1. **Instantiate the Hub**:
   ```go
   hub := hub.NewHub(dbPool, messageCache)
   ```

2. **Start the Hub Loop**:
   ```go
   go hub.Run()
   ```

3. **Clients Interact Through Channels**:
   - Register: `hub.Register <- client`
   - Unregister: `hub.Unregister <- client`
   - Send Message: `hub.Messages <- msg`


## 📝 TODO

- [ ] Graceful shutdown hook to flush sessions and broadcast disconnects.
- [ ] Implement separate broadcast and whisper queues to:
  - Prevent the hub loop from blocking if a client’s `Send` channel is full or slow
  - Decouple message delivery from message processing logic
  - Enable future enhancements like:
    - Rate limiting broadcast and whisper delivery independently
    - Filtering or moderation before delivery
    - Monitoring queue depth to detect client lag or overload
---
