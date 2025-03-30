package cache

import (
	"chatserver/internal/models"
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

func (m *MessageCache) CheckRateLimitValkey(userID string, limit int, ttlSeconds int) (bool, error) {
	ctx := context.Background()

	rateLimitScript := valkey.NewLuaScript(`
		local key = KEYS[1]
		local limit = tonumber(ARGV[1])
		local ttl   = tonumber(ARGV[2])

		-- Increment the user's counter
		local current = redis.call("INCR", key)

		-- If this is the first increment (current == 1), set expiry
		if current == 1 then
			redis.call("EXPIRE", key, ttl)
		end

		-- If above the limit, return 0 (block), else 1 (allow)
		if current > limit then
			return 0
		else
			return 1
		end
	`)

	// Build the Redis key: "ratelimit:<userID>"
	rateKey := fmt.Sprintf("ratelimit:%s", userID)

	// Execute the Lua script
	result, err := rateLimitScript.Exec(
		ctx,
		m.ValkeyClient,
		[]string{rateKey}, // KEYS
		[]string{fmt.Sprintf("%d", limit), // ARGV[1]
			fmt.Sprintf("%d", ttlSeconds)}, // ARGV[2]
	).ToInt64()

	if err != nil {
		return false, err
	}

	// result == 1 is allow, result == 0 is block
	return (result == 1), nil
}

func (m *MessageCache) AttemptCacheWithRateLimit(userID string, msg models.ChatMessage) (int, error) {
	// Allow 10 messages per 60s
	allowed, err := m.CheckRateLimitValkey(userID, m.MessageLimit, m.WindowSeconds)
	if err != nil {
		return -1, fmt.Errorf("rate limit check failed: %v", err)
	}

	if !allowed {
		return -1, fmt.Errorf("rate limit exceeded for user %s", userID)
	}

	// If allowed, proceed to cache
	cacheID := m.CacheChatMessage(msg)
	if cacheID == -1 {
		return -1, fmt.Errorf("failed to cache message for user %s", userID)
	}
	return cacheID, nil
}
