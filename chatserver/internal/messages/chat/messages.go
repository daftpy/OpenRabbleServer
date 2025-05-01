package chat

import (
	"chatserver/internal/messages"
	"chatserver/internal/models"
	"log"
	"time"
)

const (
	ChatMessageType        = "chat_message"
	BulkChatMessagesType   = "bulk_chat_messages"
	PrivateChatMessageType = "private_chat_message"
	BulkPrivateMessageType = "bulk_private_messages"
)

func NewPrivateChatMessage(ID, username, recipientID, recipient, message string, authoredAt time.Time) messages.BaseMessage {
	return messages.BaseMessage{
		Type:   PrivateChatMessageType,
		Sender: username,
		Payload: models.PrivateChatMessage{
			OwnerID:     ID,
			Username:    username,
			RecipientID: recipientID,
			Recipient:   recipient,
			Message:     message,
			Sent:        authoredAt,
		},
	}
}

func NewChatMessage(ID, username, channel, message string, authoredAt time.Time) messages.BaseMessage {
	log.Printf("DEBUG STEP LOOK %s, %s, %s", username, channel, message)
	return messages.BaseMessage{
		Type:   ChatMessageType,
		Sender: username,
		Payload: models.ChatMessage{
			Username: username,
			Channel:  channel,
			Message:  message,
			Sent:     authoredAt,
			OwnerID:  ID,
		},
	}
}

type BulkChatMessagesPayload struct {
	Messages []models.ChatMessage `json:"messages"`
}

func NewBulkChatMessages(msgs []models.ChatMessage) messages.BaseMessage {
	return messages.BaseMessage{
		Type:   BulkChatMessagesType,
		Sender: "Server",
		Payload: BulkChatMessagesPayload{
			Messages: msgs,
		},
	}
}

type BulkPrivateMessagesPayload struct {
	Messages []models.PrivateChatMessage `json:"messages"`
}

func NewBulkPrivateMessages(msgs []models.PrivateChatMessage) messages.BaseMessage {
	return messages.BaseMessage{
		Type:   BulkPrivateMessageType,
		Sender: "Server",
		Payload: BulkPrivateMessagesPayload{
			Messages: msgs,
		},
	}
}
