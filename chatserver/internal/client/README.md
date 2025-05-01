# Client Package

The `client` package represents a single active WebSocket connection to the server. It is responsible for sending and receiving messages on behalf of a user and acts as the communication bridge between the frontend and the hub.


## Architecture

### Responsibilities

- Manage the lifecycle of an individual WebSocket connection.
- Read and deserialize incoming messages from the client.
- Send messages asynchronously via a buffered channel.
- Pass valid messages to the hub for routing and broadcasting.
- Identify the source client type using OAuth client ID (e.g., `ChatClient`, `WebClient`).
- Track connection timestamps for session analytics.


### Key Struct

- **`Client`**:
  - `Username`: user-visible name.
  - `Sub`: stable unique identifier (from OAuth provider).
  - `ClientID`: identifies the type of application used (admin dashboard, chat UI, etc.).
  - `Conn`: the active WebSocket connection.
  - `Send`: channel for outgoing messages.
  - `Hub`: reference to the central `HubInterface`.
  - `ConnectedAt`: timestamp of connection start.


## 🔁 Workflow

1. **Connection Setup**:
   - A `Client` is created with an active WebSocket connection and registered with the hub.

2. **Message Receiving (ReadPump)**:
   - Runs in a goroutine.
   - Listens for JSON messages from the WebSocket.
   - Deserializes into a lightweight struct.
   - Constructs appropriate `BaseMessage` objects based on message type.
   - Sends the message to the hub for processing.

3. **Message Sending (WritePump)**:
   - Also runs in a goroutine.
   - Listens on the `Send` channel.
   - Encodes `BaseMessage` as JSON and writes it to the WebSocket.
   - Handles cleanup on failure or disconnect.

4. **Disconnection**:
   - Triggers unregistration from the hub.
   - Closes the connection and the send channel.


## Usage Example

```go
client := &Client{
    Username: "alice",
    Sub: "user-alice-123",
    ClientID: "ChatClient",
    Conn: conn, // *websocket.Conn
    Send: make(chan messages.BaseMessage, 16),
    Hub: hub,
}

go client.ReadPump()
go client.WritePump()

hub.Register <- client
```


## 📝 TODO

- [ ] Optionally support a shutdown message or ping/pong handling to detect dead clients early.
- [ ] Add per-client rate-limiting or mute functionality at the client level.
- [ ] Log dropped messages if the `Send` channel is full or unresponsive.
