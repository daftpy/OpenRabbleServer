package cache

import (
	"chatserver/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valkey-io/valkey-go"
)

type MessageCache struct {
	ValkeyClient valkey.Client
	DB           *pgxpool.Pool
	flushMutex   sync.Mutex // Syncrhonize flush operations
}

// Lua script to handle both circular and flush caches
var cacheMessageScript = valkey.NewLuaScript(`
    local recentKey = KEYS[1]      -- Circular cache for recent messages
    local flushKey = KEYS[2]       -- Flush cache for database persistence
	local counterKey = KEYS[3]	   -- Counter key for cacheID
    local message = ARGV[1]
    local maxSize = tonumber(ARGV[2])

	-- Generate the cacheID
	local cacheID = redis.call("INCR", counterKey)

	-- Attach the cacheID to the message
	local enrichedMessage = cjson.encode({cache_id = cacheID, data = cjson.decode(message)})

    -- Add to the recent circular cache
    redis.call("RPUSH", recentKey, enrichedMessage)
    redis.call("LTRIM", recentKey, -maxSize, -1)

    -- Add to the flush cache (no trimming)
    redis.call("RPUSH", flushKey, enrichedMessage)

	-- Get the current size of the flush cache
	local flushCacheSize = redis.call("LLEN", flushKey)

    return {cacheID, flushCacheSize} -- Return the cacheID and size of the flush cache
`)

const maxCacheSize = 500
const cacheTTL = 24 * 60 * 60 // 24 hours in seconds
const flushInterval = 2 * time.Minute

// Caches a chat message in Valkey and triggers a DB flush if max cache size is reached
func (m *MessageCache) CacheChatMessage(msg models.ChatMessage) int {
	recentCacheKey := "recent_messages"
	flushCacheKey := "flush_messages"
	counterKey := "cache_message_id" // Key for unique message IDs

	ctx := context.Background()

	// Ensure JSON serialization is successful before passing to Lua
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to serialize chat message: %v", err)
		return -1
	}

	// Execute Lua script and extract both values
	results, err := cacheMessageScript.Exec(
		ctx,
		m.ValkeyClient,
		[]string{recentCacheKey, flushCacheKey, counterKey},         // Keys
		[]string{string(jsonData), fmt.Sprintf("%d", maxCacheSize)}, // Arguments
	).ToArray()

	if err != nil {
		log.Printf("Failed to cache chat message with Lua script: %v", err)
		return -1
	}

	// Extract values
	cacheID, err := results[0].ToInt64()
	if err != nil {
		log.Printf("Failed to parse cacheID: %v", err)
		return -1
	}

	flushCacheSize, err := results[1].ToInt64()
	if err != nil {
		log.Printf("Failed to parse flushCacheSize: %v", err)
		return -1
	}

	// Attach the correct `cacheID` to the message
	msg.CacheID = int(cacheID)

	log.Printf("Cached chat message with ID %d: %s. Flush cache size: %d", cacheID, msg.Message, flushCacheSize)

	// Flush to DB if the actual flush cache size reaches maxCacheSize
	if flushCacheSize >= maxCacheSize {
		log.Println("Flush cache size limit reached. Flushing to the database...")
		m.FlushCacheToDB()
	}

	return int(cacheID) // Correctly return cacheID
}

// Retrieves chat messages from the circular cache
func (m *MessageCache) GetCachedChatMessages() []models.ChatMessage {
	recentCacheKey := "recent_messages"
	ctx := context.Background()

	cachedMessages, err := m.ValkeyClient.Do(
		ctx,
		m.ValkeyClient.B().Lrange().Key(recentCacheKey).Start(0).Stop(-1).Build(),
	).AsStrSlice()

	if err != nil {
		log.Printf("Failed to retrieve cached messages from Valkey: %v", err)
		return nil
	}

	var chatMessages []models.ChatMessage
	for _, jsonData := range cachedMessages {
		// New struct to extract `cache_id` alongside the `data`
		var cachedMsg struct {
			CacheID int64              `json:"cache_id"`
			Data    models.ChatMessage `json:"data"`
		}

		// Unmarshal JSON into the new struct
		if err := json.Unmarshal([]byte(jsonData), &cachedMsg); err != nil {
			log.Printf("Failed to deserialize chat message: %v", err)
			continue
		}

		// Convert int64 to int before assignment
		cachedMsg.Data.CacheID = int(cachedMsg.CacheID)

		chatMessages = append(chatMessages, cachedMsg.Data)
	}

	log.Printf("Retrieved %d messages from recent cache", len(chatMessages))
	return chatMessages
}

func (m *MessageCache) DeleteCachedMessage(cacheID int) bool {
	var deleteMessageScript = valkey.NewLuaScript(`
	local recentKey = KEYS[1]
	local flushKey = KEYS[2]
	local cacheID = ARGV[1]

	-- Function to remove the message by cacheID
	local function removeMessage(key, cacheID)
		local messages = redis.call("LRANGE", key, 0, -1)
		for i, msg in ipairs(messages) do
			local decoded = cjson.decode(msg)
			if decoded.cache_id == tonumber(cacheID) then
				redis.call("LREM", key, 1, msg)
				return 1  -- Success, message deleted
			end
		end
		return 0  -- Message not found
	end

	-- Remove from both lists
	local removedRecent = removeMessage(recentKey, cacheID)
	local removedFlush = removeMessage(flushKey, cacheID)

	return removedRecent + removedFlush  -- Return number of deletions
	`)

	recentCacheKey := "recent_messages"
	flushCacheKey := "flush_messages"
	ctx := context.Background()

	// Execute the Lua script
	deleted, err := deleteMessageScript.Exec(
		ctx,
		m.ValkeyClient,
		[]string{recentCacheKey, flushCacheKey}, // KEYS
		[]string{fmt.Sprintf("%d", cacheID)},    // ARGV (cacheID as string)
	).ToInt64()

	if err != nil {
		log.Printf("Failed to delete message with cacheID %d: %v", cacheID, err)
		return false
	}

	// If deleted count > 0, the message was found and removed
	if deleted > 0 {
		log.Printf("Deleted message with cacheID %d from cache.", cacheID)
		return true
	}

	log.Printf("Message with cacheID %d not found in cache.", cacheID)
	return false
}
