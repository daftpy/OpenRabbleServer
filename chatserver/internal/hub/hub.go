package hub

import (
	"chatserver/internal/messages"
	"log"
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
}

// Creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		Connections:     make(map[string]ClientInterface),
		Messages:        make(chan messages.Messager),
		Register:        make(chan ClientInterface),
		Unregister:      make(chan ClientInterface),
		MessageLogCount: 0,
	}
}

// Registers a new client with the hub, allowing them to send and receive messages.
func (h *Hub) RegisterClient(client ClientInterface) {
	h.Connections[client.GetUsername()] = client
	log.Printf("User registered: %s", client.GetUsername())
}

// Unregisters a client from the hub, stopping them from sending and receiving messages.
func (h *Hub) UnregisterClient(client ClientInterface) {
	if _, ok := h.Connections[client.GetUsername()]; ok {
		delete(h.Connections, client.GetUsername())
		log.Printf("User unregistered: %s", client.GetUsername())
	}
}

// Sends a message to the hub Messages channel for processing and broadcasting.
func (h *Hub) SendMessage(msg messages.Messager) {
	h.Messages <- msg
}

// Returns a list of the currently connected users.
func (h *Hub) GetConnectedUsers() []string {
	var users []string
	for k := range h.Connections {
		users = append(users, k)
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
			h.RegisterClient(client)

		case client := <-h.Unregister:
			h.UnregisterClient(client)

		case message := <-h.Messages:
			h.handleMessage(message)
		}
	}
}
