package cache

import (
	"chatserver/internal/messages"
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
    local message = ARGV[1]
    local maxSize = tonumber(ARGV[2])

    -- Add to the recent circular cache
    redis.call("RPUSH", recentKey, message)
    redis.call("LTRIM", recentKey, -maxSize, -1)

    -- Add to the flush cache (no trimming)
    redis.call("RPUSH", flushKey, message)

    return redis.call("LLEN", flushKey)  -- Return the size of the flush cache
`)

const maxCacheSize = 500
const cacheTTL = 24 * 60 * 60 // 24 hours in seconds
const flushInterval = 2 * time.Minute

// Caches a chat message in Valkey and triggers a DB flush if max cache size is reached
func (m *MessageCache) CacheChatMessage(msg messages.ChatMessagePayload) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to serialize chat message: %v", err)
		return
	}

	recentCacheKey := "recent_messages"
	flushCacheKey := "flush_messages"
	ctx := context.Background()

	// Cache the message in both caches
	cacheSize, err := cacheMessageScript.Exec(
		ctx,
		m.ValkeyClient,
		[]string{recentCacheKey, flushCacheKey},
		[]string{string(jsonData), fmt.Sprintf("%d", maxCacheSize)},
	).AsInt64()

	if err != nil {
		log.Printf("Failed to cache chat message with Lua script: %v", err)
		return
	}

	log.Printf("Cached chat message: %s. Flush cache size: %d", msg.Message, cacheSize)

	// Flush to DB if the flush cache reaches the max size
	if cacheSize >= maxCacheSize {
		log.Println("Flush cache size limit reached. Flushing to the database...")
		m.FlushCacheToDB()
	}
}

// Retrieves chat messages from the circular cache
func (m *MessageCache) GetCachedChatMessages() []messages.ChatMessagePayload {
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

	var chatMessages []messages.ChatMessagePayload
	for _, jsonData := range cachedMessages {
		var msg messages.ChatMessagePayload
		if err := json.Unmarshal([]byte(jsonData), &msg); err != nil {
			log.Printf("Failed to deserialize chat message: %v", err)
			continue
		}
		chatMessages = append(chatMessages, msg)
	}

	log.Printf("Retrieved %d messages from recent cache", len(chatMessages))
	return chatMessages
}

// Flushes cached messages to the PostgreSQL database
func (m *MessageCache) FlushCacheToDB() {
	m.flushMutex.Lock()         // Acquire the lock
	defer m.flushMutex.Unlock() // Release the lock when done

	flushCacheKey := "flush_messages"
	ctx := context.Background()

	cachedMessages, err := m.ValkeyClient.Do(
		ctx,
		m.ValkeyClient.B().Lrange().Key(flushCacheKey).Start(0).Stop(-1).Build(),
	).AsStrSlice()

	if err != nil {
		log.Printf("Failed to retrieve cached messages from Valkey: %v", err)
		return
	}

	if len(cachedMessages) == 0 {
		log.Println("No messages to flush to the database.")
		return
	}

	log.Printf("Flushing %d messages to the database.", len(cachedMessages))

	tx, err := m.DB.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start database transaction: %v", err)
		return
	}
	defer tx.Rollback(ctx)

	const tempOwnerID = "0f5e28a8-4f8e-49be-b56d-83419ab92a36"

	for _, jsonData := range cachedMessages {
		var msg messages.ChatMessagePayload
		if err := json.Unmarshal([]byte(jsonData), &msg); err != nil {
			log.Printf("Failed to deserialize chat message: %v", err)
			continue
		}

		_, err = tx.Exec(
			ctx,
			`INSERT INTO chatserver.chat_messages (owner_id, channel, message, authored_at)
			 VALUES ($1, $2, $3, $4)
			`,
			tempOwnerID, msg.Channel, msg.Message, msg.Sent,
		)

		if err != nil {
			log.Printf("Failed to insert message into the database: %v", err)
			return
		}
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return
	}

	// Clear the flush cache to avoid duplicate inserts
	if _, err := m.ValkeyClient.Do(ctx, m.ValkeyClient.B().Del().Key(flushCacheKey).Build()).AsInt64(); err != nil {
		log.Printf("Failed to clear flush cache after database flush: %v", err)
	} else {
		log.Println("Successfully cleared flush cache after database flush.")
	}

	log.Println("Successfully flushed messages to the database.")
}

// StartPeriodicFlush triggers database flush every interval
func (m *MessageCache) StartPeriodicFlush() {
	ticker := time.NewTicker(flushInterval)
	go func() {
		for range ticker.C {
			log.Println("Periodic flush triggered.")
			m.FlushCacheToDB()
		}
	}()
}
