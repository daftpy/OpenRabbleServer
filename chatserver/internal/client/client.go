package client

import (
	"chatserver/internal/hub"
	"chatserver/internal/messages"
	"log"

	"github.com/gorilla/websocket"
)

/*
Represents a single websocket connection from a user.
Manages receiving and sending messages to/from the server.
*/
type Client struct {
	Username string
	Conn     *websocket.Conn
	Send     chan messages.Messager
	Hub      hub.HubInterface
}

// Returns the clients username.
func (c *Client) GetUsername() string {
	return c.Username
}

// Places messages into the send channel to be picked up by WritePump().
func (c *Client) SendMessage(msg messages.Messager) {
	c.Send <- msg
}

/*
Listens for incoming messages from the websocket and processes the
messages which are then sent to the hub for broadcast.
*/
func (c *Client) ReadPump() {
	// Close the connection when the function exits
	defer func() {
		c.Hub.UnregisterClient(c)
		c.Conn.Close()
	}()

	for {
		// Read the next message from the websocket
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read error for %s: %v", c.Username, err)
		}

		// Use the raw text to create a new ChatMessage
		msg := messages.NewChatMessage(string(p), c.Username, "default")

		log.Printf("Message received from %s", c.Username)

		// Send the chat message to the hub for broadcast
		c.Hub.SendMessage(msg)
	}
}

/*
Listens for messages in the send channel and writes them to the
websocket. It ensures ouutgoing messages are sent asynchronously.
*/
func (c *Client) WritePump() {
	// Close the connection when the function exits
	defer func() {
		c.Hub.UnregisterClient(c)
		c.Conn.Close()
	}()

	for msg := range c.Send {
		// Coonvert the message to JSON and send it
		err := c.Conn.WriteJSON(msg)

		if err != nil {
			log.Printf("Write error for %s: %v", c.Username, err)
			break
		}
		log.Printf("Message sent for %s: %v", c.Username, msg.MessageType())
	}
}
