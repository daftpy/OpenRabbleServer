package cache

import (
	"chatserver/internal/messages"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valkey-io/valkey-go"
)

type MessageCache struct {
	ValkeyClient valkey.Client
	DB           *pgxpool.Pool
}

// Define the Lua script
var cacheMessageScript = valkey.NewLuaScript(`
    local key = KEYS[1]
    local message = ARGV[1]
    local maxSize = tonumber(ARGV[2])
    local ttl = tonumber(ARGV[3])

    -- Append the message to the list
    redis.call("RPUSH", key, message)

    -- Maintain circular cache by trimming to maxSize
    redis.call("LTRIM", key, -maxSize, -1)

    -- Set expiration time for the cache
    redis.call("EXPIRE", key, ttl)

    return true
`)

const maxCacheSize = 100
const cacheTTL = 24 * 60 * 60 // 24 hours in seconds

func (m *MessageCache) CacheChatMessage(msg messages.ChatMessage) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to serialize chat message: %v", err)
		return
	}

	cacheKey := "messages"
	ctx := context.Background()

	// Execute the Lua script atomically
	_, err = cacheMessageScript.Exec(
		ctx,
		m.ValkeyClient,
		[]string{cacheKey}, // KEYS
		[]string{string(jsonData), // ARGV[1]: The serialized message
			fmt.Sprintf("%d", maxCacheSize), // ARGV[2]: The max cache size
			fmt.Sprintf("%d", cacheTTL)},    // ARGV[3]: The cache TTL
	).ToBool()

	if err != nil {
		log.Printf("Failed to cache chat message with Lua script: %v", err)
	} else {
		log.Printf("Cached chat message in Valkey: %s", msg.Message)
	}
}

// Retrieves chat messages from Valkey and returns them as a slice of ChatMessage
func (m *MessageCache) GetCachedChatMessages() []messages.ChatMessage {
	cacheKey := "messages"
	ctx := context.Background()

	cachedMessages, err := m.ValkeyClient.Do(
		ctx,
		m.ValkeyClient.B().Lrange().Key(cacheKey).Start(0).Stop(-1).Build(),
	).AsStrSlice()

	if err != nil {
		log.Printf("Failed to retrieve cached messages from Valkey: %v", err)
		return nil
	}

	var chatMessages []messages.ChatMessage
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
