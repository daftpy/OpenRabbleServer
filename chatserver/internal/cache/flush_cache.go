package cache

import (
	"chatserver/internal/models"
	"context"
	"encoding/json"
	"log"
	"time"
)

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

	for _, jsonData := range cachedMessages {
		// Use a struct that matches the stored JSON format
		var cachedMsg struct {
			CacheID int64              `json:"cache_id"`
			Data    models.ChatMessage `json:"data"`
		}

		if err := json.Unmarshal([]byte(jsonData), &cachedMsg); err != nil {
			log.Printf("Failed to deserialize chat message: %v", err)
			continue
		}

		// Assign extracted cacheID to the payload
		cachedMsg.Data.CacheID = int(cachedMsg.CacheID)

		log.Printf("Inserting a message with the cacheID %d", cachedMsg.Data.CacheID)

		_, err = tx.Exec(
			ctx,
			`INSERT INTO chatserver.chat_messages (cache_id, owner_id, channel, message, authored_at)
			 VALUES ($1, $2, $3, $4, $5)
			`,
			cachedMsg.Data.CacheID, cachedMsg.Data.OwnerID, cachedMsg.Data.Channel, cachedMsg.Data.Message, cachedMsg.Data.Sent,
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
			m.FlushPrivateMessagesToDB()
		}
	}()
}
