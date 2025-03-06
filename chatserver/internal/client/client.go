package client

import (
	"chatserver/internal/hub"
	"chatserver/internal/messages"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

/*
Represents a single websocket connection from a user.
Manages receiving and sending messages to/from the server.
*/
type Client struct {
	Username    string
	Conn        *websocket.Conn
	Send        chan messages.BaseMessage
	Hub         hub.HubInterface
	Sub         string // Keycloak stable user ID
	ClientID    string
	ConnectedAt time.Time
}

// Returns the clients username.
func (c *Client) GetUsername() string {
	return c.Username
}

// Places messages into the send channel to be picked up by WritePump().
func (c *Client) SendMessage(msg messages.BaseMessage) {
	c.Send <- msg
}

// Closes the send channel
func (c *Client) CloseSendChannel() {
	close(c.Send)
}

func (c *Client) GetID() string {
	return c.Sub
}

func (c *Client) GetClientID() string {
	return c.ClientID
}

func (c *Client) StartConnectionTimer() {
	c.ConnectedAt = time.Now()
}

func (c *Client) GetConnectedAt() time.Time {
	return c.ConnectedAt
}

/*
Listens for incoming messages from the websocket and processes the
messages which are then sent to the hub for broadcast.
*/
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
			Type    string `json:"type"`
			Channel string `json:"channel"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(p, &receivedMessage); err != nil {
			log.Printf("Invalid message from %s: %v", c.Username, err)
			continue
		}
		log.Printf("Received type: %s", receivedMessage.Type)
		log.Printf("Received channel: %s", receivedMessage.Channel)
		log.Printf("Receiived message: %s", receivedMessage.Message)

		// Process received message
		msg := messages.NewChatMessage(c.Sub, c.Username, receivedMessage.Channel, receivedMessage.Message, time.Now())
		log.Printf("Message received from %s", c.Username)

		// Send the message to the hub
		c.Hub.SendMessage(msg)
	}
}

/*
Listens for messages in the send channel and writes them to the
websocket. It ensures ouutgoing messages are sent asynchronously.
*/
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
		log.Printf("WritePump exited for %s", c.Username)
	}()

	for msg := range c.Send {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Write error for %s: %v", c.Username, err)
			break // Exit loop if there's an error
		}
		log.Printf("Message sent for %s: %v", c.Username, msg.Type)
	}
}
