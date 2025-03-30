package db

import (
	"chatserver/internal/models"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ensureDefaultRateLimit(db *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := db.Exec(ctx, `
		INSERT INTO chatserver.rate_limiter (owner_id, message_limit, window_seconds)
		VALUES ('default', 10, 60)
		ON CONFLICT (owner_id) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("failed to insert default rate limiter row: %w", err)
	}
	log.Println("Default rate limiter row ensured.")
	return nil
}

func UpdateRateLimiter(db *pgxpool.Pool, rateLimiterID, messageLimit, windowSeconds int) error {
	ctx := context.Background()
	_, err := db.Exec(ctx, `
		UPDATE chatserver.rate_limiter
		SET message_limit = $1,
			window_seconds = $2,
			updated_at = NOW()
		WHERE id = $3
	`, messageLimit, windowSeconds, rateLimiterID)

	if err != nil {
		return fmt.Errorf("Faileed to update rate limiter %d: %w", rateLimiterID, err)
	}

	return nil
}

func GetRateLimiterByID(db *pgxpool.Pool, rateLimiterID int) (models.RateLimiter, error) {
	ctx := context.Background()
	row := db.QueryRow(ctx, `
		SELECT id, owner_id, message_limit, window_seconds
		FROM chatserver.rate_limiter
		WHERE id = $1
	`, rateLimiterID)

	var rl models.RateLimiter

	err := row.Scan(&rl.ID, &rl.OwnerID, &rl.MessageLimit, &rl.WindowSeconds)
	if err != nil {
		return models.RateLimiter{}, fmt.Errorf("Failed to retrieve rate limiter %d: %w", rateLimiterID, err)
	}

	return rl, nil
}
