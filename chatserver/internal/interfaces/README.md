# Interfaces Package

The `interfaces` package defines contracts for the core components of the chat server, specifically `ClientInterface` and `HubInterface`. These interfaces promote decoupling and testability by abstracting the concrete behavior of clients and hubs.


## Purpose

This package prevents circular dependencies between components like `client`, `hub`, and `messages` by introducing a layer of abstraction. It allows:

- The hub to interact with any type that satisfies `ClientInterface`, without needing to import the `client` package.
- The client to interact with any hub implementation that satisfies `HubInterface`.


## Interfaces

### `ClientInterface`

Represents a single user session over a WebSocket.

| Method                 | Description |
|------------------------|-------------|
| `GetUsername()`        | Returns the client's display name. |
| `SendMessage(msg)`     | Sends a message to the client’s send channel. |
| `CloseSendChannel()`   | Closes the channel used to deliver outgoing messages. |
| `GetID()`              | Returns the stable user ID (e.g., from Keycloak). |
| `GetClientID()`        | Returns the OAuth client ID indicating the source app (e.g., `WebClient`, `ChatClient`). |
| `StartConnectionTimer()` | Records the connection start time for session logging. |
| `GetConnectedAt()`     | Returns the timestamp of when the client connected. |


### `HubInterface`

Defines behavior expected from the central message dispatcher and coordinator.

| Method                       | Description |
|------------------------------|-------------|
| `Broadcast(msg)`             | Sends a message to all connected clients. |
| `Whisper(msg)`               | Sends a private message between clients. |
| `RegisterClient(client, id)` | Registers a client with a unique connection ID. |
| `UnregisterClient(client, id)` | Removes a client from the hub and ends their session. |
| `SendMessage(msg)`           | Pushes a message into the hub’s processing loop. |
| `GetConnectedUsers()`        | Returns all currently connected users. |
| `GetCachedChatMessages()`    | Retrieves recent messages from the message cache. |
| `FindUsernameByUserID(id)`   | Resolves a user ID to a username, if connected. |


## Use Cases

- **The `hub` package** depends on `ClientInterface` to avoid importing the concrete `client` package.
- **The `client` package** depends on `HubInterface` to send and receive routed messages without tight coupling.


## 📝 TODO

- [ ] Add `Disconnect()` to `ClientInterface` to enable graceful shutdowns or ban logic.
- [ ] Expand `HubInterface` with scene/channel routing if multi-room support is added.
- [ ] Create mock implementations for use in unit tests.
