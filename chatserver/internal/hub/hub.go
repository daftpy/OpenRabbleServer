package hub

import (
	"chatserver/internal/messages"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/valkey-io/valkey-go"
)

/*
Hub manages all active client connections and routes messages.
It recceives incoming messages, handles client registration/unregistration,
and broadcasts messages to clients.
*/
type Hub struct {
	Connections     map[string]ClientInterface
	Messages        chan messages.Messager
	Register        chan ClientInterface
	Unregister      chan ClientInterface
	MessageLogCount int
	cachClient      valkey.Client
}

// Creates a new Hub instance
func NewHub() *Hub {
	// Initialize Valkey client
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{"valkey:6379"},
	})
	if err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}
	return &Hub{
		Connections:     make(map[string]ClientInterface),
		Messages:        make(chan messages.Messager),
		Register:        make(chan ClientInterface),
		Unregister:      make(chan ClientInterface),
		MessageLogCount: 0,
		cachClient:      client,
	}
}

// Registers a new client with the hub, allowing them to send and receive messages.
func (h *Hub) RegisterClient(client ClientInterface, clientID string) {
	key := fmt.Sprintf("%s:%s", client.GetID(), clientID)
	log.Println("Hub Registered:", key)
	h.Connections[key] = client
	log.Printf("User registered: %s", client.GetUsername())
}

// Unregisters a client from the hub, stopping them from sending and receiving messages.
func (h *Hub) UnregisterClient(client ClientInterface, clientID string) {
	key := fmt.Sprintf("%s:%s", client.GetID(), clientID)
	if _, ok := h.Connections[key]; ok {
		delete(h.Connections, key)

		// Safely close the channel only if it's not already closed
		closeClientSendChannel(client)

		// Broadcast the disconnected message
		msg := messages.NewUserStatusMessage(client.GetUsername(), false)
		h.Broadcast(msg)

		log.Printf("User unregistered: %s", client.GetUsername())
	}
}

// Helper function to safely close the channel
func closeClientSendChannel(client ClientInterface) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic when closing channel: %v", r)
		}
	}()
	client.CloseSendChannel()
}

// Sends a message to the hub Messages channel for processing and broadcasting.
func (h *Hub) SendMessage(msg messages.Messager) {
	h.Messages <- msg
}

// Returns a list of the currently connected users.
func (h *Hub) GetConnectedUsers() []string {
	var users []string
	for _, v := range h.Connections {
		if v.GetClientID() != "WebClient" {
			users = append(users, v.GetUsername())
		}
	}
	return users
}

/*
Processes incoming messages based on their type.
It currently supports chat messages and user connection updates.
*/
func (h *Hub) handleMessage(msg messages.Messager) {
	switch msg := msg.(type) {
	case messages.ChatMessage:
		log.Printf("Handling chat message from: %s", msg.Username)
		h.cacheChatMessage(msg)
		h.Broadcast(msg)

	case messages.UserStatusMessage:
		log.Printf("Handling user status message for: %s - %v", msg.Username, msg.IsConnected)
		h.Broadcast(msg)

	case messages.ConnectedUsersMessage:
		log.Println("Sending connected users list")
		h.Broadcast(msg)

	default:
		log.Printf("Unhandled message type: %s", msg.MessageType())
	}
}

// Caches a chat message in Valkey
func (h *Hub) cacheChatMessage(msg messages.ChatMessage) {
	// Serialize the chat message to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to serialize chat message: %v", err)
		return
	}

	// Create a cache key using the channel and timestamp
	cacheKey := "messages"

	// Initialize the context
	ctx := context.Background()

	// Push the message to a list in Valkey
	err = h.cachClient.Do(
		ctx,
		h.cachClient.B().Rpush().Key(cacheKey).Element(string(jsonData)).Build(),
	).Error()

	if err != nil {
		log.Printf("Failed to cache chat message in Valkey: %v", err)
	} else {
		log.Printf("Cached chat message in Valkey: %s", msg.Message)
	}

	// Optional: Set a time-based expiration. Messages older than 24 hours are removed
	ttl := 24 * time.Hour
	err = h.cachClient.Do(
		ctx,
		h.cachClient.B().Expire().Key(cacheKey).Seconds(int64(ttl.Seconds())).Build(),
	).Error()

	if err != nil {
		log.Printf("Failed to set expiration for chat message cache: %v", err)
	}
}

// Retrieves chat messages from Valkey and returns them as a slice of ChatMessage
func (h *Hub) GetCachedChatMessages() []messages.ChatMessage {
	cacheKey := "messages" // Using a general key for all messages for now
	ctx := context.Background()

	// Retrieve all cached messages from Valkey
	cachedMessages, err := h.cachClient.Do(
		ctx,
		h.cachClient.B().Lrange().Key(cacheKey).Start(0).Stop(-1).Build(),
	).AsStrSlice()

	if err != nil {
		log.Printf("Failed to retrieve cached messages from Valkey: %v", err)
		return nil
	}

	var chatMessages []messages.ChatMessage

	// Deserialize each cached message back into a ChatMessage object
	for _, jsonData := range cachedMessages {
		var msg messages.ChatMessage
		if err := json.Unmarshal([]byte(jsonData), &msg); err != nil {
			log.Printf("Failed to deserialize chat message: %v", err)
			continue
		}
		chatMessages = append(chatMessages, msg)
	}

	log.Printf("Retrieved %d messages from cache", len(chatMessages))
	return chatMessages
}

/*
Broadcasts a message to all connected clients.
Every client in the hub receives the message.
*/
func (h *Hub) Broadcast(msg messages.Messager) {
	log.Printf("Broadcasting message of type: %s", msg.MessageType())
	for _, client := range h.Connections {
		log.Printf("Sending message to: %s", client.GetUsername())
		client.SendMessage(msg)
	}
}

/*
Listens for client registration, unregistration, and incoming messages.
Runs in a separate goroutine to handle Hub events asynchronously.
*/
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.RegisterClient(client, client.GetClientID())

		case client := <-h.Unregister:
			h.UnregisterClient(client, client.GetClientID())

		case message := <-h.Messages:
			h.handleMessage(message)
		}
	}
}
