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

// FetchChannels retrieves all channels from the database's `channels` table.
// For each channel, if the `description` column is NULL, the Channel.Description
// field will be set to nil in the returned slice. Otherwise, it contains the
// description as a pointer to a string.
// It returns:
//  1. A slice of Channel models
//  2. An error, if any
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

// FetchMessageCountByChannel returns the total number of chat messages per channel.
// It queries the `chat_messages` table, grouping by the `channel` column to produce
// a list of channels and their respective message counts.
// It returns:
//  1. A slice of ChannelMessageCount, where each item contains a channel name and a message count
//  2. An error, if the query or row scanning fails
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
