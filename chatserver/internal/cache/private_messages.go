package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	database "onrabble.com/chatserver/internal/db"
	"onrabble.com/chatserver/internal/models"

	"github.com/valkey-io/valkey-go"
)

var cachePrivateMessageScript = valkey.NewLuaScript(`
	local senderKey = KEYS[1]
	local recipientKey = KEYS[2]
	local flushKey = KEYS[3]
	local counterKey = KEYS[4]

	local message = ARGV[1]
	local maxSize = tonumber(ARGV[2])
	local isSelf = ARGV[3] == "1"

	-- Generate unique cache ID
	local cacheID = redis.call("INCR", counterKey)

	-- Enrich the message with the ID
	local enrichedMessage = cjson.encode({cache_id = cacheID, data = cjson.decode(message)})

	-- Push to circular caches
	redis.call("RPUSH", senderKey, enrichedMessage)
	redis.call("LTRIM", senderKey, -maxSize, -1)

	if not isSelf then
		redis.call("RPUSH", recipientKey, enrichedMessage)
		redis.call("LTRIM", recipientKey, -maxSize, -1)
	end

	-- Push to flush queue
	redis.call("RPUSH", flushKey, enrichedMessage)

	-- Return cache ID and flush size
	return {cacheID, redis.call("LLEN", flushKey)}
`)

func (m *MessageCache) CachePrivateMessage(msg models.PrivateChatMessage) int {
	senderKey := fmt.Sprintf("recent_private_messages:%s", msg.OwnerID)
	recipientKey := fmt.Sprintf("recent_private_messages:%s", msg.RecipientID)
	flushKey := "flush_private_messages"
	counterKey := "cache_private_message_id"

	ctx := context.Background()

	// Serialize the private message to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to serialize private message: %v", err)
		return -1
	}

	isSelf := "0"
	if msg.OwnerID == msg.RecipientID {
		isSelf = "1"
	}

	results, err := cachePrivateMessageScript.Exec(
		ctx,
		m.ValkeyClient,
		[]string{senderKey, recipientKey, flushKey, counterKey},
		[]string{string(jsonData), fmt.Sprintf("%d", maxCacheSize), isSelf},
	).ToArray()

	if err != nil {
		log.Printf("Failed to cache private message: %v", err)
		return -1
	}

	// Extract cacheID
	cacheID, err := results[0].ToInt64()
	if err != nil {
		log.Printf("Failed to parse cacheID from Lua result: %v", err)
		return -1
	}

	flushCacheSize, err := results[1].ToInt64()
	if err != nil {
		log.Printf("Failed to parse flush cache size from Lua result: %v", err)
		return -1
	}

	// Attach cacheID to message
	msg.CacheID = int(cacheID)
	log.Printf("Cached private message with ID %d. Flush cache size: %d", cacheID, flushCacheSize)

	// Trigger a flush if the flush cache is full
	if flushCacheSize >= maxCacheSize {
		log.Println("Flush cache size limit reached for private messages. Flushing to database...")
		m.FlushPrivateMessagesToDB()
	}

	return int(cacheID)
}

func (m *MessageCache) GetCachedPrivateMessages(userID string) []models.PrivateChatMessage {
	cacheKey := fmt.Sprintf("recent_private_messages:%s", userID)
	ctx := context.Background()

	cachedMessages, err := m.ValkeyClient.Do(
		ctx,
		m.ValkeyClient.B().Lrange().Key(cacheKey).Start(0).Stop(-1).Build(),
	).AsStrSlice()

	if err != nil {
		log.Printf("Failed to retrieve private messages for user %s: %v", userID, err)
		return nil
	}

	var privateMessages []models.PrivateChatMessage
	for _, jsonData := range cachedMessages {
		var cachedMsg struct {
			CacheID int64                     `json:"cache_id"`
			Data    models.PrivateChatMessage `json:"data"`
		}
		if err := json.Unmarshal([]byte(jsonData), &cachedMsg); err != nil {
			log.Printf("Failed to deserialize private message: %v", err)
			continue
		}
		cachedMsg.Data.CacheID = int(cachedMsg.CacheID)
		privateMessages = append(privateMessages, cachedMsg.Data)
	}

	log.Printf("Retrieved %d private messages from cache for user %s", len(privateMessages), userID)
	return privateMessages
}

// FlushPrivateMessagesToDB writes Valkey-cached private messages to PostgreSQL and clears the flush list
func (m *MessageCache) FlushPrivateMessagesToDB() {
	m.flushMutex.Lock()
	defer m.flushMutex.Unlock()

	ctx := context.Background()
	flushKey := "flush_private_messages"

	cachedMessages, err := m.ValkeyClient.Do(
		ctx,
		m.ValkeyClient.B().Lrange().Key(flushKey).Start(0).Stop(-1).Build(),
	).AsStrSlice()
	if err != nil {
		log.Printf("Failed to retrieve private messages from cache: %v", err)
		return
	}
	if len(cachedMessages) == 0 {
		log.Println("No private messages to flush to the database.")
		return
	}

	var messages []models.PrivateChatMessage
	for _, jsonData := range cachedMessages {
		var cached struct {
			CacheID int64                     `json:"cache_id"`
			Data    models.PrivateChatMessage `json:"data"`
		}
		if err := json.Unmarshal([]byte(jsonData), &cached); err != nil {
			log.Printf("Failed to deserialize private message: %v", err)
			continue
		}
		cached.Data.CacheID = int(cached.CacheID)
		messages = append(messages, cached.Data)
	}

	err = database.FlushPrivateMessages(m.DB, messages)
	if err != nil {
		log.Printf("Failed to flush private messages to database: %v", err)
		return
	}

	_, err = m.ValkeyClient.Do(ctx, m.ValkeyClient.B().Del().Key(flushKey).Build()).AsInt64()
	if err != nil {
		log.Printf("Failed to clear private message flush cache: %v", err)
		return
	}

	log.Printf("Successfully flushed %d private messages to the database.", len(messages))
}
