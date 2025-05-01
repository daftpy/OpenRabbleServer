package hub

import (
	"chatserver/internal/cache"
	"chatserver/internal/db"
	"chatserver/internal/messages"
	"chatserver/internal/messages/chat"
	"chatserver/internal/models"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
Hub manages all active client connections and routes messages.
It recceives incoming messages, handles client registration/unregistration,
and broadcasts messages to clients.
*/
type Hub struct {
	Connections  map[string]ClientInterface
	Messages     chan messages.BaseMessage
	Register     chan ClientInterface
	Unregister   chan ClientInterface
	MessageCache *cache.MessageCache
	db           *pgxpool.Pool
}

// Creates a new Hub instance
func NewHub(db *pgxpool.Pool, cache *cache.MessageCache) *Hub {

	return &Hub{
		Connections:  make(map[string]ClientInterface),
		Messages:     make(chan messages.BaseMessage),
		Register:     make(chan ClientInterface),
		Unregister:   make(chan ClientInterface),
		MessageCache: cache,
		db:           db,
	}
}

// Registers a new client with the hub, allowing them to send and receive messages.
func (h *Hub) RegisterClient(client ClientInterface, clientID string) {
	key := fmt.Sprintf("%s:%s", client.GetID(), clientID)
	log.Println("Hub Registered:", key)
	h.Connections[key] = client
	client.StartConnectionTimer()
	log.Printf("User registered: %s", client.GetUsername())
}

// Unregisters a client from the hub, stopping them from sending and receiving messages.
func (h *Hub) UnregisterClient(client ClientInterface, clientID string) {
	key := fmt.Sprintf("%s:%s", client.GetID(), clientID)
	if _, ok := h.Connections[key]; ok {
		delete(h.Connections, key)

		sessionStart := client.GetConnectedAt()
		sessionEnd := time.Now()

		err := db.RecordUserSession(h.db, client.GetID(), sessionStart, sessionEnd)
		if err != nil {
			log.Printf("Failed to record session for %s: %v", client.GetUsername(), err)
		} else {
			log.Printf("Session recorded for %s (duration: %v)", client.GetUsername(), (sessionEnd.Sub(sessionStart)))
		}

		// Safely close the channel only if it's not already closed
		closeClientSendChannel(client)

		// Broadcast the disconnected message
		msg := chat.NewUserStatusMessage(client.GetUsername(), client.GetID(), false)
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
func (h *Hub) SendMessage(msg messages.BaseMessage) {
	h.Messages <- msg
}

// Returns a list of the currently connected users.
func (h *Hub) GetConnectedUsers() []chat.UserStatusPayload {
	// var users []string
	var users []chat.UserStatusPayload
	for _, v := range h.Connections {
		if v.GetClientID() != "WebClient" {
			users = append(users, chat.UserStatusPayload{
				Username:    v.GetUsername(),
				ID:          v.GetID(),
				IsConnected: true,
			})
		}
	}
	return users
}

/*
Processes incoming messages based on their type.
It currently supports chat messages and user connection updates.
*/
func (h *Hub) handleMessage(msg messages.BaseMessage) {
	switch msg.Type {
	case chat.ChatMessageType:
		log.Printf("Handling chat message from: %s", msg.Sender)
		log.Printf("DEBUG: msg.Payload actual type: %T", msg.Payload)

		// Extract chat message payload
		payload, ok := msg.Payload.(models.ChatMessage)
		if !ok {
			log.Println("invalid chat message payload")
			break
		}

		// Get cacheID from CacheChatMessage
		cacheID, err := h.MessageCache.AttemptCacheWithRateLimit(payload.OwnerID, payload)
		if err != nil {
			// The user is blocked by rate limit or something else went wrong
			log.Printf("Rate limited or error: %v", err)
			// Possibly send a “rate limit exceeded” message back to the client
			break
		}

		// Attach cacheID to msg.Payload for broadcasting
		payload.CacheID = cacheID
		msg.Payload = payload // Update BaseMessage with new payload

		log.Printf("Broadcasting message with cacheID %d", cacheID)

		// Now broadcast with cacheID included
		h.Broadcast(msg)

	case chat.UserStatusMessageType:
		log.Printf("Handling user status message for: %s - %v", msg.Sender, msg.Payload)
		h.Broadcast(msg)

	case chat.ConnectedUsersMessageType:
		log.Println("Sending connected users list")
		h.Broadcast(msg)

	case chat.PrivateChatMessageType:
		log.Println("Received a private chat message")
		h.Whisper(msg)

	default:
		log.Printf("Unhandled message type: %s", msg.Type)
	}
}

// Retrieves chat messages from the MessageCache and returns them as a slice of ChatMessage
func (h *Hub) GetCachedChatMessages() []models.ChatMessage {
	chatMessages := h.MessageCache.GetCachedChatMessages()
	for i := 0; i < len(chatMessages); i++ {
		log.Printf("Message %d: %v", i, chatMessages[i])
	}
	return chatMessages
}

/*
Broadcasts a message to all connected clients.
Every client in the hub receives the message.
*/
func (h *Hub) Broadcast(msg messages.BaseMessage) {
	log.Printf("Broadcasting message of type: %s", msg.Type)
	for _, client := range h.Connections {
		log.Printf("Sending message to: %s", client.GetUsername())
		client.SendMessage(msg)
	}
}

func (h *Hub) Whisper(msg messages.BaseMessage) {
	log.Printf("Whispering message of type: %s", msg.Type)

	// Extract private message payload
	payload, ok := msg.Payload.(models.PrivateChatMessage)
	if !ok {
		log.Println("Invalid private chat message payload")
		return
	}

	senderID := payload.OwnerID
	recipientID := payload.RecipientID

	for key, client := range h.Connections {
		clientID := client.GetID()

		if clientID == senderID || clientID == recipientID {
			log.Printf("Sending whisper to: %s (key: %s)", client.GetUsername(), key)
			client.SendMessage(msg)
		}
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

// FindUsernameByUserID searches all connections and returns the username associated with the given user ID.
func (h *Hub) FindUsernameByUserID(userID string) (string, bool) {
	for _, client := range h.Connections {
		if client.GetID() == userID {
			return client.GetUsername(), true
		}
	}
	return "", false // not found
}
