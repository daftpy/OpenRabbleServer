package db

import (
	"chatserver/internal/messages"
	"chatserver/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func FetchChannels(db *pgxpool.Pool) ([]models.Channel, error) {
	rows, err := db.Query(context.Background(), "SELECT name, description FROM chatserver.channels")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channels: %w", err)
	}
	defer rows.Close()

	var channels []models.Channel
	for rows.Next() {
		var channel models.Channel
		var description sql.NullString

		if err := rows.Scan(&channel.Name, &description); err != nil {
			return nil, fmt.Errorf("failed to scan channel row: %w", err)
		}

		// Convert sql.NullString to *string
		if description.Valid {
			channel.Description = &description.String
		} else {
			channel.Description = nil // Keep as nil if it was NULL in DB
		}

		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	log.Printf("Loaded %d channels from database", len(channels))
	return channels, nil
}

func FetchMessageCountByChannel(db *pgxpool.Pool) ([]messages.ChannelMessageCount, error) {
	rows, err := db.Query(context.Background(), `
		SELECT channel, COUNT(*) AS message_count
		FROM chatserver.chat_messages
		GROUP BY channel
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch message counts: %w", err)
	}
	defer rows.Close()

	var counts []messages.ChannelMessageCount
	for rows.Next() {
		var count messages.ChannelMessageCount
		if err := rows.Scan(&count.Channel, &count.MessageCount); err != nil {
			return nil, fmt.Errorf("failed to scan message count row: %w", err)
		}
		counts = append(counts, count)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	log.Println("Loaded message counts")
	return counts, nil
}
