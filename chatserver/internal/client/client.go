package client

import (
	"chatserver/internal/interfaces"
	"chatserver/internal/messages"
	"chatserver/internal/messages/chat"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection from a user.
// It manages receiving and sending messages to/from the server.
type Client struct {
	Username    string
	Conn        *websocket.Conn
	Send        chan messages.BaseMessage
	Hub         interfaces.HubInterface
	Sub         string // Keycloak stable user ID
	ClientID    string // OAuth client ID, e.g., "ChatClient" or "WebClient"
	ConnectedAt time.Time
}

// GetUsername returns the client's username.
func (c *Client) GetUsername() string {
	return c.Username
}

// SendMessage places a message into the send channel to be picked up by WritePump().
func (c *Client) SendMessage(msg messages.BaseMessage) {
	c.Send <- msg
}

// CloseSendChannel closes the client's outgoing message channel.
func (c *Client) CloseSendChannel() {
	close(c.Send)
}

// GetID returns the stable user ID (usually from Keycloak).
func (c *Client) GetID() string {
	return c.Sub
}

// GetClientID returns the OAuth client ID used to identify the source application.
func (c *Client) GetClientID() string {
	return c.ClientID
}

// StartConnectionTimer records the time when the client connects.
func (c *Client) StartConnectionTimer() {
	c.ConnectedAt = time.Now()
}

// GetConnectedAt returns the timestamp when the client connected.
func (c *Client) GetConnectedAt() time.Time {
	return c.ConnectedAt
}

// ReadPump listens for incoming messages from the WebSocket and processes them.
// Parsed messages are sent to the hub for broadcast or private delivery.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.UnregisterClient(c, c.ClientID)
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected disconnect from %s: %v", c.Username, err)
			} else {
				log.Printf("Client %s disconnected: %v", c.Username, err)
			}
			break
		}

		// Unmarshal the JSON message into a struct
		var receivedMessage struct {
			Type        string `json:"type"`
			Channel     string `json:"channel,omitempty"`
			RecipientID string `json:"recipient_id,omitempty"`
			Message     string `json:"message"`
		}
		if err := json.Unmarshal(p, &receivedMessage); err != nil {
			log.Printf("Invalid message from %s: %v", c.Username, err)
			continue
		}
		log.Printf("Received type: %s", receivedMessage.Type)
		log.Printf("Received channel: %s", receivedMessage.Channel)
		log.Printf("Receiived message: %s", receivedMessage.Message)

		var msg messages.BaseMessage
		// Process received message
		if receivedMessage.Type == chat.ChatMessageType {
			msg = chat.NewChatMessage(c.Sub, c.Username, receivedMessage.Channel, receivedMessage.Message, time.Now())
		} else if receivedMessage.Type == chat.PrivateChatMessageType {
			username, ok := c.Hub.FindUsernameByUserID(receivedMessage.RecipientID)
			if ok {
				msg = chat.NewPrivateChatMessage(c.Sub, c.Username, receivedMessage.RecipientID, username, receivedMessage.Message, time.Now())
			} else {
				continue
			}
		}
		log.Printf("Message received from %s", c.Username)

		// Send the message to the hub
		c.Hub.SendMessage(msg)
	}
}

// WritePump listens for messages on the send channel and writes them to the WebSocket.
// It ensures that outgoing messages are sent asynchronously.
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
		log.Printf("WritePump exited for %s", c.Username)
	}()

	for msg := range c.Send {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Write error for %s: %v", c.Username, err)
			break
		}
		log.Printf("Message sent for %s: %v", c.Username, msg.Type)
	}
}
