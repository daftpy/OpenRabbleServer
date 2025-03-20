package messages

import (
	"chatserver/internal/models"
	"log"
	"time"
)

const (
	ChatMessageType         = "chat_message"
	BulkChatMessagesType    = "bulk_chat_messages"
	MessageSearchResultType = "message_search_result"
)

func NewChatMessage(ID, username, channel, message string, authoredAt time.Time) BaseMessage {
	log.Printf("DEBUG STEP LOOK %s, %s, %s", username, channel, message)
	return BaseMessage{
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

func NewBulkChatMessages(messages []models.ChatMessage) BaseMessage {
	return BaseMessage{
		Type:   BulkChatMessagesType,
		Sender: "Server",
		Payload: BulkChatMessagesPayload{
			Messages: messages,
		},
	}
}

type MessageSearchResultPayload struct {
	Messages []models.ChatMessage `json:"messages"`
	HasMore  bool                 `json:"has_more"`
}

func NewMessageSearchResultMessage(payload MessageSearchResultPayload) BaseMessage {
	return BaseMessage{
		Type:    MessageSearchResultType,
		Sender:  "server",
		Payload: payload, // Should be a single struct, not a slice
	}
}
