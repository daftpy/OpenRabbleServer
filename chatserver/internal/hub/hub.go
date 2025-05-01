package hub

import (
	"chatserver/internal/cache"
	"chatserver/internal/db"
	"chatserver/internal/interfaces"
	"chatserver/internal/messages"
	"chatserver/internal/messages/chat"
	"chatserver/internal/models"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Hub manages all active client connections, routes messages,
// and handles broadcasting, registration, and unregistration.
type Hub struct {
	Connections  map[string]interfaces.ClientInterface
	Messages     chan messages.BaseMessage
	Register     chan interfaces.ClientInterface
	Unregister   chan interfaces.ClientInterface
	MessageCache *cache.MessageCache
	db           *pgxpool.Pool
}

// NewHub creates and returns a new Hub instance.
func NewHub(db *pgxpool.Pool, cache *cache.MessageCache) *Hub {
	return &Hub{
		Connections:  make(map[string]interfaces.ClientInterface),
		Messages:     make(chan messages.BaseMessage),
		Register:     make(chan interfaces.ClientInterface),
		Unregister:   make(chan interfaces.ClientInterface),
		MessageCache: cache,
		db:           db,
	}
}

// RegisterClient adds a client to the hub and tracks its connection start time.
func (h *Hub) RegisterClient(client interfaces.ClientInterface, clientID string) {
	key := fmt.Sprintf("%s:%s", client.GetID(), clientID)
	log.Println("Hub Registered:", key)
	h.Connections[key] = client
	client.StartConnectionTimer()
	log.Printf("User registered: %s", client.GetUsername())

	privateMessages := h.MessageCache.GetCachedPrivateMessages(client.GetID())
	if len(privateMessages) > 0 {
		bulk := chat.NewBulkPrivateMessages(privateMessages)
		client.SendMessage(bulk)
		log.Printf("Sent %d cached private messages to %s", len(privateMessages), client.GetUsername())
	}
}

// UnregisterClient removes a client from the hub and logs the session duration.
func (h *Hub) UnregisterClient(client interfaces.ClientInterface, clientID string) {
	key := fmt.Sprintf("%s:%s", client.GetID(), clientID)
	if _, ok := h.Connections[key]; ok {
		delete(h.Connections, key)

		sessionStart := client.GetConnectedAt()
		sessionEnd := time.Now()

		err := db.RecordUserSession(h.db, client.GetID(), sessionStart, sessionEnd)
		if err != nil {
			log.Printf("Failed to record session for %s: %v", client.GetUsername(), err)
		} else {
			log.Printf("Session recorded for %s (duration: %v)", client.GetUsername(), sessionEnd.Sub(sessionStart))
		}

		// Safely close the channel only if it's not already closed
		closeClientSendChannel(client)

		// Broadcast the disconnected message
		msg := chat.NewUserStatusMessage(client.GetUsername(), client.GetID(), false)
		h.Broadcast(msg)

		log.Printf("User unregistered: %s", client.GetUsername())
	}
}

// closeClientSendChannel safely closes a client’s send channel, recovering from any panic.
func closeClientSendChannel(client interfaces.ClientInterface) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic when closing channel: %v", r)
		}
	}()
	client.CloseSendChannel()
}

// SendMessage sends a message into the hub’s internal message loop for handling.
func (h *Hub) SendMessage(msg messages.BaseMessage) {
	h.Messages <- msg
}

// GetConnectedUsers returns a list of currently connected user payloads,
// excluding clients identified as "WebClient".
func (h *Hub) GetConnectedUsers() []chat.UserStatusPayload {
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

// handleMessage processes an incoming message based on its type,
// including public chat, private chat, and user status updates.
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
		h.Broadcast(msg)

	case chat.UserStatusMessageType:
		log.Printf("Handling user status message for: %s - %v", msg.Sender, msg.Payload)
		h.Broadcast(msg)

	case chat.ConnectedUsersMessageType:
		log.Println("Sending connected users list")
		h.Broadcast(msg)

	case chat.PrivateChatMessageType:
		// log.Println("Received a private chat message")
		// h.Whisper(msg)
		log.Println("Handling private chat message")

		payload, ok := msg.Payload.(models.PrivateChatMessage)
		if !ok {
			log.Println("invalid private chat message payload")
			break
		}

		cacheID, err := h.MessageCache.AttemptCachePrivateWithRateLimit(payload.OwnerID, payload)
		if err != nil {
			log.Printf("Rate limited or error (private): %v", err)
			break
		}

		payload.CacheID = cacheID
		msg.Payload = payload

		h.Whisper(msg)

	default:
		log.Printf("Unhandled message type: %s", msg.Type)
	}
}

// GetCachedChatMessages returns a slice of chat messages from the message cache.
func (h *Hub) GetCachedChatMessages() []models.ChatMessage {
	chatMessages := h.MessageCache.GetCachedChatMessages()
	for i := 0; i < len(chatMessages); i++ {
		log.Printf("Message %d: %v", i, chatMessages[i])
	}
	return chatMessages
}

// Broadcast sends the given message to all connected clients in the hub.
func (h *Hub) Broadcast(msg messages.BaseMessage) {
	log.Printf("Broadcasting message of type: %s", msg.Type)
	for _, client := range h.Connections {
		log.Printf("Sending message to: %s", client.GetUsername())
		client.SendMessage(msg)
	}
}

// Whisper sends a private message only to the sender and recipient clients.
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

// Run starts the hub's main loop and handles registration, unregistration, and messages.
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

// FindUsernameByUserID returns the username for a given user ID, if connected.
func (h *Hub) FindUsernameByUserID(userID string) (string, bool) {
	for _, client := range h.Connections {
		if client.GetID() == userID {
			return client.GetUsername(), true
		}
	}
	return "", false
}
