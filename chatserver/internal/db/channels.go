package db

import (
	"chatserver/internal/messages"
	"chatserver/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateChannel(db *pgxpool.Pool, name string, description string) error {
	placeholderOwner := "00000000-0000-0000-0000-000000000000"

	// Get current max sort_order
	var maxOrder int
	err := db.QueryRow(context.Background(),
		`SELECT COALESCE(MAX(sort_order), 0) FROM chatserver.channels`,
	).Scan(&maxOrder)
	if err != nil {
		return fmt.Errorf("failed to get max sort_order: %w", err)
	}

	_, err = db.Exec(context.Background(), `
		INSERT INTO chatserver.channels (name, description, owner_id, sort_order)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO NOTHING
	`, name, description, placeholderOwner, maxOrder+1)

	return err
}

// FetchChannels retrieves all channels from the database's `channels` table.
// For each channel, if the `description` column is NULL, the Channel.Description
// field will be set to nil in the returned slice. Otherwise, it contains the
// description as a pointer to a string.
// It returns:
//  1. A slice of Channel models
//  2. An error, if any
func FetchChannels(db *pgxpool.Pool) ([]models.Channel, error) {
	rows, err := db.Query(context.Background(), "SELECT id, name, description FROM chatserver.channels ORDER BY sort_order, id")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channels: %w", err)
	}
	defer rows.Close()

	var channels []models.Channel
	for rows.Next() {
		var channel models.Channel
		var description sql.NullString

		if err := rows.Scan(&channel.ID, &channel.Name, &description); err != nil {
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

func UpdateChannel(db *pgxpool.Pool, ID int, name *string, description *string) error {
	clauses := []string{}
	params := []interface{}{}
	paramIndex := 1

	if name != nil {
		clauses = append(clauses, "name = $"+strconv.Itoa(paramIndex))
		params = append(params, *name)
		paramIndex++
	}

	if description != nil {
		clauses = append(clauses, "description = $"+strconv.Itoa(paramIndex))
		params = append(params, *description)
		paramIndex++
	}

	if len(clauses) == 0 {
		return errors.New("no fields provided to update")
	}

	query := `
        UPDATE chatserver.channels
        SET ` + strings.Join(clauses, ", ") + ` 
        WHERE id = $` + strconv.Itoa(paramIndex)

	params = append(params, ID)

	_, err := db.Exec(context.Background(), query, params...)
	return err
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

// MoveChannelBefore reorders a channel to appear before another channel.
// If beforeID is nil, it moves the channel to the end.
func MoveChannelBefore(db *pgxpool.Pool, movedID int, beforeID *int) error {
	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get current sort_order of the moved channel
	var currentOrder int
	err = tx.QueryRow(ctx, `SELECT sort_order FROM chatserver.channels WHERE id = $1`, movedID).Scan(&currentOrder)
	if err != nil {
		return fmt.Errorf("failed to get current sort_order: %w", err)
	}

	var newOrder int
	if beforeID != nil {
		// Get target sort_order of the channel we want to go before
		err = tx.QueryRow(ctx, `SELECT sort_order FROM chatserver.channels WHERE id = $1`, *beforeID).Scan(&newOrder)
		if err != nil {
			return fmt.Errorf("failed to get target sort_order: %w", err)
		}
	} else {
		// Move to the end
		err = tx.QueryRow(ctx, `SELECT COALESCE(MAX(sort_order), 0) + 1 FROM chatserver.channels`).Scan(&newOrder)
		if err != nil {
			return fmt.Errorf("failed to get max sort_order: %w", err)
		}
	}

	// Normalize the direction
	if newOrder > currentOrder {
		// Shift channels between current and new down by 1
		_, err = tx.Exec(ctx, `
			UPDATE chatserver.channels
			SET sort_order = sort_order - 1
			WHERE sort_order > $1 AND sort_order < $2
		`, currentOrder, newOrder)
	} else if newOrder < currentOrder {
		// Shift channels between new and current up by 1
		_, err = tx.Exec(ctx, `
			UPDATE chatserver.channels
			SET sort_order = sort_order + 1
			WHERE sort_order >= $1 AND sort_order < $2
		`, newOrder, currentOrder)
	}

	if err != nil {
		return fmt.Errorf("failed to shift surrounding channels: %w", err)
	}

	// Set moved channel to new sort_order
	_, err = tx.Exec(ctx, `
		UPDATE chatserver.channels
		SET sort_order = $1
		WHERE id = $2
	`, newOrder, movedID)

	if err != nil {
		return fmt.Errorf("failed to update moved channel: %w", err)
	}

	return tx.Commit(ctx)
}
